package main

import (
	"dtdao/greenlight/internal/data"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type TemplateData struct {
	Movie []*data.Movie
}

type MovieData struct {
	Movie *data.Movie
}

type movieCreateForm struct {
	Title   string       `form:"title"`
	Year    int32        `form:"year"`
	Runtime data.Runtime `form:"runtime"`
	Genres  string       `form:"genres"`
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
		"./internal/ui/html/partials/table.tmpl",
		"./internal/ui/html/partials/form.tmpl",
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

func (app *application) form(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "add" {
		http.NotFound(w, r)
		return
	}
	ts, err := template.ParseFiles("./internal/ui/html/partials/form.tmpl")

	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "form", nil)

	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func (app *application) create(w http.ResponseWriter, r *http.Request) {
	var form movieCreateForm

	err := r.ParseForm()

	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = app.formDecoder.Decode(&form, r.PostForm)

	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	movie := &data.Movie{
		Title:   form.Title,
		Year:    form.Year,
		Runtime: form.Runtime,
		Genres:  strings.Split(form.Genres, ","),
	}

	ts, err := template.ParseFiles("./internal/ui/html/partials/table.tmpl")

	app.Movies.Insert(movie)

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

	err = ts.ExecuteTemplate(w, "table", data)

	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func (app *application) edit(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		log.Print(err.Error())
		http.NotFound(w, r)
		return
	}

	movie, err := app.Movies.Get(id)

	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	data := &MovieData{
		Movie: movie,
	}

	ts, err := template.New("edit").Parse(`
	{{define "table-edit-item"}}
<tr hx-trigger='cancel' class='editing' hx-get="/contact/${contact.id}">
  <td>
    <button class="btn btn-danger" hx-get="/contact/${contact.id}">
      Cancel
    </button>
    <button class="btn btn-danger" hx-put="/contact/${contact.id}" hx-include="closest tr">
      Save
    </button>
  </td>
</tr>
{{end}}
	`)

	err = ts.ExecuteTemplate(w, "table-edit-item", data)

	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

}
