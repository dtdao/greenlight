package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.Print(err)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {

	errMessage := envelope{"error": message}

	err := app.writeJSON(w, status, errMessage, nil)

	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {

	app.logError(r, err)
	message := "the server encountered a problem and could not process your request"

	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {

	message := "the requested resource could not be found"

	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)

	app.errorResponse(w, r, http.StatusInternalServerError, message)
}