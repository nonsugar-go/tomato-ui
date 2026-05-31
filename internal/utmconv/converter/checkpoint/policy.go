package checkpoint

import (
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"github.com/nonsugar-go/tomato-ui/internal/utmconv/model"
)

func GroupPoliciesByScope(policies []model.Policy) map[string][]*model.Policy {
	grouped := make(map[string][]*model.Policy)

	for i := range policies {
		p := &policies[i]
		scope := p.Scope

		if scope == "" {
			scope = "default"
		}

		grouped[scope] = append(grouped[scope], p)
	}

	return grouped
}

func ConvertPolicies(policies []model.Policy, ctx *Context) ([]string, error) {
	var results []string
	var errs []error

	grouped := GroupPoliciesByScope(policies)
	scopes := make([]string, 0, len(grouped))
	for s := range grouped {
		scopes = append(scopes, s)
	}
	sort.Strings(scopes)

	for _, scope := range scopes {
		policiesInScope := grouped[scope]
		numOfPoliciesInScope := len(policiesInScope)

		results = append(results, fmt.Sprintf("#\n# --- %s (%d) ---", scope, numOfPoliciesInScope))

		slog.Info("Check Point のアクセス コントロール ポリシーを処理中",
			slog.String("scope", scope),
			slog.Int("count", numOfPoliciesInScope),
		)

		for _, policy := range policiesInScope {
			line, err := ConvertPolicy(*policy, ctx)
			if err != nil {
				errs = append(errs, fmt.Errorf("%s: %w", policy.Name, err))
				continue
			}
			results = append(results, line)
		}
	}

	return results, errors.Join(errs...)
}

// ConvertPolicy converts a model.Policy to a Check Point access rule command.
func ConvertPolicy(p model.Policy, ctx *Context) (string, error) {
	var errs []error
	var sb strings.Builder

	sb.WriteString(`add access-rule layer "`)
	sb.WriteString(ctx.App.AppConfig.CheckPoint.Cli.AccessRuleLayer.Value)
	sb.WriteString(`" position.bottom "`)
	sb.WriteString(ctx.App.AppConfig.CheckPoint.Cli.AccessRuleSection.Value)
	sb.WriteString(`" name "`)
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
	if p.Match.NegateSource {
		sb.WriteString(` source-negate true`)
	}
	sb.WriteString(buildIndexedKVWithDefaultAny("source", src))
	if p.Match.NegateDestination {
		sb.WriteString(` destination-negate true`)
	}
	sb.WriteString(buildIndexedKVWithDefaultAny("destination", dst))
	sb.WriteString(buildIndexedKVWithDefaultAny("service", svc))
	sb.WriteString(buildKV("action", normalizeAction(p.Action.Type)))
	if p.Logging.LogAtStart || p.Logging.LogAtEnd {
		sb.WriteString(` track.type "Log"`)
		if p.Action.Type == model.ActionAllow {
			sb.WriteString(` track.accounting true`)
		} else {
			sb.WriteString(` track.accounting false`)
		}
	}
	sb.WriteString(buildComment(p.Description))
	// タグが存在する場合には、最初のタグをカスタム フィールド-1 として設定する
	if len(p.Tags) > 0 {
		// sb.WriteString(buildIndexedKV("tags", p.Tags))
		sb.WriteString(buildKV("custom-fields.field-1", p.Tags[0]))
	}

	return sb.String(), errors.Join(errs...)
}

func normalizeAction(a model.ActionType) string {
	/*
	 * "Accept", "Drop", "Ask", "Inform", "Reject", "User Auth", "Client Auth", "Apply Layer".
	 */
	switch a {
	case model.ActionAllow:
		return "Accept"
	default:
		return "Drop"
	}
}
