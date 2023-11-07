package main

import (
	"dtdao/greenlight/internal/data"
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type TemplateData struct {
	Movie []*data.Movie
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}

	input.Title = ""
	input.Genres = []string{}
	input.Page = 1
	input.PageSize = 20
	input.Sort = "id"

	input.Filters.SortSafeList = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	movies, _, err := app.Movies.GetAll(input.Title, input.Genres, input.Filters)

	data := &TemplateData{
		Movie: movies,
	}

	files := []string{
		"./internal/ui/html/base.tmpl",
		"./internal/ui/html/partials/nav.tmpl",
		"./internal/ui/html/pages/home.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.ServeFiles("/static/*filepath", http.Dir("./internal/ui/static/"))
	return router
}
