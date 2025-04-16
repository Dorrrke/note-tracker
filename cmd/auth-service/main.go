package main

import "github.com/Dorrrke/note-tracker/internal/auth-service/app"

func main() {
	app := app.NewApp()
	if err := app.StartApp(); err != nil {
		panic(err)
	}
}
