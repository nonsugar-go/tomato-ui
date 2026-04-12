package checkpoint

import (
	"errors"
	"fmt"
	"net"
	"sort"
	"strings"

	"github.com/nonsugar-go/tomato-ui/internal/utmconv/ip"
	"github.com/nonsugar-go/tomato-ui/internal/utmconv/model"
)

func ConvertAddresses(addrs []model.Address) ([]string, error) {
	var results []string
	var errs []error

	sort.SliceStable(addrs, func(i, j int) bool {
		a, b := addrs[i], addrs[j]
		if a.Type != b.Type {
			return a.Type < b.Type
		}
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})

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
	tags := buildTags(addr.Tags)

	switch addr.Type {

	case "ip-netmask":
		if ip.IsIPv4(addr.Value) {
			if ip.IsIPv4Host(addr.Value) {
				ipaddr := addr.Value
				ip, _, err := net.ParseCIDR(ipaddr)
				if err == nil {
					ipaddr = ip.String()
				}
				// host
				return fmt.Sprintf(
					"add host name \"%s\" ipv4-address %s%s%s",
					addr.Name, ipaddr, comment, tags,
				), nil
			}
			ip, ipnet, err := net.ParseCIDR(addr.Value)
			if err == nil {
				length, _ := ipnet.Mask.Size()
				return fmt.Sprintf(
					"add network name \"%s\" subnet4 %s mask-length4 %d%s%s",
					addr.Name, ip.String(), length, comment, tags,
				), nil
			} else {
				return "", fmt.Errorf("invalid IP address: %s", addr.Value)
			}
		} else if ip.IsIPv6(addr.Value) {
			if ip.IsIPv6Host(addr.Value) { // NOTE: untested
				ipaddr := addr.Value
				ip, _, err := net.ParseCIDR(ipaddr)
				if err == nil {
					ipaddr = ip.String()
				}
				// host
				return fmt.Sprintf(
					"add host name \"%s\" ipv6-address %s%s%s",
					addr.Name, ipaddr, comment, tags,
				), nil
			}
			ip, ipnet, err := net.ParseCIDR(addr.Value)
			if err == nil {
				length, _ := ipnet.Mask.Size()
				return fmt.Sprintf(
					"add network name \"%s\" subnet6 %s mask-length6 %d%s%s",
					addr.Name, ip.String(), length, comment, tags,
				), nil
			} else {
				return "", fmt.Errorf("invalid IP address: %s", addr.Value)
			}
		} else {
			return "", fmt.Errorf("invalid IP address: %s", addr.Value)
		}

	case "fqdn":
		return fmt.Sprintf(
			"add dns-domain name \".%s\" is-sub-domain false%s%s",
			addr.Value, comment, tags,
		), nil

	default:
		return "", fmt.Errorf("unsupported address type: %s", addr.Type)
	}
}

func ConvertAddressGroups(groups []model.AddressGroup) ([]string, error) {
	var results []string
	var errs []error

	sort.SliceStable(groups, func(i, j int) bool {
		return strings.ToLower(groups[i].Name) < strings.ToLower(groups[j].Name)
	})

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
