package auth

import (
	"goshortly/models"
	"goshortly/modules/web"
	"goshortly/services/context"
	"goshortly/services/forms"
	"log"
	"net/http"
)

const (
	tplSignIn = "auth/signin"
)

func SignIn(ctx *context.Context) {
	if ctx.IsSignedIn {
		ctx.Redirect("/admin")
		return
	}

	ctx.Data["Title"] = "Sign In"
	ctx.HTML(http.StatusOK, tplSignIn)
}

func SignInPost(ctx *context.Context) {
	log.Printf("SignInPost: %v", ctx.Req.Form)
	form := web.GetForm[forms.SignInForm](ctx.Req)
	log.Printf("SignInPost: %v", form)

	if form == nil {
		ctx.Data["ErrorMsg"] = "Invalid submission."
		ctx.HTML(http.StatusBadRequest, tplSignIn)
		return
	}

	ctx.Data["Title"] = "Sign In"
	ctx.Data["Form"] = form

	if form.UserName != "admin" || form.Password != "password" {
		ctx.Data["ErrorMsg"] = "Invalid username or password."
		ctx.HTML(http.StatusUnauthorized, tplSignIn)
		return
	}

	// Dummy signed-in user
	ctx.Doer = &models.User{
		Username: form.UserName,
		IsAdmin:  true,
	}
	ctx.IsSignedIn = true

	ctx.Redirect("/admin")
}
