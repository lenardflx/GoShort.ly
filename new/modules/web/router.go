package web

import (
	"encoding/json"
	"fmt"
	"goshortly/services/context"
	"log"
	"net/http"
	"reflect"
	"strings"

	stdctx "context"
	"github.com/go-chi/chi/v5"
)

type contextKey string

const formKey contextKey = "__form"

// Bind parses form or JSON into a struct and stores it in the request context
func Bind[T any]() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			form := new(T)

			contentType := r.Header.Get("Content-Type")
			method := r.Method

			if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
				switch {
				case strings.Contains(contentType, "form-urlencoded"):
					if err := r.ParseForm(); err != nil {
						http.Error(w, "Invalid form data", http.StatusBadRequest)
						return
					}
					bindForm(r.Form, form)
				case strings.Contains(contentType, "multipart/form-data"):
					if err := r.ParseMultipartForm(32 << 20); err != nil {
						http.Error(w, "Invalid multipart form data", http.StatusBadRequest)
						return
					}
					bindForm(r.MultipartForm.Value, form)
				case strings.Contains(contentType, "json"):
					if err := json.NewDecoder(r.Body).Decode(form); err != nil {
						http.Error(w, "Invalid JSON", http.StatusBadRequest)
						return
					}
				default:
					log.Printf("[Bind] Unsupported content-type: %q", contentType)
				}
			}

			ctx := stdctx.WithValue(r.Context(), formKey, form)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetForm extracts the parsed form from request context.
func GetForm[T any](r *http.Request) *T {
	val := r.Context().Value(formKey)
	if typed, ok := val.(*T); ok {
		return typed
	}
	return nil
}

func bindForm(values map[string][]string, dst any) {
	v := reflect.ValueOf(dst).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		key := field.Tag.Get("form")
		if key == "" {
			key = field.Name
		}
		vals, ok := values[key]
		if !ok || len(vals) == 0 {
			continue
		}
		val := vals[0]
		fv := v.Field(i)
		if !fv.CanSet() {
			continue
		}
		switch fv.Kind() {
		case reflect.String:
			fv.SetString(val)
		case reflect.Bool:
			fv.SetBool(val == "true" || val == "on" || val == "1")
		default:
			panic("unhandled default case")
		}
	}
}

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

func (r *Router) Use(middlewares ...any) {
	for _, m := range middlewares {
		if m != nil {
			r.chiRouter.Use(toMiddleware(m))
		}
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
	r.Methods(http.MethodGet, pattern, handler...)
}
func (r *Router) Post(pattern string, handler ...any) {
	r.Methods(http.MethodPost, pattern, handler...)
}
func (r *Router) Put(pattern string, handler ...any) {
	r.Methods(http.MethodPut, pattern, handler)
}
func (r *Router) Delete(pattern string, handler ...any) {
	r.Methods(http.MethodPost, pattern, handler...)
}
func (r *Router) Patch(pattern string, handler ...any) {
	r.Methods(http.MethodPost, pattern, handler...)
}

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
		fmt.Printf("DEBUG args: %#v\n", args)
		panic(fmt.Sprintf("last argument must be http.HandlerFunc, http.Handler, or func(*context.Context), got %T", last))
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
