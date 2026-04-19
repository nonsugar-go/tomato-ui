package checkpoint

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/nonsugar-go/tomato-ui/internal/utmconv/model"
)

func ConvertServices(svcs []model.Service) ([]string, error) {
	var results []string
	var errs []error

	sort.SliceStable(svcs, func(i, j int) bool {
		return strings.ToLower(svcs[i].Name) < strings.ToLower(svcs[j].Name)
	})

	for _, svc := range svcs {
		line, err := ConvertService(svc)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", svc.Name, err))
			continue
		}
		results = append(results, line)
	}

	return results, errors.Join(errs...)
}

func ConvertService(svc model.Service) (string, error) {
	sourcePort := ""
	if svc.SourcePorts != "" {
		sourcePort = fmt.Sprintf(" source-port %s", svc.SourcePorts)
	}
	comment := buildComment(svc.Description)
	tags := buildTags(svc.Tags)

	switch svc.Type {

	case model.ServiceTypeTCP:
		return fmt.Sprintf(
			"add service-tcp name \"%s\" port %s%s%s%s",
			svc.Name, svc.Ports, sourcePort, comment, tags,
		), nil

	case model.ServiceTypeUDP:
		return fmt.Sprintf(
			"add service-udp name \"%s\" port %s%s%s%s",
			svc.Name, svc.Ports, sourcePort, comment, tags,
		), nil

	default:
		return "", fmt.Errorf("unsupported service type: %s", svc.Type)
	}
}

func ConvertServiceGroups(groups []model.ServiceGroup) ([]string, error) {
	var results []string
	var errs []error

	sort.SliceStable(groups, func(i, j int) bool {
		return strings.ToLower(groups[i].Name) < strings.ToLower(groups[j].Name)
	})

	for _, g := range groups {
		line, err := ConvertServiceGroup(g)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", g.Name, err))
			continue
		}
		results = append(results, line)
	}

	return results, errors.Join(errs...)
}

func ConvertServiceGroup(g model.ServiceGroup) (string, error) {
	if len(g.Members) == 0 {
		return "", fmt.Errorf("group has no members")
	}

	comment := buildComment(g.Description)
	tags := buildTags(g.Tags)

	var sb strings.Builder
	sb.WriteString(`add service-group name "`)
	sb.WriteString(g.Name)
	sb.WriteString(`"`)
	sb.WriteString(buildIndexedKV("members", g.Members))
	sb.WriteString(comment)
	sb.WriteString(tags)

	return sb.String(), nil
}
