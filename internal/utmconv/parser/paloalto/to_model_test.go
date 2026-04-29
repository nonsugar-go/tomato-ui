package paloalto

import (
	"reflect"
	"testing"

	"github.com/nonsugar-go/tomato-ui/internal/utmconv/model"
)

func TestToModelTags(t *testing.T) {
	tests := []struct {
		name    string
		in      []ScopedTagObject
		want    []model.Tag
		wantErr bool
	}{
		{
			name: "Valid tags",
			in: []ScopedTagObject{
				{Scope: "shared", TagObject: TagObject{Name: "tag1", Color: "color15", Comments: "This is tag1"}},
				{Scope: "shared", TagObject: TagObject{Name: "tag2", Color: "color2", Comments: "This is tag2"}},
			},
			want: []model.Tag{
				{Value: "tag1", Color: "color15", Description: "This is tag1"},
				{Value: "tag2", Color: "color2", Description: "This is tag2"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToModelTags(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ToModelTags() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ToModelTags() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToModelAddresses(t *testing.T) {
	tests := []struct {
		name    string
		in      []ScopedAddress
		want    []model.Address
		wantErr bool
	}{
		{
			name: "Valid addresses",
			in: []ScopedAddress{
				{Scope: "shared", Address: Address{Name: "addr1", IPNetmask: "192.168.1.1/24", Description: "This is addr1"}},
				{Scope: "shared", Address: Address{Name: "addr2", FQDN: "example.com", Description: "This is addr2"}},
			},
			want: []model.Address{
				{Name: "addr1", Type: model.AddressTypeIPNetmask, Value: "192.168.1.1/24", Description: "This is addr1"},
				{Name: "addr2", Type: model.AddressTypeFQDN, Value: "example.com", Description: "This is addr2"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToModelAddresses(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ToModelAddresses() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ToModelAddresses() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToModelServices(t *testing.T) {
	tests := []struct {
		name    string
		in      []ScopedService
		want    []model.Service
		wantErr bool
	}{
		{
			name: "Valid services",
			in: []ScopedService{
				{Scope: "shared", Service: Service{Name: "svc1", Protocol: Protocol{TCP: &TCP{Port: "80"}}, Description: "This is svc1"}},
				{Scope: "shared", Service: Service{Name: "svc2", Protocol: Protocol{UDP: &UDP{Port: "53"}}, Description: "This is svc2"}},
			},
			want: []model.Service{
				{Name: "service-http", Type: model.ServiceTypeTCP, Ports: "80"},
				{Name: "service-https", Type: model.ServiceTypeTCP, Ports: "443"},
				{Name: "svc1", Type: model.ServiceTypeTCP, Ports: "80", Description: "This is svc1"},
				{Name: "svc2", Type: model.ServiceTypeUDP, Ports: "53", Description: "This is svc2"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToModelServices(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ToModelServices() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ToModelServices() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToModelServiceGroups(t *testing.T) {
	tests := []struct {
		name    string
		in      []ScopedServiceGroup
		want    []model.ServiceGroup
		wantErr bool
	}{
		{
			name: "Valid service groups",
			in: []ScopedServiceGroup{
				{Scope: "shared", Group: ServiceGroup{Name: "sg1", Members: []string{"svc1", "svc2"}, Description: "This is sg1", Tags: []string{"tag1"}}},
			},
			want: []model.ServiceGroup{
				{Name: "sg1", Members: []string{"svc1", "svc2"}, Description: "This is sg1", Tags: []string{"tag1"}},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToModelServiceGroups(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ToModelServiceGroups() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ToModelServiceGroups() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToAddrRefs(t *testing.T) {
	in := []string{" host1 ", "any", "", "host1"}

	got := toAddrRefs(in)

	want := []model.AddressRef{
		{Name: "host1"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestToSvcRefs(t *testing.T) {
	tests := []struct {
		name  string
		names []string
		apps  []string
		want  []model.ServiceRef
	}{
		{
			name:  "Valid service refs",
			names: []string{" smtp ", "any", "http", "svc1", "https"},
			apps:  []string{""},
			want: []model.ServiceRef{
				{Name: "http"},
				{Name: "https"},
				{Name: "smtp"},
				{Name: "svc1"},
			},
		},
		{
			name:  "Valid service refs with application-default",
			names: []string{" application-default "},
			apps:  []string{" traceroute ", "icmp"},
			want: []model.ServiceRef{
				{Name: "icmp-proto"},
				{Name: "traceroute"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toSvcRefs(tt.names, tt.apps)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got=%v want=%v", got, tt.want)
			}
		})
	}
}

func TestToAction(t *testing.T) {
	tests := []struct {
		in   string
		want model.ActionType
	}{
		{"allow", model.ActionAllow},
		{"deny", model.ActionDeny},
		{"drop", model.ActionDrop},
		{"reset-client", model.ActionReset},
		{"RESET-SERVER", model.ActionReset},
		{"unknown", model.ActionDeny},
	}

	for _, tt := range tests {
		got := toAction(tt.in)
		if got != tt.want {
			t.Fatalf("in=%s got=%v want=%v", tt.in, got, tt.want)
		}
	}
}

func TestExtractProfiles(t *testing.T) {
	ps := &ProfileSetting{
		Group: []string{"default"},
	}

	got := extractProfiles(ps)
	want := []string{"default"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v want=%v", got, want)
	}
}
