package context

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
	"urlshort-backend/modules/setting"
)

const tplStatus500 string = "status/500"

func (ctx *Context) ServerError(logMsg string, logErr error) {
	if logErr != nil {
		var opError *net.OpError
		if errors.As(logErr, &opError) || errors.Is(logErr, &net.OpError{}) {
			return
		}

		if !setting.IsProd || (ctx.Doer != nil && ctx.Doer.IsAdmin) {
			ctx.Data["ErrorMsg"] = fmt.Sprintf("%s, %s", logMsg, logErr)
		}
	}

	ctx.Data["Title"] = "Internal Server Error"
	ctx.HTML(http.StatusInternalServerError, tplStatus500)
}

func (ctx *Context) HTML(status int, name string) {
	tmplStart := time.Now()

	if !setting.IsProd {
		ctx.Data["TemplateName"] = name
	}

	ctx.Data["TemplateLoadTime"] = func() string {
		return strconv.FormatInt(time.Since(tmplStart).Milliseconds(), 10) + "ms"
	}

	ctx.Resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.Resp.WriteHeader(status)

	err := ctx.Render.ExecuteTemplate(ctx.Resp, name, ctx.Data)

	if err == nil {
		return
	}

	log.Printf("template render failed: %v", err)
	ctx.ServerError("Render failed", err)
}

func (ctx *Context) NotFound() {
	ctx.Data["Title"] = "Page Not Found"
	ctx.HTML(http.StatusNotFound, "status/404")
}
