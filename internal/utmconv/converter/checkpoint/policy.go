package checkpoint

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nonsugar-go/tomato-ui/internal/utmconv/model"
)

func ConvertPolicies(policies []model.Policy, ctx *Context) ([]string, error) {
	var results []string
	var errs []error

	for _, policy := range policies {
		line, err := ConvertPolicy(policy, ctx)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", policy.Name, err))
			continue
		}
		results = append(results, line)
	}

	return results, errors.Join(errs...)
}

// ConvertPolicy converts a model.Policy to a Check Point access rule command.
// NOTE: untested
func ConvertPolicy(p model.Policy, ctx *Context) (string, error) {
	var sb strings.Builder
	layer := "Network" // NOTE: untested

	sb.WriteString(`add access-rule layer "`)
	sb.WriteString(layer)
	// sb.WriteString(`" position 1`)
	sb.WriteString(`" position bottom`)
	sb.WriteString(` name "`)
	sb.WriteString(p.Name)
	sb.WriteString(`"`)
	if p.Enabled == false {
		sb.WriteString(` enabled false`)
	}
	src := mapStrings(p.Match.Sources.Names(), ctx.AddrMap)
	dst := mapStrings(p.Match.Destinations.Names(), ctx.AddrMap)
	if p.Match.NegateSource {
		sb.WriteString(` source-negate true`)
	}
	sb.WriteString(buildIndexedKV("source", src))
	if p.Match.NegateDestination {
		sb.WriteString(` destination-negate true`)
	}
	sb.WriteString(buildIndexedKV("destination", dst))
	sb.WriteString(buildIndexedKV("service", p.Match.Services.Names()))
	action := normalizeAction(string(p.Action.Type))
	sb.WriteString(buildKV("action", action))
	sb.WriteString(buildComment(p.Description))
	sb.WriteString(buildIndexedKV("tags", p.Tags))

	return sb.String(), nil
	// return "", fmt.Errorf("invalid policy: %s", p.Name)
}

func normalizeAction(a string) string {
	/*
	 * "Accept", "Drop", "Ask", "Inform", "Reject", "User Auth", "Client Auth", "Apply Layer".
	 */
	switch strings.ToLower(a) {
	case "accept", "allow", "permit":
		return "Accept"
	case "drop":
		return "Drop"
	case "deny", "reject":
		return "Reject"
	default:
		return "Accept"
	}
}
