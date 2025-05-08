package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"urlshort-backend/services/context"
)

type Router struct {
	chiRouter      *chi.Mux
	curGroupPrefix string
	curMiddlewares []any
}

func NewRouter() *Router {
	return &Router{
		chiRouter: chi.NewRouter(),
	}
}

// Global middleware
func (r *Router) Use(middlewares ...any) {
	for _, m := range middlewares {
		r.chiRouter.Use(toMiddleware(m))
	}
}

// Group routing
func (r *Router) Group(pattern string, fn func(), middlewares ...any) {
	prevPrefix := r.curGroupPrefix
	prevMiddlewares := r.curMiddlewares

	r.curGroupPrefix += pattern
	r.curMiddlewares = append(r.curMiddlewares, middlewares...)

	fn()

	r.curGroupPrefix = prevPrefix
	r.curMiddlewares = prevMiddlewares
}

// Shortcut methods
func (r *Router) Get(pattern string, handler ...any) {
	r.Method(http.MethodGet, pattern, handler...)
}
func (r *Router) Post(pattern string, handler ...any) {
	r.Method(http.MethodPost, pattern, handler)
}
func (r *Router) Put(pattern string, handler ...any) {
	r.Method(http.MethodPut, pattern, handler)
}
func (r *Router) Delete(pattern string, handler ...any) {
	r.Method(http.MethodDelete, pattern, handler)
}
func (r *Router) Patch(pattern string, handler ...any) {
	r.Method(http.MethodPatch, pattern, handler)
}

// Define a single method
func (r *Router) Method(method string, pattern string, args ...any) {
	mws, handler := wrapMiddlewareAndHandler(r.curMiddlewares, args)
	final := chainMiddleware(handler, mws)
	r.chiRouter.MethodFunc(method, r.fullPath(pattern), final)
}

// Define multiple methods
func (r *Router) Methods(methods string, pattern string, args ...any) {
	mws, handler := wrapMiddlewareAndHandler(r.curMiddlewares, args)
	final := chainMiddleware(handler, mws)

	for _, method := range strings.Split(methods, ",") {
		r.chiRouter.MethodFunc(strings.TrimSpace(method), r.fullPath(pattern), final)
	}
}

// Middleware + handler parsing
func wrapMiddlewareAndHandler(cur []any, args []any) ([]func(http.Handler) http.Handler, http.HandlerFunc) {
	if len(args) == 0 {
		panic("no handler provided")
	}

	var chain []func(http.Handler) http.Handler

	for _, m := range append(cur, args[:len(args)-1]...) {
		if m != nil {
			chain = append(chain, toMiddleware(m))
		}
	}

	last := args[len(args)-1]

	switch h := last.(type) {
	case http.HandlerFunc:
		return chain, h
	case http.Handler:
		return chain, h.ServeHTTP
	case func(*context.Context):
		wrapped := func(w http.ResponseWriter, r *http.Request) {
			ctx := context.Get(r)
			h(ctx)
		}
		return chain, wrapped
	default:
		panic(fmt.Sprintf("last argument must be http.HandlerFunc, http.Handler, or func(*context.Context), got %T", last))
	}

}

// Converts Gitea-style or standard middleware
func toMiddleware(m any) func(http.Handler) http.Handler {
	switch mw := m.(type) {
	case func(http.Handler) http.Handler:
		return mw

	case func(*context.Context):
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Wrap ResponseWriter only if it's not already wrapped
				resp := context.WrapResponseWriter(w)

				// Ensure ctx is available
				ctx := context.Get(r)
				if ctx == nil {
					panic("context not found: context.Middleware must be applied")
				}

				// Update ctx.Resp
				ctx.Resp = resp

				// Run the middleware
				mw(ctx)

				// Stop if written
				if ctx.Resp.Written() {
					return
				}

				next.ServeHTTP(ctx.Resp, r)
			})
		}

	default:
		panic("invalid middleware: must be func(http.Handler) or func(*context.Context)")
	}
}

// Chains middleware
func chainMiddleware(h http.HandlerFunc, mws []func(http.Handler) http.Handler) http.HandlerFunc {
	handler := http.Handler(h)
	for i := len(mws) - 1; i >= 0; i-- {
		handler = mws[i](handler)
	}
	return handler.ServeHTTP
}

// Path prefix utility
func (r *Router) fullPath(pattern string) string {
	full := r.curGroupPrefix + pattern
	if !strings.HasPrefix(full, "/") {
		full = "/" + full
	}
	if full == "/" {
		return full
	}
	return strings.TrimSuffix(full, "/")
}

// Extras
func (r *Router) Mount(path string, sub *Router) {
	sub.Use(r.curMiddlewares...)
	r.chiRouter.Mount(r.fullPath(path), sub.chiRouter)
}

func (r *Router) NotFound(h http.HandlerFunc) {
	r.chiRouter.NotFound(h)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, r2 *http.Request) {
	r.chiRouter.ServeHTTP(w, r2)
}
