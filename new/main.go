package main

import (
	"goshortly/cmd"
	"goshortly/modules/setting"
	"log"
)

var (
	Version = "v0.0.1dev"
	Author  = "Lenard Felix"
)

func init() {
	setting.AppVersion = Version
	setting.AppAuthor = Author
}

func main() {
	app := cmd.NewApp()
	if err := cmd.RunApp(app); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
