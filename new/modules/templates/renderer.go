package templates

import (
	"html/template"
	"log"
	"path/filepath"
	"sync"
)

var (
	tmpl     *template.Template
	initOnce sync.Once
)

// HTMLRenderer returns a lazily loaded *template.Template instance.
func HTMLRenderer() *template.Template {
	initOnce.Do(func() {
		funcs := template.FuncMap{
			// `ctx` will be called in templates as `{{ctx}}`
			"ctx": func() any {
				// fallback: return nil or dummy for now
				return nil
			},
			// add your other global helpers here
			"AppName": func() string { return "GoShort.ly" },
			// "AssetUrlPrefix", "AppSubUrl", etc. can also go here
		}

		var err error
		tmpl, err = template.New("").Funcs(funcs).ParseGlob(filepath.Join("templates", "**", "*.gohtml"))
		if err != nil {
			log.Fatalf("failed to parse templates: %v", err)
		}
	})
	return tmpl
}
