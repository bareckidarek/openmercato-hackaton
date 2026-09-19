package main

import (
	"math/rand"
	"net/http"
	"time"

	"org.melements/skills-test/internal/response"
)

type kittyPicker struct {
	now  func() time.Time
	pick func(int) int
}

func newKittyPicker() kittyPicker {
	return kittyPicker{
		now:  time.Now,
		pick: rand.Intn,
	}
}

func (p kittyPicker) path() string {
	color := "black"
	if p.now().Minute()%2 != 0 {
		color = "red"
	}

	paths := kittyPaths[color]
	return paths[p.pick(len(paths))]
}

var kittyPaths = map[string][]string{
	"black": {
		"/static/img/kitty-black-1.svg",
		"/static/img/kitty-black-2.svg",
	},
	"red": {
		"/static/img/kitty-red-1.svg",
		"/static/img/kitty-red-2.svg",
	},
}

func (app *application) kittyPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data["KittyPath"] = app.kitty.path()

	if err := response.Page(w, http.StatusOK, data, "pages/kitty.tmpl"); err != nil {
		app.serverError(w, r, err)
	}
}
