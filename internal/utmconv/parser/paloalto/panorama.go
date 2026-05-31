package paloalto

import (
	"encoding/xml"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strings"

	"github.com/nonsugar-go/tomato-ui/internal/utmconv/exporter/excel"
	"github.com/nonsugar-go/tomato-ui/internal/utmconv/model"
)

func BuildPaloAltoConfig(c *Config) *PaloAltoConfig {
	var result PaloAltoConfig

	// Tags
	for _, tag := range c.Shared.Tags {
		result.TagObject = append(result.TagObject, ScopedTagObject{
			Scope:     "shared",
			TagObject: tag,
		})
	}

	for _, dg := range c.Devices.DeviceGroups {
		for _, tag := range dg.Tags {
			result.TagObject = append(result.TagObject, ScopedTagObject{
				Scope:     dg.Name,
				TagObject: tag,
			})
		}
	}

	// Addresses
	for _, addr := range c.Shared.Addresses {
		result.Addresses = append(result.Addresses, ScopedAddress{
			Scope:   "shared",
			Address: addr,
		})
	}

	for _, dg := range c.Devices.DeviceGroups {
		for _, addr := range dg.Addresses {
			result.Addresses = append(result.Addresses, ScopedAddress{
				Scope:   dg.Name,
				Address: addr,
			})
		}
	}

	// Address Groups
	for _, addrGrp := range c.Shared.AddressGroups {
		result.AddressGroups = append(result.AddressGroups, ScopedAddressGroup{
			Scope: "shared",
			Group: addrGrp,
		})
	}

	for _, dg := range c.Devices.DeviceGroups {
		for _, addrGrp := range dg.AddressGroups {
			result.AddressGroups = append(result.AddressGroups, ScopedAddressGroup{
				Scope: dg.Name,
				Group: addrGrp,
			})
		}
	}

	// Services
	for _, svc := range c.Shared.Services {
		result.Services = append(result.Services, ScopedService{
			Scope:   "shared",
			Service: svc,
		})
	}

	for _, dg := range c.Devices.DeviceGroups {
		for _, svc := range dg.Services {
			result.Services = append(result.Services, ScopedService{
				Scope:   dg.Name,
				Service: svc,
			})
		}
	}

	// Service Groups
	for _, svcGrp := range c.Shared.ServiceGroups {
		result.ServiceGroups = append(result.ServiceGroups, ScopedServiceGroup{
			Scope: "shared",
			Group: svcGrp,
		})
	}

	for _, dg := range c.Devices.DeviceGroups {
		for _, svcGrp := range dg.ServiceGroups {
			result.ServiceGroups = append(result.ServiceGroups, ScopedServiceGroup{
				Scope: dg.Name,
				Group: svcGrp,
			})
		}
	}

	// Security Rules
	for _, rule := range c.Shared.PreRulebase.Security.Rules {
		result.SecurityRules = append(result.SecurityRules, ScopedSecurity{
			Scope:        "shared",
			Rulebase:     "pre",
			SecurityRule: rule,
		})
	}

	for _, rule := range c.Shared.PostRulebase.Security.Rules {
		result.SecurityRules = append(result.SecurityRules, ScopedSecurity{
			Scope:        "shared",
			Rulebase:     "post",
			SecurityRule: rule,
		})
	}

	for _, dg := range c.Devices.DeviceGroups {
		for _, rule := range dg.PreRulebase.Security.Rules {
			result.SecurityRules = append(result.SecurityRules, ScopedSecurity{
				Scope:        dg.Name,
				Rulebase:     "pre",
				SecurityRule: rule,
			})
		}
	}

	for _, dg := range c.Devices.DeviceGroups {
		for _, rule := range dg.PostRulebase.Security.Rules {
			result.SecurityRules = append(result.SecurityRules, ScopedSecurity{
				Scope:        dg.Name,
				Rulebase:     "post",
				SecurityRule: rule,
			})
		}
	}

	// NAT Rules
	for _, rule := range c.Shared.PreRulebase.Nat.Rules {
		result.NATRules = append(result.NATRules, ScopedNAT{
			Scope:    "shared",
			Rulebase: "pre",
			NATRule:  rule,
		})
	}

	for _, rule := range c.Shared.PostRulebase.Nat.Rules {
		result.NATRules = append(result.NATRules, ScopedNAT{
			Scope:    "shared",
			Rulebase: "post",
			NATRule:  rule,
		})
	}

	for _, dg := range c.Devices.DeviceGroups {
		for _, rule := range dg.PreRulebase.Nat.Rules {
			result.NATRules = append(result.NATRules, ScopedNAT{
				Scope:    dg.Name,
				Rulebase: "pre",
				NATRule:  rule,
			})
		}
	}

	for _, dg := range c.Devices.DeviceGroups {
		for _, rule := range dg.PostRulebase.Nat.Rules {
			result.NATRules = append(result.NATRules, ScopedNAT{
				Scope:    dg.Name,
				Rulebase: "post",
				NATRule:  rule,
			})
		}
	}

	return &result
}

func printTags(scope string, tags []TagObject, e *excel.Excel) {
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Name < tags[j].Name
	})

	for _, tag := range tags {
		e.Println(scope, tag.Name, tag.Color, colorName(tag.Color), tag.Comments)
	}
}

func colorName(code string) string {
	if name, ok := colorMap[code]; ok {
		return name
	}
	return code
}

func printAddresses(scope string, addrs []Address, e *excel.Excel) {
	for _, addr := range addrs {
		typ, value := resolveAddress(addr)

		e.Println(scope, addr.Name, typ, value, strings.Join(addr.Tags, ";"), addr.Description)
	}
}

func resolveAddress(addr Address) (string, string) {
	switch {
	case addr.IPNetmask != "":
		return "ip-netmask", addr.IPNetmask
	case addr.IPRange != "": // NOTE: untested
		return "ip-range", addr.IPRange
	case addr.IPWildcard != "": // NOTE: untested
		return "ip-wildcard", addr.IPWildcard
	case addr.FQDN != "":
		return "fqdn", addr.FQDN
	default:
		return "unknown", ""
	}
}

func printAddressGroups(scope string, groups []AddressGroup, e *excel.Excel) {
	for _, g := range groups {
		typ, value := resolveGroup(g)

		e.Println(scope, g.Name, typ, value, strings.Join(g.Tags, ";"), g.Description)
	}
}

func resolveGroup(g AddressGroup) (string, string) {
	if len(g.Static) > 0 {
		return "static", strings.Join(g.Static, ";")
	}
	if g.Dynamic != nil { // NOTE: untested
		return "dynamic", g.Dynamic.Filter
	}
	return "unknown", ""
}

func printServices(scope string, services []Service, e *excel.Excel) {
	for _, svc := range services {
		proto, port, srcPort := resolveService(svc)

		e.Println(scope, svc.Name, proto, port, srcPort, strings.Join(svc.Tags, ","), svc.Description)
	}
}

func resolveService(svc Service) (string, string, string) {
	if svc.Protocol.TCP != nil {
		return "tcp", svc.Protocol.TCP.Port, svc.Protocol.TCP.SourcePort
	}
	if svc.Protocol.UDP != nil {
		return "udp", svc.Protocol.UDP.Port, svc.Protocol.UDP.SourcePort
	}
	return "unknown", "", ""
}

func printServiceGroups(scope string, groups []ServiceGroup, e *excel.Excel) {
	for _, g := range groups {
		e.Println(
			scope,
			g.Name,
			strings.Join(g.Members, ";"),
			strings.Join(g.Tags, ";"),
			g.Description,
		)
	}
}

func printApplications(scope string, apps []Application, e *excel.Excel) {
	for _, app := range apps {
		port := ""
		if app.Default != nil {
			port = app.Default.Port
		}

		e.Println(scope, app.Name, app.Category, app.Subcategory, app.Technology,
			app.Risk, port, strings.Join(app.Tags, ";"), app.Description)
	}
}

func printApplicationGroups(scope string, groups []ApplicationGroup, e *excel.Excel) {
	for _, g := range groups {
		e.Println(scope, g.Name, strings.Join(g.Members, ";"), strings.Join(g.Tags, ","), g.Description)
	}
}

func printSecurityRules(scope string, rules []SecurityRule, e *excel.Excel) {
	for _, rule := range rules {

		from := strings.Join(rule.FromZones, ";")
		to := strings.Join(rule.ToZones, ";")
		src := strings.Join(rule.Sources, ";")
		dst := strings.Join(rule.Destinations, ";")
		app := strings.Join(rule.Applications, ";")
		svc := strings.Join(rule.Services, ";")
		sourceUsers := strings.Join(rule.SourceUsers, ";")
		categories := strings.Join(rule.Categories, ";")
		tags := strings.Join(rule.Tags, ";")

		if rule.NegateSource == "yes" {
			src = "NOT(" + src + ")"
		}
		if rule.NegateDestination == "yes" {
			dst = "NOT(" + dst + ")"
		}

		srcHIP := strings.Join(rule.SourceHIP, ";")
		dstHIP := strings.Join(rule.DestinationHIP, ";")

		profileGroup := ""
		if rule.ProfileSetting != nil {
			if len(rule.ProfileSetting.Group) > 0 {
				profileGroup = strings.Join(rule.ProfileSetting.Group, ";")
			} else if rule.ProfileSetting.Profiles != nil {
				profileGroup = expandProfiles(rule.ProfileSetting.Profiles)
			}
		}

		target := "all"

		if rule.Target != nil {
			var devs []string
			for _, d := range rule.Target.Devices {
				devs = append(devs, d.Name)
			}

			if len(devs) > 0 {
				if rule.Target.Negate == "yes" {
					target = "NOT(" + strings.Join(devs, ";") + ")"
				} else {
					target = strings.Join(devs, ";")
				}
			}
		}

		e.Println(scope, rule.Name, tags, rule.GroupTag,
			from, to, src, dst, app, svc, rule.Action,
			rule.Disabled, rule.LogStart, rule.LogEnd, rule.LogSetting,
			rule.Schedule, sourceUsers, categories, srcHIP, dstHIP,
			profileGroup, target, rule.Description)
	}
}

func expandProfiles(p *Profiles) string {
	if p == nil {
		return ""
	}

	var parts []string

	add := func(name string, v []string) {
		if len(v) == 0 {
			return
		}
		parts = append(parts, name+"="+strings.Join(v, ";"))
	}

	add("AV", p.AV)
	add("AS", p.AS)
	add("VP", p.VP)
	add("URL", p.URL)
	add("FB", p.FB)
	add("DF", p.DF)
	add("WFA", p.WFA)

	return strings.Join(parts, " ")
}

func printNatRules(scope string, rules []NATRule, e *excel.Excel) {
	for _, rule := range rules {

		// ----------------------------
		// Source NAT
		// ----------------------------
		srcTrans := ""

		if rule.SourceTranslation != nil {
			switch {
			case rule.SourceTranslation.DynamicIPAndPort != nil:
				d := rule.SourceTranslation.DynamicIPAndPort

				// Interface NAT
				if d.InterfaceAddress != nil {
					srcTrans = fmt.Sprintf(
						"Interface:%s IP:%s Floating:%s",
						d.InterfaceAddress.Interface,
						d.InterfaceAddress.Ip,
						d.InterfaceAddress.FloatingIp,
					)
				} else {
					// Dynamic IP and Port NAT
					srcTrans = "DIPP:" + strings.Join(d.TranslatedAddress, ",")
				}

			case rule.SourceTranslation.StaticIP != nil:
				s := rule.SourceTranslation.StaticIP
				srcTrans = fmt.Sprintf(
					"Static:%s BiDir:%s",
					s.TranslatedAddress,
					s.BiDirectional,
				)
			}
		}

		// ----------------------------
		// Destination NAT
		// ----------------------------
		dstTrans := ""

		if rule.DestinationTranslation != nil {
			d := rule.DestinationTranslation

			switch {
			case d.TranslatedAddress != "" && d.TranslatedPort != "":
				dstTrans = fmt.Sprintf("%s:%s",
					d.TranslatedAddress,
					d.TranslatedPort,
				)

			case d.TranslatedAddress != "":
				dstTrans = d.TranslatedAddress
			}
		}

		// ----------------------------
		// Output row
		// ----------------------------
		e.Println(
			scope,
			rule.Name,
			rule.Disabled,
			strings.Join(rule.FromZones, ";"),
			strings.Join(rule.ToZones, ";"),
			strings.Join(rule.Sources, ";"),
			strings.Join(rule.Destinations, ";"),
			rule.Service,
			srcTrans,
			dstTrans,
			rule.Description,
		)
	}
}

func parseXML(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := xml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func ParsePanorama(app *model.App) {
	var config *Config
	var err error

	if config, err = parseXML(app.Filename); err != nil {
		slog.Error("Failed to parse XML", "error", err)
		os.Exit(1)
	}

	panoramaConfig := BuildPaloAltoConfig(config)
	fillModel(app, panoramaConfig)

	e := excel.NewExcel("panorama.xlsx")
	defer e.Close()

	e.NewSheet("Tags")
	e.Println("Scope", "Name", "Color", "Color Name", "Comments")
	printTags("shared", config.Shared.Tags, e)
	for _, dg := range config.Devices.DeviceGroups {
		printTags(dg.Name, dg.Tags, e)
	}
	e.AddTable()

	e.NewSheet("Addresses")
	e.Println("Scope", "Name", "Type", "Value", "Tags", "Description")
	printAddresses("shared", config.Shared.Addresses, e) // NOTE: untested
	for _, dg := range config.Devices.DeviceGroups {
		printAddresses(dg.Name, dg.Addresses, e)
	}
	e.AddTable()

	e.NewSheet("Address Groups")
	e.Println("Scope", "Name", "Type", "Value", "Tags", "Description")
	printAddressGroups("shared", config.Shared.AddressGroups, e)
	for _, dg := range config.Devices.DeviceGroups {
		printAddressGroups(dg.Name, dg.AddressGroups, e)
	}
	e.AddTable()

	e.NewSheet("Services")
	e.Println("Scope", "Name", "Protocol", "Port", "SourcePort", "Tags", "Description")
	printServices("shared", config.Shared.Services, e)
	for _, dg := range config.Devices.DeviceGroups {
		printServices(dg.Name, dg.Services, e)
	}
	e.AddTable()

	e.NewSheet("Service Groups")
	e.Println("Scope", "Name", "Value", "Tags", "Description")
	printServiceGroups("shared", config.Shared.ServiceGroups, e)
	for _, dg := range config.Devices.DeviceGroups {
		printServiceGroups(dg.Name, dg.ServiceGroups, e)
	}
	e.AddTable()

	e.NewSheet("Applications")
	e.Println("Scope", "Name", "Category", "Subcategory", "Technology", "Risk",
		"Port", "Tags", "Description",
	)
	printApplications("shared", config.Shared.Applications, e)
	for _, dg := range config.Devices.DeviceGroups {
		printApplications(dg.Name, dg.Applications, e)
	}
	e.AddTable()

	e.NewSheet("Application Groups")
	e.Println("Scope", "Name", "Members", "Tags", "Description")
	printApplicationGroups("shared", config.Shared.ApplicationGroups, e)
	for _, dg := range config.Devices.DeviceGroups {
		printApplicationGroups(dg.Name, dg.ApplicationGroups, e)
	}
	e.AddTable()

	e.NewSheet("Security Rules")
	e.Println("Scope", "Name", "Tags", "GroupTag", "From", "To", "Source", "Destination", "Application",
		"Service", "Action", "Disabled", "LogStart", "LogEnd", "LogSetting", "Schedule", "SourceUser",
		"Category", "SourceHIP", "DestinationHIP", "Profile", "Target", "Description")

	printSecurityRules("shared-pre", config.Shared.PreRulebase.Security.Rules, e)
	printSecurityRules("shared-post", config.Shared.PostRulebase.Security.Rules, e)
	for _, dg := range config.Devices.DeviceGroups {
		printSecurityRules(dg.Name+"-pre", dg.PreRulebase.Security.Rules, e)
		printSecurityRules(dg.Name+"-post", dg.PostRulebase.Security.Rules, e)
	}
	e.AddTable()

	e.NewSheet("NAT Rules")
	e.Println("Scope", "Name", "Disabled", "From Zones", "To Zones", "Sources", "Destinations",
		"Service", "Source Translation", "Destination Translation", "Description")
	printNatRules("shared-pre", config.Shared.PreRulebase.Nat.Rules, e)
	printNatRules("shared-post", config.Shared.PostRulebase.Nat.Rules, e)
	for _, dg := range config.Devices.DeviceGroups {
		printNatRules(dg.Name+"-pre", dg.PreRulebase.Nat.Rules, e)
		printNatRules(dg.Name+"-post", dg.PostRulebase.Nat.Rules, e)
	}
	e.AddTable()
}
