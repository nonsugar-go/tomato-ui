package ip

import "testing"

func TestIsIPv4(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"192.168.0.1", true},
		{"0.0.0.0", true},
		{"255.255.255.255", true},
		{"192.168.0.0/24", true},
		{"192.168.168.168/33", false},
		{"10.0.0.0/8", true},

		{"::1", false},
		{"2001:db8::/32", false},
		{"example.com", false},
		{"", false},
		{"999.999.999.999", false},
	}

	for _, tt := range tests {
		got := IsIPv4(tt.input)
		if got != tt.want {
			t.Errorf("IsIPv4(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestIsIPv6(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"::1", true},
		{"2001:db8::1", true},
		{"fe80::1", true},
		{"::", true},

		{"2001:db8::/32", true},
		{"2001:db8::1/129", false},
		{"fe80::/64", true},

		{"192.168.0.1", false},
		{"10.0.0.0/8", false},
		{"192.168.0.0/24", false},

		{"example.com", false},
		{"", false},
		{"999.999.999.999", false},
	}

	for _, tt := range tests {
		got := IsIPv6(tt.input)
		if got != tt.want {
			t.Errorf("IsIPv6(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestIsIPv4Host(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"192.168.0.1", true},
		{"0.0.0.0", true},
		{"255.255.255.255", true},

		{"192.168.0.0/24", false},
		{"10.0.0.0/8", false},

		{"::1", false},
		{"2001:db8::/32", false},

		{"example.com", false},
		{"", false},
		{"999.999.999.999", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := IsIPv4Host(tt.input)
			if got != tt.want {
				t.Errorf("IsIPv4Host(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsIPv6Host(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"::1", true},
		{"2001:db8::1", true},
		{"fe80::1", true},
		{"::", true},

		{"::1/128", true},
		{"2001:db8::1/128", true},

		{"2001:db8::/32", false},
		{"fe80::/64", false},

		{"192.168.0.1", false},
		{"10.0.0.0/8", false},

		{"example.com", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := IsIPv6Host(tt.input)
			if got != tt.want {
				t.Errorf("IsIPv6Host(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
