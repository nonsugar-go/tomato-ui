package checkpoint

import (
	"testing"

	"github.com/nonsugar-go/tomato-ui/internal/utmconv/model"
)

func TestNormalizeAction(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"Accept", "Accept", "Accept"},
		{"Allow", "Allow", "Accept"},
		{"Permit", "Permit", "Accept"},
		{"Drop", "Drop", "Drop"},
		{"Deny", "Deny", "Reject"},
		{"Reject", "Reject", "Reject"},
		{"Unknown", "Unknown", "Accept"}, // default case
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeAction(tt.input); got != tt.want {
				t.Errorf("normalizeAction(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestConvertPolicy(t *testing.T) {
	tests := []struct {
		name    string
		policy  model.Policy
		want    string
		wantErr bool
	}{
		{
			name: "Valid policy",
			policy: model.Policy{
				Name:        "Test Policy",
				Enabled:     true,
				Description: "This is a test policy",
				Action:      model.PolicyAction{Type: model.ActionAllow},
			},
			want: `add access-rule layer "Network" position bottom name "Test Policy" action "Accept" comments "This is a test policy"`,
		},
		{
			name: "Policy with sources and destinations",
			policy: model.Policy{
				Name:    "Policy with Sources and Destinations",
				Enabled: true,
				Match: model.PolicyMatch{
					Sources:      model.AddressRefs{{Name: "Source1"}, {Name: "Source2"}},
					Destinations: model.AddressRefs{{Name: "Destination1"}},
				},
				Action: model.PolicyAction{Type: model.ActionDeny},
			},
			want: `add access-rule layer "Network" position bottom name "Policy with Sources and Destinations" source.1 "Source1" source.2 "Source2" destination.1 "Destination1" action "Reject"`,
		},
		{
			name: "Policy with tags",
			policy: model.Policy{
				Name:    "Policy with Tags",
				Enabled: true,
				Tags:    []string{"tag1", "tag2"},
				Action:  model.PolicyAction{Type: model.ActionAllow},
			},
			want: `add access-rule layer "Network" position bottom name "Policy with Tags" action "Accept" tags.1 "tag1" tags.2 "tag2"`,
		},
		{
			name: "Policy with enabled flag",
			policy: model.Policy{
				Name:    "Policy with Enabled Flag",
				Enabled: false,
				Action:  model.PolicyAction{Type: model.ActionAllow},
			},
			want: `add access-rule layer "Network" position bottom name "Policy with Enabled Flag" enabled false action "Accept"`,
		},
		{
			name: "Policy with services",
			policy: model.Policy{
				Name:    "Policy with Services",
				Enabled: true,
				Match: model.PolicyMatch{
					Services: model.ServiceRefs{{Name: "Service1"}, {Name: "Service2"}},
				},
				Action: model.PolicyAction{Type: model.ActionAllow},
			},
			want: `add access-rule layer "Network" position bottom name "Policy with Services" service.1 "Service1" service.2 "Service2" action "Accept"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertPolicy(tt.policy)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConvertPolicy() = %q, want %q", got, tt.want)
			}
		})
	}
}
