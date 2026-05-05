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

func ConvertAddresses(addrs []model.Address, ctx *Context) ([]string, error) {
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
		line, err := ConvertAddress(addr, ctx)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", addr.Name, err))
			continue
		}
		results = append(results, line)
	}

	return results, errors.Join(errs...)
}

func ConvertAddress(addr model.Address, ctx *Context) (string, error) {
	comment := buildComment(addr.Description)
	tags := buildTags(addr.Tags)

	switch addr.Type {

	case model.AddressTypeIPNetmask:
		if ip.IsIPv4(addr.Value) {
			if ip.IsIPv4Host(addr.Value) {
				ipaddr := addr.Value
				ip, _, err := net.ParseCIDR(ipaddr)
				if err == nil {
					ipaddr = ip.String()
				}
				// host
				if _, ok := ctx.AddrMap[addr.Name]; !ok {
					ctx.AddrMap[addr.Name] = addr.Name
				}
				return fmt.Sprintf(
					"add host name \"%s\" ipv4-address %s%s%s",
					addr.Name, ipaddr, comment, tags,
				), nil
			}
			ip, ipnet, err := net.ParseCIDR(addr.Value)
			if err == nil {
				length, _ := ipnet.Mask.Size()
				if _, ok := ctx.AddrMap[addr.Name]; !ok {
					ctx.AddrMap[addr.Name] = addr.Name
				}
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
				if _, ok := ctx.AddrMap[addr.Name]; !ok {
					ctx.AddrMap[addr.Name] = addr.Name
				}
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

	case model.AddressTypeFQDN:
		var comment string
		if addr.Description != "" {
			comment = buildComment(addr.Description)
		} else {
			comment = buildComment(addr.Name)
		}

		if _, ok := ctx.AddrMap[addr.Name]; !ok {
			ctx.AddrMap[addr.Name] = "." + addr.Value
		}
		return fmt.Sprintf(
			"add dns-domain name \".%s\" is-sub-domain false%s%s",
			addr.Value, comment, tags,
		), nil

	default:
		return "", fmt.Errorf("unsupported address type: %s", addr.Type)
	}
}

func ConvertAddressGroups(groups []model.AddressGroup, ctx *Context) ([]string, error) {
	cmds1, err := BuildEmptyGroups(groups, ctx)
	if err != nil {
		return nil, err
	}

	cmds2, err := BuildGroupMembers(groups, ctx)
	if err != nil {
		return nil, err
	}

	cmds1 = append(cmds1, cmds2...)
	return cmds1, nil
}

func BuildEmptyGroups(groups []model.AddressGroup, ctx *Context) ([]string, error) {
	var results []string
	var errs []error

	sort.SliceStable(groups, func(i, j int) bool {
		return strings.ToLower(groups[i].Name) < strings.ToLower(groups[j].Name)
	})

	for _, g := range groups {
		if g.Name == "" {
			errs = append(errs, fmt.Errorf("group name is empty"))
			continue
		}

		comment := buildComment(g.Description)
		tags := buildTags(g.Tags)

		var sb strings.Builder
		sb.WriteString(`add group name "`)
		sb.WriteString(g.Name)
		sb.WriteString(`"`)
		sb.WriteString(comment)
		sb.WriteString(tags)

		results = append(results, sb.String())
		if _, ok := ctx.AddrMap[g.Name]; !ok {
			ctx.AddrMap[g.Name] = g.Name
		}
	}

	return results, errors.Join(errs...)
}

func BuildGroupMembers(groups []model.AddressGroup, ctx *Context) ([]string, error) {
	const chunkSize = 20

	var results []string
	var errs []error

	for _, g := range groups {
		if len(g.Members) == 0 {
			continue
		}

		for i := 0; i < len(g.Members); i += chunkSize {
			end := i + chunkSize
			if end > len(g.Members) {
				end = len(g.Members)
			}

			chunk := g.Members[i:end]
			mapped := make([]string, 0, len(chunk))
			for _, m := range chunk {
				if v, ok := ctx.AddrMap[m]; ok {
					m = v
				} else {
					errs = append(errs, fmt.Errorf("could not find member '%s' to add to group '%s'", m, g.Name))
				}
				mapped = append(mapped, m)
			}

			var sb strings.Builder
			sb.WriteString(`set group name "`)
			sb.WriteString(g.Name)
			sb.WriteString(`"`)
			sb.WriteString(buildIndexedKV("members.add", mapped))

			results = append(results, sb.String())
		}
	}

	return results, errors.Join(errs...)
}
