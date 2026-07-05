package fortigate

import (
	"log/slog"

	"github.com/nonsugar-go/tomato-ui/internal/utmconv/model"
)

func ParseFortiGate(app *model.App) {
	slog.Error("unsupported vendor", "vendor", app.To)
}
