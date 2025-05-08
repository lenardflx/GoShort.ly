package web

import (
	"net/http"
	"urlshort-backend/modules/templates"
	"urlshort-backend/modules/web"
	"urlshort-backend/routers/common"
	"urlshort-backend/routers/web/admin"
	"urlshort-backend/services/context"
)

// Routes returns all web routes
func Routes() *web.Router {
	webRoutes := web.NewRouter()

	_ = templates.HTMLRenderer()

	var mid []any

	mid = append(mid, context.Contexter())

	webRoutes.Use(mid...)

	webRoutes.Group("", func() { registerWebRoutes(webRoutes) })

	return webRoutes
}

func registerWebRoutes(r *web.Router) {
	adminReq := VerifyAuthWithOptions(&common.VerifyOptions{
		SignInRequired: true,
		AdminRequired:  true,
	})

	r.Group("/admin", func() {
		r.Get("", admin.Dashboard)
	}, adminReq)

	r.Get("/test", admin.Test)

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		ctx := context.Get(req)
		ctx.NotFound()
	})

}

func VerifyAuthWithOptions(opt *common.VerifyOptions) func(ctx *context.Context) {
	return func(ctx *context.Context) {
		if opt.SignInRequired && !ctx.IsSignedIn {
			http.Redirect(ctx.Resp, ctx.Req, "/login", http.StatusFound)
			return
		}

		if opt.AdminRequired {
			if !ctx.Doer.IsAdmin {
				http.Error(ctx.Resp, "Forbidden", http.StatusForbidden)
				return
			}
			ctx.Data["PageAdmin"] = true
		}
	}
}
