package checkpoint

import (
	"fmt"
	"testing"

	"github.com/nonsugar-go/tomato-ui/internal/utmconv/model"
)

func TestNormalizeAction(t *testing.T) {
	tests := []struct {
		name  string
		input model.ActionType
		want  string
	}{
		{"Allow", model.ActionAllow, "Accept"},
		{"Deny", model.ActionDeny, "Drop"},
		{"Drop", model.ActionDrop, "Drop"},
		{"Reset", model.ActionReset, "Drop"},
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
	app := model.App{AppConfig: *model.NewDefaultAppConfig()}
	ctx := NewContext(&app)
	ctx.AddrMap = make(map[string]string)
	for i := 1; i <= 2; i++ {
		member := fmt.Sprintf("Source%d", i)
		ctx.AddrMap[member] = member
		member = fmt.Sprintf("Destination%d", i)
		ctx.AddrMap[member] = member
		member = fmt.Sprintf("Service%d", i)
		ctx.SvcMap[member] = member
	}

	tests := []struct {
		name    string
		app     model.App
		policy  model.Policy
		want    string
		wantErr bool
	}{
		{
			name: "Valid policy",
			app:  app,
			policy: model.Policy{
				Name:        "Test Policy",
				Enabled:     true,
				Description: "This is a test policy",
				Action:      model.PolicyAction{Type: model.ActionAllow},
			},
			want: `add access-rule layer "Network" position.bottom "New rules" name "Test Policy" source.1 "Any" destination.1 "Any" service.1 "Any" action "Accept" comments "This is a test policy"`,
		},
		{
			name: "Policy with sources and destinations",
			app:  app,
			policy: model.Policy{
				Name:    "Policy with Sources and Destinations",
				Enabled: true,
				Match: model.PolicyMatch{
					Sources:      model.AddressRefs{{Name: "Source1"}, {Name: "Source2"}},
					Destinations: model.AddressRefs{{Name: "Destination1"}},
				},
				Action:  model.PolicyAction{Type: model.ActionDeny},
				Logging: model.Logging{LogAtStart: true, LogAtEnd: true},
			},
			want: `add access-rule layer "Network" position.bottom "New rules" name "Policy with Sources and Destinations" source.1 "Source1" source.2 "Source2" destination.1 "Destination1" service.1 "Any" action "Drop" track.type "Log"`,
		},
		{
			name: "Policy with tags",
			app:  app,
			policy: model.Policy{
				Name:    "Policy with Tags",
				Enabled: true,
				Tags:    []string{"tag1", "tag2"},
				Action:  model.PolicyAction{Type: model.ActionAllow},
			},
			want: `add access-rule layer "Network" position.bottom "New rules" name "Policy with Tags" source.1 "Any" destination.1 "Any" service.1 "Any" action "Accept" custom-fields.field-1 "tag1"`,
		},
		{
			name: "Policy with enabled Flag and LogAtEnd",
			app:  app,
			policy: model.Policy{
				Name:    "Policy with Enabled Flag and LogAtEnd",
				Enabled: false,
				Tags:    []string{"tag1", "tag2"},
				Action:  model.PolicyAction{Type: model.ActionAllow},
				Logging: model.Logging{LogAtEnd: true},
			},
			want: `add access-rule layer "Network" position.bottom "New rules" name "Policy with Enabled Flag and LogAtEnd" enabled false source.1 "Any" destination.1 "Any" service.1 "Any" action "Accept" track.type "Log" track.accounting true custom-fields.field-1 "tag1"`,
		},
		{
			name: "Policy with services and LogAtStart",
			app:  app,
			policy: model.Policy{
				Name:    "Policy with Services",
				Enabled: true,
				Match: model.PolicyMatch{
					Services: model.ServiceRefs{{Name: "Service1"}, {Name: "Service2"}},
				},
				Action:  model.PolicyAction{Type: model.ActionAllow},
				Logging: model.Logging{LogAtStart: true},
			},
			want: `add access-rule layer "Network" position.bottom "New rules" name "Policy with Services" source.1 "Any" destination.1 "Any" service.1 "Service1" service.2 "Service2" action "Accept" track.type "Log" track.accounting true`,
		},
		{
			name: "Policy with source-negate flag",
			app:  app,
			policy: model.Policy{
				Name:    "Policy with Source-Negate Flag",
				Enabled: true,
				Match: model.PolicyMatch{
					Sources:      model.AddressRefs{{Name: "Source1"}, {Name: "Source2"}},
					Destinations: model.AddressRefs{{Name: "Destination1"}},
					NegateSource: true,
				},
				Action: model.PolicyAction{Type: model.ActionDeny},
			},
			want: `add access-rule layer "Network" position.bottom "New rules" name "Policy with Source-Negate Flag" source-negate true source.1 "Source1" source.2 "Source2" destination.1 "Destination1" service.1 "Any" action "Drop"`,
		},
		{
			name: "Policy with destination-negate flag",
			app:  app,
			policy: model.Policy{
				Name:    "Policy with Destination-Negate Flag",
				Enabled: true,
				Match: model.PolicyMatch{
					Sources:           model.AddressRefs{{Name: "Source1"}},
					Destinations:      model.AddressRefs{{Name: "Destination1"}, {Name: "Destination2"}},
					NegateDestination: true,
				},
				Action: model.PolicyAction{Type: model.ActionDeny},
			},
			want: `add access-rule layer "Network" position.bottom "New rules" name "Policy with Destination-Negate Flag" source.1 "Source1" destination-negate true destination.1 "Destination1" destination.2 "Destination2" service.1 "Any" action "Drop"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertPolicy(tt.policy, ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertPolicy() error = %v,\nwantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConvertPolicy() = %q,\nwant %q", got, tt.want)
			}
		})
	}
}
