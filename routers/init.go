package routers

import (
	"urlshort-backend/modules/web"
	web_routers "urlshort-backend/routers/web"
)

func Routes() *web.Router {
	r := web.NewRouter()

	r.Use()
	r.Mount("/", web_routers.Routes())

	return r
}
