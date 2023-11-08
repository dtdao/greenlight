package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/add", app.form)
	router.HandlerFunc(http.MethodPost, "/create", app.create)

	router.HandlerFunc(http.MethodGet, "/edit/:id", app.edit)

	router.ServeFiles("/static/*filepath", http.Dir("./internal/ui/static/"))
	return router
}
