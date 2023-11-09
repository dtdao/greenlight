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
	router.HandlerFunc(http.MethodGet, "/movie/:id", app.getItem)
	router.HandlerFunc(http.MethodPut, "/movie/:id", app.update)

	router.ServeFiles("/static/*filepath", http.Dir("./internal/ui/static/"))
	return router
}
