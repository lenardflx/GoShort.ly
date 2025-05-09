package web

import (
	"goshortly/modules/templates"
	"goshortly/modules/web"
	"goshortly/routers/common"
	"goshortly/routers/web/admin"
	"goshortly/routers/web/auth"
	"goshortly/services/context"
	"goshortly/services/forms"
	"net/http"
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
	reqSignIn := verifyAuthWithOptions(&common.VerifyOptions{SignInRequired: true})
	reqSignOut := verifyAuthWithOptions(&common.VerifyOptions{SignOutRequired: true})
	reqAdmin := verifyAuthWithOptions(&common.VerifyOptions{SignInRequired: true, AdminRequired: true})

	// User
	r.Get("/user/login", auth.SignIn)
	r.Group("/user", func() {
		r.Post("/login", web.Bind[forms.SignInForm](), auth.SignInPost)
	}, reqSignOut)

	r.Group("/user", func() {

	}, reqSignIn)

	// Admin
	r.Group("/admin", func() {
		r.Get("", admin.Dashboard)
	}, reqAdmin)

	r.Get("/test", admin.Test)

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		ctx := context.Get(req)
		ctx.NotFound()
	})

}

func verifyAuthWithOptions(opt *common.VerifyOptions) func(ctx *context.Context) {
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
