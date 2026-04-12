package checkpoint

import (
	"errors"
	"fmt"
	"net"

	"github.com/nonsugar-go/tomato-ui/internal/utmconv/model"
)

func ConvertAddresses(addrs []model.Address) ([]string, error) {
	var results []string
	var errs []error

	for _, addr := range addrs {
		line, err := ConvertAddress(addr)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", addr.Name, err))
			continue
		}
		results = append(results, line)
	}

	return results, errors.Join(errs...)
}

func ConvertAddress(addr model.Address) (string, error) {
	comment := buildComment(addr.Description)

	switch addr.Type {

	case "ip-netmask":
		ip, ipnet, err := net.ParseCIDR(addr.Value)
		if err == nil {
			ones, _ := ipnet.Mask.Size()
			if ones == 32 {
				// host
				return fmt.Sprintf(
					"add host name \"%s\" ip-address %s%s",
					addr.Name, ip.String(), comment,
				), nil
			}
			return fmt.Sprintf(
				"add network name \"%s\" subnet %s mask-length %d%s",
				addr.Name, ip.String(), ones, comment,
			), nil
		} else {
			// host
			return fmt.Sprintf(
				"add host name \"%s\" ip-address %s%s",
				addr.Name, addr.Value, comment,
			), nil
		}

	case "fqdn":
		return fmt.Sprintf(
			"add dns-domain name \".%s\" is-sub-domain false %s",
			addr.Value, comment,
		), nil

	default:
		return "", fmt.Errorf("unsupported address type: %s", addr.Type)
	}
}

func ConvertAddressGroups(groups []model.AddressGroup) ([]string, error) {
	var results []string
	var errs []error

	for _, g := range groups {
		line, err := ConvertAddressGroup(g)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", g.Name, err))
			continue
		}
		results = append(results, line)
	}

	return results, errors.Join(errs...)
}

func ConvertAddressGroup(g model.AddressGroup) (string, error) {
	if len(g.Members) == 0 {
		return "", fmt.Errorf("group has no members")
	}

	comment := buildComment(g.Description)

	cmd := fmt.Sprintf("add group name \"%s\"", g.Name)

	if len(g.Members) == 1 {
		cmd += fmt.Sprintf(" members \"%s\"", g.Members[0])
	} else {
		for i, m := range g.Members {
			cmd += fmt.Sprintf(" members.%d \"%s\"", i+1, m)
		}
	}

	cmd += comment

	return cmd, nil
}
