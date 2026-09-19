package main

import (
	"net/http"

	"org.melements/skills-test/assets"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.FS(assets.EmbeddedFiles))
	mux.Handle("GET /static/", fileServer)

	mux.HandleFunc("GET /", app.home)

	return app.logRequest(app.recoverPanic(app.securityHeaders(mux)))
}
