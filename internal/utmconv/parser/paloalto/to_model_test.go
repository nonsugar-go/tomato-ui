package paloalto

import (
	"reflect"
	"testing"

	"github.com/nonsugar-go/tomato-ui/internal/utmconv/model"
)

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
	in := []string{" service-http ", "any", "service-http", "application-default"}

	got := toSvcRefs(in)

	want := []model.ServiceRef{
		{Name: "service-http"},
		{Name: "application-default"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v want=%v", got, want)
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
