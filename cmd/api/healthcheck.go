package main

import (
	"fmt"
	"net/http"
)

func (app *application) healtcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environemtn %s\n", app.config.env)
	fmt.Fprintf(w, "version %s\n", version)
}
