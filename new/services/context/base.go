package context

import "net/http"

func (ctx *Context) Redirect(location string) {
	http.Redirect(ctx.Resp, ctx.Req, location, http.StatusFound)
}
