package paloalto

import (
	"cmp"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/nonsugar-go/tomato-ui/internal/utmconv/model"
)

func fillModel(app *model.App, palo *PaloAltoConfig) error {
	var err error

	if app.Tag, err = ToModelTags(palo.TagObject); err != nil {
		return err
	}
	slog.Info("タグが変換されました", "count", len(app.Tag))

	if app.Addresses, err = ToModelAddresses(palo.Addresses); err != nil {
		return err
	}
	slog.Info("アドレスが変換されました", "count", len(app.Addresses))

	if app.AddressGroups, err = ToModelAddressGroups(palo.AddressGroups); err != nil {
		return err
	}
	slog.Info("アドレスグループが変換されました", "count", len(app.AddressGroups))

	if app.Services, err = ToModelServices(palo.Services); err != nil {
		return err
	}
	slog.Info("サービスが変換されました", "count", len(app.Services))

	if app.ServiceGroups, err = ToModelServiceGroups(palo.ServiceGroups); err != nil {
		return err
	}
	slog.Info("サービスグループが変換されました", "count", len(app.ServiceGroups))

	if app.Policies, err = ToModelPolicies(palo.SecurityRules, app.AppConfig.PaloAlto.Conf.ApplicationDefaultReplacementMap.Value); err != nil {
		return err
	}
	slog.Info("ポリシーが変換されました", "count", len(app.Policies))

	// app.NATRules, err = ...

	return nil
}

// ToModelTags converts Palo Alto tag objects to model tags.
// NOTE: untested
func ToModelTags(tagObjects []ScopedTagObject) ([]model.Tag, error) {
	var result []model.Tag

	for _, st := range tagObjects {
		t := st.TagObject

		tag := model.Tag{
			Value:       t.Name,
			Color:       t.Color,
			Description: t.Comments,
		}

		result = append(result, tag)
	}

	return result, nil
}

func ToModelAddresses(scopedAddrs []ScopedAddress) ([]model.Address, error) {
	var result []model.Address

	for _, sa := range scopedAddrs {
		a := sa.Address

		addr := model.Address{
			Name:        a.Name,
			Description: a.Description,
			Tags:        a.Tags,
		}

		switch {
		case a.IPNetmask != "":
			addr.Type = model.AddressTypeIPNetmask
			addr.Value = a.IPNetmask

		case a.FQDN != "":
			addr.Type = model.AddressTypeFQDN
			addr.Value = a.FQDN

		default:
			addr.Type = model.AddressTypeUnknown
			addr.Value = ""
		}

		result = append(result, addr)
	}

	return result, nil
}

func ToModelAddressGroups(scopedAddrGrps []ScopedAddressGroup) ([]model.AddressGroup, error) {
	var result []model.AddressGroup

	for _, sg := range scopedAddrGrps {
		g := sg.Group

		grp := model.AddressGroup{
			Name:        g.Name,
			Description: g.Description,
			Tags:        g.Tags,
		}

		switch {
		case len(g.Static) != 0:
			grp.Members = g.Static
		case g.Dynamic != nil:
			grp.Members = []string{g.Dynamic.Filter}
		default:
			grp.Members = []string{}
		}

		result = append(result, grp)
	}

	return result, nil
}

func ToModelServices(scopedSvcs []ScopedService) ([]model.Service, error) {
	var result []model.Service

	// result = append(result, model.Service{
	// 	Name:  "service-http",
	// 	Type:  model.ServiceTypeTCP,
	// 	Ports: "80",
	// })

	// result = append(result, model.Service{
	// 	Name:  "service-https",
	// 	Type:  model.ServiceTypeTCP,
	// 	Ports: "443",
	// })

	for _, ss := range scopedSvcs {
		s := ss.Service

		svc := model.Service{
			Name:        s.Name,
			Description: s.Description,
			Tags:        s.Tags,
		}

		switch {
		case s.Protocol.TCP != nil:
			svc.Type = model.ServiceTypeTCP
			svc.Ports = s.Protocol.TCP.Port
			svc.SourcePorts = s.Protocol.TCP.SourcePort

		case s.Protocol.UDP != nil:
			svc.Type = model.ServiceTypeUDP
			svc.Ports = s.Protocol.UDP.Port
			svc.SourcePorts = s.Protocol.UDP.SourcePort

		default:
			svc.Type = model.ServiceTypeUnknown
			svc.Ports = ""
			svc.SourcePorts = ""

		}

		result = append(result, svc)
	}

	return result, nil
}

func ToModelServiceGroups(scopedSvcGrps []ScopedServiceGroup) ([]model.ServiceGroup, error) {
	var result []model.ServiceGroup

	for _, sg := range scopedSvcGrps {
		g := sg.Group

		grp := model.ServiceGroup{
			Name:        g.Name,
			Members:     g.Members,
			Description: g.Description,
			Tags:        g.Tags,
		}

		result = append(result, grp)
	}

	return result, nil
}

func getUniqueTags(groupTag string, tags []string) []string {
	result := make([]string, 0, len(tags)+1)
	seen := make(map[string]struct{})

	add := func(t string) {
		if t != "" {
			if _, ok := seen[t]; !ok {
				seen[t] = struct{}{}
				result = append(result, t)
			}
		}
	}

	add(groupTag)

	for _, t := range tags {
		add(t)
	}

	if len(result) == 0 {
		return nil
	}
	return result
}

func ToModelPolicies(scopedSecurities []ScopedSecurity, appSvcMap []model.AppSvcMap) ([]model.Policy, error) {
	var result []model.Policy

	scopeRulebase := func(scope, rulebase string) string {
		switch rulebase {
		case "":
			return scope
		}
		return fmt.Sprintf("%s-%s", scope, rulebase)
	}

	for _, ss := range scopedSecurities {
		rule := ss.SecurityRule

		policy := model.Policy{
			Name:        rule.Name,
			Description: rule.Description,
			Enabled:     rule.Disabled != "yes",

			Match: model.PolicyMatch{
				FromZones: rule.FromZones,
				ToZones:   rule.ToZones,

				Sources:      toAddrRefs(rule.Sources),
				Destinations: toAddrRefs(rule.Destinations),

				Applications: rule.Applications,
				Services:     toSvcRefs(rule.Services, rule.Applications, appSvcMap),

				Users: rule.SourceUsers,
				HIPs:  append(rule.SourceHIP, rule.DestinationHIP...),

				NegateSource:      rule.NegateSource == "yes",
				NegateDestination: rule.NegateDestination == "yes",
			},

			Action: model.PolicyAction{
				Type:     toAction(rule.Action),
				Profiles: extractProfiles(rule.ProfileSetting),
			},

			Logging: model.Logging{
				LogAtStart: rule.LogStart == "yes", // default: no
				LogAtEnd:   rule.LogEnd != "no",    // default: yes
				LogProfile: rule.LogSetting,
			},

			Schedule: rule.Schedule,
			Tags:     getUniqueTags(rule.GroupTag, rule.Tags),
			Group:    rule.GroupTag,

			Scope: scopeRulebase(ss.Scope, ss.Rulebase),
		}

		result = append(result, policy)
	}

	return result, nil
}

func toAddrRefs(names []string) []model.AddressRef {
	if len(names) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(names))
	refs := make([]model.AddressRef, 0, len(names))

	for _, n := range names {
		n = normalizeName(n)
		if n == "" {
			continue
		}

		// 重複排除
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}

		refs = append(refs, model.AddressRef{Name: n})
	}

	return refs
}

func normalizeName(s string) string {
	s = strings.TrimSpace(s)

	// PaloAlto XMLでたまに来る
	if s == "any" {
		return ""
	}

	return s
}

func toSvcRefs(names []string, apps []string, appSvcMap []model.AppSvcMap) []model.ServiceRef {
	if len(names) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(names))
	refs := make([]model.ServiceRef, 0, len(names))

	for _, n := range names {
		n = normalizeServiceName(n)

		// application-default or any
		if n == "application-default" || n == "" {
			if len(apps) == 0 {
				continue
			}
			for _, a := range apps {
				normalizeApp := normalizeName(a)
				appToService := []string{}

				// switch normalizeApp {
				// case "":
				// 	// "any"
				// 	continue
				// case "icmp":
				// 	appToService = append(appToService, "icmp-proto")
				// case "ping":
				// 	appToService = append(appToService, "echo-request")
				// case "traceroute":
				// 	appToService = append(appToService, "traceroute")
				// case "ssh":
				// 	appToService = append(appToService, "ssh")
				// case "syslog":
				// 	appToService = append(appToService, "syslog")
				// case "ipsec-esp", "ipsec-esp-udp":
				// 	appToService = append(appToService, "ESP")
				// 	appToService = append(appToService, "IKE")
				// 	appToService = append(appToService, "IKE_NAT_TRAVERSAL")
				// default:
				// 	slog.Warn("cannot handle application-default for app", "app", a)
				// }
				if normalizeApp == "" {
					// "any"
					continue
				} else {
					isFound := false
					for _, v := range appSvcMap {
						if normalizeApp == v.Application {
							isFound = true
							appToService = append(appToService, v.Services...)
							break
						}
					}
					if !isFound {
						slog.Warn("cannot handle application-default for app", "app", a)
					}
				}
				for _, s := range appToService {
					if _, ok := seen[s]; ok {
						continue
					}
					seen[s] = struct{}{}

					refs = append(refs, model.ServiceRef{Name: s})
				}
			}
			continue
		}

		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}

		refs = append(refs, model.ServiceRef{Name: n})
	}

	slices.SortFunc(refs, func(a, b model.ServiceRef) int {
		return cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})

	return refs
}

func normalizeServiceName(s string) string {
	s = strings.TrimSpace(s)

	// PaloAltoの "any"
	if s == "any" {
		return ""
	}

	return s
}

func toAction(s string) model.ActionType {
	s = normalizeAction(s)

	switch s {
	case "allow":
		return model.ActionAllow
	case "deny":
		return model.ActionDeny
	case "drop":
		return model.ActionDrop

	// PaloAlto系
	case "reset-client", "reset-server", "reset-both":
		return model.ActionReset

	default:
		// 不明は deny 扱いに寄せる（安全側）
		return model.ActionDeny
	}
}

func normalizeAction(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// extractProfiles extracts profile names from ProfileSetting.
// It handles both group-based and individual profile settings.
// NOTE: untested
func extractProfiles(ps *ProfileSetting) []string {
	if ps == nil {
		return nil
	}

	var out []string

	// パターン1: group指定
	if len(ps.Group) != 0 {
		return ps.Group
	}

	// パターン2: 個別プロファイル
	if ps.Profiles != nil { // NOTE: untested
		if ps.Profiles.AV != nil {
			for _, m := range ps.Profiles.AV {
				m = strings.TrimSpace(m)
				if m == "" {
					continue
				}
				out = append(out, m)
			}
		}
		if ps.Profiles.VP != nil {
			for _, m := range ps.Profiles.VP {
				m = strings.TrimSpace(m)
				if m == "" {
					continue
				}
				out = append(out, m)
			}
		}
		if ps.Profiles.AS != nil {
			for _, m := range ps.Profiles.AS {
				m = strings.TrimSpace(m)
				if m == "" {
					continue
				}
				out = append(out, m)
			}
		}
		if ps.Profiles.URL != nil {
			for _, m := range ps.Profiles.URL {
				m = strings.TrimSpace(m)
				if m == "" {
					continue
				}
				out = append(out, m)
			}
		}
		if ps.Profiles.FB != nil {
			for _, m := range ps.Profiles.FB {
				m = strings.TrimSpace(m)
				if m == "" {
					continue
				}
				out = append(out, m)
			}
		}
		if ps.Profiles.DF != nil {
			for _, m := range ps.Profiles.DF {
				m = strings.TrimSpace(m)
				if m == "" {
					continue
				}
				out = append(out, m)
			}
		}
		if ps.Profiles.WFA != nil {
			for _, m := range ps.Profiles.WFA {
				m = strings.TrimSpace(m)
				if m == "" {
					continue
				}
				out = append(out, m)
			}
		}
	}

	if len(out) == 0 {
		return nil
	}

	return out
}
