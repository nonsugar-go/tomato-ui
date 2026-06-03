package checkpoint

import (
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"github.com/nonsugar-go/tomato-ui/internal/utmconv/model"
)

// PaloAlto から check point への変換に特化したコードであるため、将来的に削除する可能性がある
// 脅威防御ポリシーを生成するためのコードを記述する
// PaloAlto の Seciryt Policy で action が allow かつ profile group が設定されているルールを対象とする
// PaloAlto の source, destinattion, service を脅威防御ルールの条件とする
func ConvertThreatPolicies(policies []model.Policy, ctx *Context) ([]string, error) {
	var results []string
	var errs []error

	grouped := GroupPoliciesByScope(policies)
	scopes := make([]string, 0, len(grouped))
	for s := range grouped {
		scopes = append(scopes, s)
	}
	sort.Strings(scopes)

	existProfileGroup := func(p *model.Policy) bool {
		if _, ok := p.Extensions["paloalto-profile-setting-group"]; !ok {
			return false
		}
		return true
	}

	for _, scope := range scopes {
		// action が allow かつ profile group が設定されているルールを対象とする
		var filteredPolicies []*model.Policy
		for _, p := range grouped[scope] {
			if p.Action.Type == model.ActionAllow && existProfileGroup(p) {
				filteredPolicies = append(filteredPolicies, p)
			}
		}
		numOfPoliciesInScope := len(filteredPolicies)

		results = append(results, fmt.Sprintf("#\n# --- %s (%d) ---", scope, numOfPoliciesInScope))

		slog.Info("Check Point の脅威防御ポリシーを処理中",
			slog.String("scope", scope),
			slog.Int("count", numOfPoliciesInScope),
		)

		for _, policy := range filteredPolicies {
			line, err := ConvertThreatPolicy(*policy, ctx)
			if err != nil {
				errs = append(errs, fmt.Errorf("%s: %w", policy.Name, err))
				continue
			}
			results = append(results, line)
		}
	}

	return results, errors.Join(errs...)
}

// ConvertThreatPolicy is a function that converts a PaloAlto Security Policy to a Check Point Threat Prevention Policy.
func ConvertThreatPolicy(p model.Policy, ctx *Context) (string, error) {
	var errs []error
	var sb strings.Builder

	sb.WriteString(`add threat-rule layer "`)
	sb.WriteString(ctx.App.AppConfig.CheckPoint.Cli.ThreatRuleLayer.Value)
	sb.WriteString(`" position "bottom" name "`)
	sb.WriteString(p.Name)
	sb.WriteString(`"`)
	if p.Enabled == false {
		sb.WriteString(` enabled false`)
	}

	src, err := mapStrings(p.Match.Sources.Names(), ctx.AddrMap)
	errs = append(errs, err)
	dst, err := mapStrings(p.Match.Destinations.Names(), ctx.AddrMap)
	errs = append(errs, err)
	svc, err := mapStrings(p.Match.Services.Names(), ctx.SvcMap)
	errs = append(errs, err)

	// TODO: Workaround for "ICMP Protocol" はアプリケーションなのでダメ
	var svc2 []string
	for _, s := range svc {
		if s == "ICMP Protocol" {
			// 何もしない
		} else {
			svc2 = append(svc2, s)
		}
	}

	if p.Match.NegateSource {
		sb.WriteString(` source-negate true`)
	}
	sb.WriteString(` protected-scope "Any"`)
	sb.WriteString(buildIndexedKVWithDefaultAny("source", src))
	if p.Match.NegateDestination {
		sb.WriteString(` destination-negate true`)
	}
	sb.WriteString(buildIndexedKVWithDefaultAny("destination", dst))
	// service-negate is not implemented yet.
	sb.WriteString(buildIndexedKVWithDefaultAny("service", svc2))
	sb.WriteString(` track "Log"`)
	sb.WriteString(buildComment(p.Description))
	if p.Description != "" && len(p.Tags) > 0 {
		buildComment(p.Tags[0])
	}

	return sb.String(), errors.Join(errs...)
}
