package admin

import (
	"net/http"
	"urlshort-backend/services/context"
)

const (
	tplAdminDashboard = "admin"
	tplTest           = "test"
)

func Dashboard(ctx *context.Context) {
	ctx.Data["Title"] = "Admin Dashboard"
	ctx.Data["Username"] = ctx.Doer.Username

	ctx.HTML(http.StatusOK, tplAdminDashboard)
}

func Test(ctx *context.Context) {
	ctx.Data["Title"] = "Debug Info"
	ctx.Data["Method"] = ctx.Req.Method
	ctx.Data["URL"] = ctx.Req.URL.String()
	ctx.Data["IsSignedIn"] = ctx.IsSignedIn
	ctx.Data["IsAdmin"] = ctx.IsAdmin

	if ctx.Doer != nil {
		ctx.Data["User"] = map[string]any{
			"ID":       ctx.Doer.ID,
			"Username": ctx.Doer.Username,
			"Email":    ctx.Doer.Email,
		}
	} else {
		ctx.Data["User"] = "nil"
	}

	ctx.HTML(http.StatusOK, "debug/info")
}
