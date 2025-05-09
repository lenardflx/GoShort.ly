package routers

import (
	"goshortly/modules/web"
	web_routers "goshortly/routers/web"
)

func Routes() *web.Router {
	r := web.NewRouter()

	r.Use()
	r.Mount("/", web_routers.Routes())

	return r
}
