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
				// {Name: "service-http", Type: model.ServiceTypeTCP, Ports: "80"},
				// {Name: "service-https", Type: model.ServiceTypeTCP, Ports: "443"},
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
	app := model.App{AppConfig: *model.NewDefaultAppConfig()}

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
				{Name: "ICMP Protocol"},
				{Name: "traceroute"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toSvcRefs(tt.names, tt.apps, app.AppConfig.PaloAlto.Conf.ApplicationDefaultReplacementMap.Value)
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

func TestGetUniqueTags(t *testing.T) {
	tests := []struct {
		name     string
		groupTag string
		tags     []string
		want     []string
	}{
		{
			name:     "正常系：重複なし",
			groupTag: "GroupA",
			tags:     []string{"Tag1", "Tag2"},
			want:     []string{"GroupA", "Tag1", "Tag2"},
		},
		{
			name:     "重複あり：GroupTagがTagsの中にもある",
			groupTag: "Tag1",
			tags:     []string{"Tag1", "Tag2"},
			want:     []string{"Tag1", "Tag2"},
		},
		{
			name:     "重複あり：Tagsの中で重複している",
			groupTag: "GroupA",
			tags:     []string{"Tag1", "Tag1", "Tag2"},
			want:     []string{"GroupA", "Tag1", "Tag2"},
		},
		{
			name:     "境界値：GroupTagが空文字",
			groupTag: "",
			tags:     []string{"Tag1", "Tag2"},
			want:     []string{"Tag1", "Tag2"},
		},
		{
			name:     "境界値：Tagsが空またはnil",
			groupTag: "GroupA",
			tags:     nil,
			want:     []string{"GroupA"},
		},
		{
			name:     "境界値：すべて空",
			groupTag: "",
			tags:     []string{},
			want:     nil,
		},
		{
			name:     "順番の維持：入力された順序で返るか",
			groupTag: "Z",
			tags:     []string{"A", "B", "C"},
			want:     []string{"Z", "A", "B", "C"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getUniqueTags(tt.groupTag, tt.tags)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getUniqueTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToModelPolicies(t *testing.T) {
	app := model.App{AppConfig: *model.NewDefaultAppConfig()}

	tests := []struct {
		name    string
		in      []ScopedSecurity
		want    []model.Policy
		wantErr bool
	}{
		{
			name: "Valid policies 1",
			in: []ScopedSecurity{
				{Scope: "shared", SecurityRule: SecurityRule{
					Name: "rhel 1", Tags: []string{"linux", "rhel"}, Sources: []string{"rhel_1"}, Destinations: []string{"any"},
					Applications: []string{"any"}, Services: []string{"http", "https"}, Action: "allow",
					Description: "Redhat",
				}},
			},
			want: []model.Policy{
				{Name: "rhel 1", Description: "Redhat", Enabled: true,
					Match: model.PolicyMatch{
						Sources:      model.AddressRefs{{Name: "rhel_1"}},
						Destinations: model.AddressRefs{},
						Applications: []string{"any"},
						Services:     model.ServiceRefs{{Name: "http"}, {Name: "https"}},
					},
					Action:  model.PolicyAction{Type: model.ActionAllow},
					Logging: model.Logging{LogAtStart: false, LogAtEnd: true},
					Tags:    []string{"linux", "rhel"},
					Scope:   "shared",
				},
			},
			wantErr: false,
		},
		{
			name: "LogStart=yes, LogEnd=default",
			in: []ScopedSecurity{
				{Scope: "dg1", Rulebase: "pre", SecurityRule: SecurityRule{
					Name: "Log Start yes", Action: "allow", LogStart: "yes", LogSetting: "log_prof1"}},
			},
			want: []model.Policy{
				{Name: "Log Start yes", Enabled: true,
					Action:  model.PolicyAction{Type: model.ActionAllow},
					Logging: model.Logging{LogAtStart: true, LogAtEnd: true, LogProfile: "log_prof1"},
					Scope:   "dg1-pre",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToModelPolicies(tt.in, app.AppConfig.PaloAlto.Conf.ApplicationDefaultReplacementMap.Value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ToModelPolicies() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ToModelPolicies() \ngot = %#v,\nwant %#v", got, tt.want)
			}
		})
	}
}
