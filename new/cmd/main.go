package cmd

import (
	"errors"
	"fmt"
	"goshortly/routers"
	"log"
	"net/http"
	"time"
)

type App struct {
	Name      string
	Usage     string
	UsageText string
	Compiled  time.Time
}

func NewApp() *App {
	app := &App{
		Name:      "GoShort.ly",
		Usage:     "A URL shortener written in Go",
		UsageText: "GoShort.ly is a self-hosted URL shortener that is easy to set up and use.",
		Compiled:  time.Now(),
	}

	return app
}

func RunApp(app *App) error {
	fmt.Printf("%s — %s\nCompiled: %s\n\n", app.Name, app.Usage, app.Compiled.Format(time.RFC1123))

	routes := routers.Routes()

	server := &http.Server{
		Addr:    ":8080",
		Handler: routes,
	}

	log.Println("🚀 Server running at http://localhost:8080")
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
