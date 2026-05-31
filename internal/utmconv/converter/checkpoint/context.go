package checkpoint

import "github.com/nonsugar-go/tomato-ui/internal/utmconv/model"

type Context struct {
	App     *model.App
	ZoneMap map[string]string
	AddrMap map[string]string
	SvcMap  map[string]string
}

func NewContext(app *model.App) *Context {
	ctx := Context{
		App:     app,
		ZoneMap: make(map[string]string),
		AddrMap: make(map[string]string),
		SvcMap:  make(map[string]string),
	}

	for _, v := range app.AppConfig.CheckPoint.Cli.ZoneReplacementMap.Value {
		ctx.ZoneMap[v.Before] = v.After
	}

	for _, v := range app.AppConfig.CheckPoint.Cli.AddressReplacementMap.Value {
		ctx.AddrMap[v.Before] = v.After
	}

	for _, v := range app.AppConfig.CheckPoint.Cli.PredefinedServices.Value {
		ctx.SvcMap[v] = v
	}

	for _, v := range app.AppConfig.CheckPoint.Cli.ServiceReplacementMap.Value {
		ctx.SvcMap[v.Before] = v.After
	}

	return &ctx
}
