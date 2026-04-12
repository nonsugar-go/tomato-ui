package ip

import "net/netip"

func IsIPv6(s string) bool {
	if prefix, err := netip.ParsePrefix(s); err == nil {
		return prefix.Addr().Is6()
	}
	ip, err := netip.ParseAddr(s)
	return err == nil && ip.Is6()
}

func IsIPv6Host(s string) bool {
	if prefix, err := netip.ParsePrefix(s); err == nil {
		return prefix.Addr().Is6() && prefix.Bits() == 128
	}
	ip, err := netip.ParseAddr(s)
	return err == nil && ip.Is6()
}
