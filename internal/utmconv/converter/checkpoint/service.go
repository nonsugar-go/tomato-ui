package checkpoint

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/nonsugar-go/tomato-ui/internal/utmconv/model"
)

func ConvertServices(svcs []model.Service, ctx *Context) ([]string, error) {
	var results []string
	var errs []error

	customSvcs := make([]model.Service, 0, len(svcs))
	for _, svc := range svcs {
		if _, ok := ctx.SvcMap[svc.Name]; !ok {
			customSvcs = append(customSvcs, svc)
		}
	}

	sort.SliceStable(customSvcs, func(i, j int) bool {
		return strings.ToLower(customSvcs[i].Name) < strings.ToLower(customSvcs[j].Name)
	})

	for _, svc := range customSvcs {
		line, err := ConvertService(svc, ctx)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", svc.Name, err))
			continue
		}
		results = append(results, line)
	}

	return results, errors.Join(errs...)
}

func ConvertService(svc model.Service, ctx *Context) (string, error) {
	sourcePort := ""
	if svc.SourcePorts != "" {
		sourcePort = fmt.Sprintf(" source-port %s", svc.SourcePorts)
	}
	comment := buildComment(svc.Description)
	tags := buildTags(svc.Tags)

	switch svc.Type {

	case model.ServiceTypeTCP:
		if _, ok := ctx.SvcMap[svc.Name]; !ok {
			ctx.SvcMap[svc.Name] = svc.Name
		}
		return fmt.Sprintf(
			"add service-tcp name \"%s\" port %s%s%s%s",
			svc.Name, svc.Ports, sourcePort, comment, tags,
		), nil

	case model.ServiceTypeUDP:
		if _, ok := ctx.SvcMap[svc.Name]; !ok {
			ctx.SvcMap[svc.Name] = svc.Name
		}
		return fmt.Sprintf(
			"add service-udp name \"%s\" port %s%s%s%s",
			svc.Name, svc.Ports, sourcePort, comment, tags,
		), nil

	default:
		return "", fmt.Errorf("unsupported service type: %s", svc.Type)
	}
}

func ConvertServiceGroups(groups []model.ServiceGroup, ctx *Context) ([]string, error) {
	var results []string
	var errs []error

	customGroups := make([]model.ServiceGroup, 0, len(groups))
	for _, group := range groups {
		if _, ok := ctx.SvcMap[group.Name]; !ok {
			customGroups = append(customGroups, group)
		}
	}

	sort.SliceStable(customGroups, func(i, j int) bool {
		return strings.ToLower(customGroups[i].Name) < strings.ToLower(customGroups[j].Name)
	})

	for _, g := range customGroups {
		line, err := ConvertServiceGroup(g, ctx)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", g.Name, err))
			continue
		}
		results = append(results, line)
	}

	return results, errors.Join(errs...)
}

func ConvertServiceGroup(g model.ServiceGroup, ctx *Context) (string, error) {
	var errs []error

	if len(g.Members) == 0 {
		return "", fmt.Errorf("group has no members")
	}

	comment := buildComment(g.Description)
	tags := buildTags(g.Tags)

	mapped := make([]string, 0, len(g.Members))
	for _, m := range g.Members {
		if v, ok := ctx.SvcMap[m]; ok {
			m = v
		} else {
			errs = append(errs, fmt.Errorf("could not find member '%s' to add to group '%s'", m, g.Name))
		}
		mapped = append(mapped, m)
	}

	var sb strings.Builder
	sb.WriteString(`add service-group name "`)
	sb.WriteString(g.Name)
	sb.WriteString(`"`)
	sb.WriteString(buildIndexedKV("members", mapped))
	sb.WriteString(comment)
	sb.WriteString(tags)

	if _, ok := ctx.SvcMap[g.Name]; !ok {
		ctx.SvcMap[g.Name] = g.Name
	}

	return sb.String(), errors.Join(errs...)
}
