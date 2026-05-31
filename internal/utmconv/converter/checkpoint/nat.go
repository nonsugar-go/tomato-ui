package checkpoint

import (
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"github.com/nonsugar-go/tomato-ui/internal/utmconv/model"
)

func GroupNATByScope(natRules []model.NATRule) map[string][]*model.NATRule {
	grouped := make(map[string][]*model.NATRule)

	for i := range natRules {
		n := &natRules[i]
		scope := n.Scope

		if scope == "" {
			scope = "default"
		}

		grouped[scope] = append(grouped[scope], n)
	}

	return grouped
}

// TODO:
func ConvertNATPolicies(natRules []model.NATRule, ctx *Context) ([]string, error) {
	var results []string
	var errs []error

	// TODO: NAT タイプごとにグループ化して、セクションを分ける
	grouped := GroupNATByScope(natRules)
	scopes := make([]string, 0, len(grouped))
	for s := range grouped {
		scopes = append(scopes, s)
	}
	sort.Strings(scopes)

	for _, scope := range scopes {
		natRulesInScope := grouped[scope]
		numOfNATRulesInScope := len(natRulesInScope)

		results = append(results, fmt.Sprintf("#\n# --- %s (%d) ---", scope, numOfNATRulesInScope))

		slog.Info("Check Point の NAT ルールを処理中",
			slog.String("scope", scope),
			slog.Int("count", numOfNATRulesInScope),
		)

		for _, natRule := range natRulesInScope {
			line, err := ConvertNATPolicy(*natRule, ctx)
			if err != nil {
				errs = append(errs, fmt.Errorf("%s: %w", natRule.Name, err))
				continue
			}
			lines := strings.Split(line, "\n")
			results = append(results, lines...)
		}
	}

	return results, errors.Join(errs...)
}

func buildNATRule(
	pkg, section, name string, enabled bool, method, origSrc, origDst, origSvc, trSrc, trDst, trSvc, description string,
	tags []string,
) string {
	var sb strings.Builder

	sb.WriteString(`add nat-rule package "`)
	sb.WriteString(pkg)
	sb.WriteString(`" position.bottom "`)
	sb.WriteString(section)
	sb.WriteString(`" name "`)
	sb.WriteString(name)
	sb.WriteString(`"`)
	if !enabled {
		sb.WriteString(` enabled false`)
	}
	// sb.WriteString(` "install-on "Policy Targets"`)
	sb.WriteString(` method "`)
	sb.WriteString(method)
	sb.WriteString(`"`)
	sb.WriteString(buildKV("original-source", origSrc))
	sb.WriteString(buildKV("original-destination", origDst))
	sb.WriteString(buildKV("original-service", origSvc))
	sb.WriteString(buildKV("translated-source", trSrc))
	sb.WriteString(buildKV("translated-destination", trDst))
	sb.WriteString(buildKV("translated-service", trSvc))
	sb.WriteString(buildComment(description))
	sb.WriteString(buildIndexedKV("tags", tags))

	return sb.String()
}

// ConvertNATPolicy converts a single NAT rule to Check Point CLI format.
// TODO: unchecked
func ConvertNATPolicy(n model.NATRule, ctx *Context) (string, error) {
	mustObjSingle := func(objs []string, objType string, errs *[]error) {
		if len(objs) > 1 {
			slog.Warn("NAT ルールの変換: 複数のオブジェクトが指定されているため、最初のオブジェクトのみを使用します",
				slog.String("obj_type", objType),
				slog.Int("num_objs", len(objs)),
				slog.String("objs", strings.Join(objs, ";")),
			)
		} else if len(objs) != 1 {
			*errs = append(*errs, fmt.Errorf("%s must be a single object, got %d", objType, len(objs)))
		}
	}

	AnyAddrToZone := func(addrStr string, zones []string, errs *[]error) string {
		if addrStr == "Any" {
			if len(zones) == 1 {
				if zones[0] != "any" {
					mappedZone, err := mapStrings(zones, ctx.ZoneMap)
					if err != nil {
						*errs = append(*errs, err)
					} else {
						addrStr = mappedZone[0]
					}
				}
			} else if len(zones) > 1 {
				*errs = append(*errs, fmt.Errorf("zones must be a single zone, got %d", len(zones)))
			}
		}
		return addrStr
	}

	var errs []error
	var sb strings.Builder

	pkg := ctx.App.AppConfig.CheckPoint.Cli.NatRulePackage.Value

	// slog.Info("Check Point のNATルールを処理中",
	// 	slog.String("name", n.Name),
	// 	slog.Bool("enabled", n.Enabled),
	// 	slog.String("from_zone", strings.Join(n.FromZones, ", ")),
	// 	slog.String("to_zone", strings.Join(n.ToZones, ", ")),
	// 	slog.String("original_source", strings.Join(n.OriginalSource, ", ")),
	// 	slog.String("original_destination", strings.Join(n.OriginalDestination, ", ")),
	// 	slog.String("original_service", strings.Join(n.OriginalService, ", ")),
	// 	slog.String("translated_source", strings.Join(n.TranslatedSource, ", ")),
	// 	slog.String("translated_destination", strings.Join(n.TranslatedDestination, ", ")),
	// 	slog.String("translated_service", strings.Join(n.TranslatedService, ", ")),
	// )

	origSrc, err := mapStringsOrDefault(n.OriginalSource, ctx.AddrMap, "Any")
	if err != nil {
		errs = append(errs, err)
	}
	mustObjSingle(origSrc, "original source", &errs)
	origDst, err := mapStringsOrDefault(n.OriginalDestination, ctx.AddrMap, "Any")
	if err != nil {
		errs = append(errs, err)
	}
	mustObjSingle(origDst, "original destination", &errs)
	origSvc, err := mapStringsOrDefault(n.OriginalService, ctx.SvcMap, "Any")
	if err != nil {
		errs = append(errs, err)
	}
	mustObjSingle(origSvc, "original service", &errs)
	trSrc, err := mapStringsOrDefault(n.TranslatedSource, ctx.AddrMap, "Original")
	if err != nil {
		errs = append(errs, err)
	}
	mustObjSingle(trSrc, "translated source", &errs)
	trDst, err := mapStringsOrDefault(n.TranslatedDestination, ctx.AddrMap, "Original")
	if err != nil {
		errs = append(errs, err)
	}
	mustObjSingle(trDst, "translated destination", &errs)
	trSvc, err := mapStringsOrDefault(n.TranslatedService, ctx.SvcMap, "Original")
	if err != nil {
		errs = append(errs, err)
	}
	mustObjSingle(trSvc, "translated service", &errs)

	origSrcStr := AnyAddrToZone(origSrc[0], n.FromZones, &errs)
	origDstStr := AnyAddrToZone(origDst[0], n.ToZones, &errs)
	origSvcStr := origSvc[0]
	trSrcStr := trSrc[0]
	trDstStr := trDst[0]
	trSvcStr := trSvc[0]

	switch n.Type {
	case model.NATTypeStatic:
		if n.BiDirectional {
			// 双方向の静的 NAT ルールは、元のルールと逆のルールを追加するが、サービスは Any / Original として扱う
			sb.WriteString(buildNATRule(
				pkg, "Static NAT Rules", n.Name+"_In", n.Enabled, "static",
				origDstStr, trSrcStr, "Any", "Original", origSrcStr, "Original", n.Description, n.Tags,
			))
			sb.WriteString("\n")
		}
		sb.WriteString(buildNATRule(
			pkg, "Static NAT Rules", n.Name+"_Out", n.Enabled, "static",
			origSrcStr, origDstStr, origSvcStr, trSrcStr, trDstStr, trSvcStr, n.Description, n.Tags,
		))
	case model.NATTypeHide:
		// TODO:
		sb.WriteString(buildNATRule(
			pkg, "Hide NAT Rules", n.Name, n.Enabled, "hide",
			origSrcStr, origDstStr, origSvcStr, trSrcStr, trDstStr, trSvcStr, n.Description, n.Tags,
		))
	default:
		errs = append(errs, fmt.Errorf("unsupported NAT type: %s", n.Type))
	}
	return sb.String(), errors.Join(errs...)
}
