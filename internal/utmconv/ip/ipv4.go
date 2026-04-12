package ip

import "net/netip"

func IsIPv4(s string) bool {
	if prefix, err := netip.ParsePrefix(s); err == nil {
		return prefix.Addr().Is4()
	}
	ip, err := netip.ParseAddr(s)
	return err == nil && ip.Is4()
}

func IsIPv4Host(s string) bool {
	if prefix, err := netip.ParsePrefix(s); err == nil {
		return prefix.Addr().Is4() && prefix.Bits() == 32
	}
	ip, err := netip.ParseAddr(s)
	return err == nil && ip.Is4()
}
