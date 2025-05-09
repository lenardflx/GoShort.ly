package web

import (
	"fmt"
	"goshortly/services/context"
	"net/http"
)

func toMiddleware(m any) func(next http.Handler) http.Handler {
	switch mw := m.(type) {

	case func(http.Handler) http.Handler:
		return mw

	case func(*context.Context):
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				resp := context.WrapResponseWriter(w)
				ctx := context.Get(r)

				if ctx == nil {
					panic("context.Context is missing: make sure context.Contexter() middleware is used")
				}

				ctx.Resp = resp

				mw(ctx)

				if ctx.Resp.Written() {
					return
				}

				next.ServeHTTP(resp, r)
			})
		}

	default:
		panic(fmt.Sprintf("invalid middleware type: got %T, want func(http.Handler) or func(*context.Context)", m))
	}
}
