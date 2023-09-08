package main

import (
	"net/http"
)

func (app *application) healtcheckHandler(w http.ResponseWriter, r *http.Request) {

	systemInfo := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	data := envelope{
		"status":      "available",
		"system_info": systemInfo,
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
