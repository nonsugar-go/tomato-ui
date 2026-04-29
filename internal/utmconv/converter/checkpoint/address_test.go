package checkpoint

import (
	"strings"
	"testing"

	"github.com/nonsugar-go/tomato-ui/internal/utmconv/model"
)

func TestConvertAddress(t *testing.T) {
	tests := []struct {
		name    string
		address model.Address
		want    string
		wantErr bool
	}{
		{
			name: "IPv4 host",
			address: model.Address{
				Type:        model.AddressTypeIPNetmask,
				Name:        "Host1",
				Value:       "192.168.1.1/32",
				Description: "Host address",
				Tags:        []string{"tag1", "tag2"},
			},
			want:    `add host name "Host1" ipv4-address 192.168.1.1 comments "Host address" tags.1 "tag1" tags.2 "tag2"`,
			wantErr: false,
		},
		{
			name: "IPv4 network",
			address: model.Address{
				Type:        model.AddressTypeIPNetmask,
				Name:        "Net1",
				Value:       "192.168.1.0/24",
				Description: "Net1 address",
				Tags:        []string{"tag1", "tag2"},
			},
			want:    `add network name "Net1" subnet4 192.168.1.0 mask-length4 24 comments "Net1 address" tags.1 "tag1" tags.2 "tag2"`,
			wantErr: false,
		},
		{
			name: "IPv6 host",
			address: model.Address{
				Type:        model.AddressTypeIPNetmask,
				Name:        "Host2",
				Value:       "2001:db8::1/128",
				Description: "IPv6 host address",
				Tags:        []string{"tag1", "tag2"},
			},
			want:    `add host name "Host2" ipv6-address 2001:db8::1 comments "IPv6 host address" tags.1 "tag1" tags.2 "tag2"`,
			wantErr: false,
		},
		{
			name: "IPv6 network",
			address: model.Address{
				Type:        model.AddressTypeIPNetmask,
				Name:        "Net2",
				Value:       "2001:db8::/64",
				Description: "IPv6 network address",
				Tags:        []string{"tag1", "tag2"},
			},
			want:    `add network name "Net2" subnet6 2001:db8:: mask-length6 64 comments "IPv6 network address" tags.1 "tag1" tags.2 "tag2"`,
			wantErr: false,
		},
		{
			name: "Domain name",
			address: model.Address{
				Type:        model.AddressTypeFQDN,
				Name:        "Domain1",
				Value:       "example.com",
				Description: "Domain address",
				Tags:        []string{"tag1", "tag2"},
			},
			want:    `add dns-domain name ".example.com" is-sub-domain false comments "Domain address" tags.1 "tag1" tags.2 "tag2"`,
			wantErr: false,
		},
		{
			name: "Domain name with empty description",
			address: model.Address{
				Type:        model.AddressTypeFQDN,
				Name:        "Domain2",
				Value:       "example.org",
				Description: "",
				Tags:        []string{"tag1", "tag2"},
			},
			want:    `add dns-domain name ".example.org" is-sub-domain false comments "Domain2" tags.1 "tag1" tags.2 "tag2"`,
			wantErr: false,
		},
		{
			name: "Invalid type",
			address: model.Address{
				Type:        model.AddressTypeUnknown,
				Name:        "Invalid1",
				Value:       "value",
				Description: "Invalid address",
				Tags:        []string{"tag1", "tag2"},
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewContext()
			got, err := ConvertAddress(tt.address, ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("convertAddress() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestConvertAddressGroups(t *testing.T) {
	tests := []struct {
		name    string
		groups  []model.AddressGroup
		want    []string
		wantErr bool
	}{
		{
			name: "Single group with members",
			groups: []model.AddressGroup{
				{Name: "Group1", Members: []string{"Member1", "Member2"}, Description: "Test group", Tags: []string{"tag1", "tag2"}},
			},
			want: []string{
				`add group name "Group1" comments "Test group" tags.1 "tag1" tags.2 "tag2"`,
				`set group name "Group1" members.add.1 "Member1" members.add.2 "Member2"`,
			},
			wantErr: false,
		},
		{
			name: "Multiple groups with members",
			groups: []model.AddressGroup{
				{Name: "Group1", Members: []string{"Member1"}, Description: "First group"},
				{Name: "Group2", Members: []string{"Member2", "Member3"}, Description: "Second group", Tags: []string{"tag1"}},
			},
			want: []string{
				`add group name "Group1" comments "First group"`,
				`add group name "Group2" comments "Second group" tags.1 "tag1"`,
				`set group name "Group1" members.add.1 "Member1"`,
				`set group name "Group2" members.add.1 "Member2" members.add.2 "Member3"`,
			},
			wantErr: false,
		},
		{
			name:    "Group with no members",
			groups:  []model.AddressGroup{{Name: "EmptyGroup", Description: "This group has no members"}},
			want:    []string{`add group name "EmptyGroup" comments "This group has no members"`},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewContext()
			got, err := ConvertAddressGroups(tt.groups, ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertAddressGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if strings.Join(got, "\n") != strings.Join(tt.want, "\n") {
				t.Errorf("ConvertAddressGroups() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuildEmptyGroups(t *testing.T) {
	tests := []struct {
		name   string
		groups []model.AddressGroup
		want   []string
	}{
		{
			name:   "No groups",
			groups: []model.AddressGroup{},
			want:   []string{},
		},
		{
			name: "One group",
			groups: []model.AddressGroup{
				{Name: "Group1"},
			},
			want: []string{`add group name "Group1"`},
		},
		{
			name: "Multiple groups",
			groups: []model.AddressGroup{
				{Name: "Group1"},
				{Name: "Group2"},
			},
			want: []string{
				`add group name "Group1"`,
				`add group name "Group2"`,
			},
		},
		{
			name:   "Group with description",
			groups: []model.AddressGroup{{Name: "Group1", Description: "Test group"}},
			want:   []string{`add group name "Group1" comments "Test group"`},
		},
		{
			name:   "Group with tags",
			groups: []model.AddressGroup{{Name: "Group1", Tags: []string{"tag1", "tag2"}}},
			want:   []string{`add group name "Group1" tags.1 "tag1" tags.2 "tag2"`},
		},
		{
			name:   "Group with description and tags",
			groups: []model.AddressGroup{{Name: "Group1", Description: "Test group", Tags: []string{"tag1", "tag2"}}},
			want:   []string{`add group name "Group1" comments "Test group" tags.1 "tag1" tags.2 "tag2"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewContext()
			got, err := BuildEmptyGroups(tt.groups, ctx)
			if err != nil {
				t.Errorf("BuildEmptyGroups() error = %v", err)
				return
			}
			if strings.Join(got, "\n") != strings.Join(tt.want, "\n") {
				t.Errorf("BuildEmptyGroups() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuildGroupMembers(t *testing.T) {
	tests := []struct {
		name   string
		groups []model.AddressGroup
		want   []string
	}{
		{
			name: "One group with members",
			groups: []model.AddressGroup{
				{Name: "Group1", Members: []string{"Member1", "Member2"}},
			},
			want: []string{
				`set group name "Group1" members.add.1 "Member1" members.add.2 "Member2"`,
			},
		},
		{
			name: "Multiple groups with members",
			groups: []model.AddressGroup{
				{Name: "Group1", Members: []string{"Member1"}},
				{Name: "Group2", Members: []string{"Member2", "Member3"}},
			},
			want: []string{
				`set group name "Group1" members.add.1 "Member1"`,
				`set group name "Group2" members.add.1 "Member2" members.add.2 "Member3"`,
			},
		},
		{
			name: "Group with no members",
			groups: []model.AddressGroup{
				{Name: "Group1", Members: []string{}},
			},
			want: []string{},
		},
		{
			name: "Groups with and without members",
			groups: []model.AddressGroup{
				{Name: "Group1", Members: []string{"Member1"}},
				{Name: "Group2", Members: []string{}},
			},
			want: []string{
				`set group name "Group1" members.add.1 "Member1"`,
			},
		},
		{
			name: "Group with 44 members",
			groups: []model.AddressGroup{{Name: "Group1", Members: []string{
				"m1", "m2", "m3", "m4", "m5", "m6", "m7", "m8", "m9", "m10",
				"m11", "m12", "m13", "m14", "m15", "m16", "m17", "m18", "m19", "m20",
				"m21", "m22", "m23", "m24", "m25", "m26", "m27", "m28", "m29", "m30",
				"m31", "m32", "m33", "m34", "m35", "m36", "m37", "m38", "m39", "m40",
				"m41", "m42", "m43", "m44"}}},
			want: []string{
				`set group name "Group1" members.add.1 "m1" members.add.2 "m2" members.add.3 "m3" members.add.4 "m4" members.add.5 "m5" members.add.6 "m6" members.add.7 "m7" members.add.8 "m8" members.add.9 "m9" members.add.10 "m10" members.add.11 "m11" members.add.12 "m12" members.add.13 "m13" members.add.14 "m14" members.add.15 "m15" members.add.16 "m16" members.add.17 "m17" members.add.18 "m18" members.add.19 "m19" members.add.20 "m20"`,
				`set group name "Group1" members.add.1 "m21" members.add.2 "m22" members.add.3 "m23" members.add.4 "m24" members.add.5 "m25" members.add.6 "m26" members.add.7 "m27" members.add.8 "m28" members.add.9 "m29" members.add.10 "m30" members.add.11 "m31" members.add.12 "m32" members.add.13 "m33" members.add.14 "m34" members.add.15 "m35" members.add.16 "m36" members.add.17 "m37" members.add.18 "m38" members.add.19 "m39" members.add.20 "m40"`,
				`set group name "Group1" members.add.1 "m41" members.add.2 "m42" members.add.3 "m43" members.add.4 "m44"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewContext()
			got, err := BuildGroupMembers(tt.groups, ctx)
			if err != nil {
				t.Errorf("BuildGroupMembers() error = %v", err)
				return
			}
			if strings.Join(got, "\n") != strings.Join(tt.want, "\n") {
				t.Errorf("BuildGroupMembers() = %q, want %q", got, tt.want)
			}
		})
	}
}
