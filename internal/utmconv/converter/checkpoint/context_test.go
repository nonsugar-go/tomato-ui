package checkpoint

import (
	"testing"

	"github.com/nonsugar-go/tomato-ui/internal/utmconv/model"
)

func TestNewContext_WithInitializedApp(t *testing.T) {
	var app model.App
	app.AppConfig = *model.NewDefaultAppConfig()
	ctx := NewContext(&app)

	expectedSvc := map[string]string{
		"service-http":  "http",
		"service-https": "https",
	}

	for _, v := range app.AppConfig.CheckPoint.Cli.PredefinedServices.Value {
		expectedSvc[v] = v
	}

	if len(ctx.SvcMap) != len(expectedSvc) {
		t.Fatalf("SvcMap length mismatch: got %d, want %d", len(ctx.SvcMap), len(expectedSvc))
	}

	for k, wantV := range expectedSvc {
		if gotV, ok := ctx.SvcMap[k]; !ok || gotV != wantV {
			t.Errorf("SvcMap[%s]: got %s, want %s", k, gotV, wantV)
		}
	}
}
