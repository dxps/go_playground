package comps

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type Logo struct {
	app.Compo
}

func (l *Logo) Render() app.UI {
	return app.Div().Class("w-[26px] h-[26px] -mt-1").Body(
		app.Raw(LogoIcon),
	)
}
