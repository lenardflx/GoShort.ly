package context

import (
	"context"
	"goshortly/models"
	"goshortly/modules/templates"
	"html/template"
	"net/http"
)

type contextKeyType string

const contextKey = contextKeyType("request-context")

type Context struct {
	Resp *Response
	Req  *http.Request

	Render *template.Template

	Doer       *models.User
	IsSignedIn bool
	IsAdmin    bool

	Data map[string]any
}

func Get(r *http.Request) *Context {
	val := r.Context().Value(contextKey)
	if ctx, ok := val.(*Context); ok {
		return ctx
	}
	return nil
}

func Contexter() func(http.Handler) http.Handler {
	tmpl := templates.HTMLRenderer()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := WrapResponseWriter(w)
			ctx := &Context{
				Resp:   resp,
				Req:    r,
				Render: tmpl,
				Data:   map[string]any{},
			}
			ctxReq := r.WithContext(context.WithValue(r.Context(), contextKey, ctx))
			ctx.Req = ctxReq

			next.ServeHTTP(resp, ctxReq)
		})
	}
}
