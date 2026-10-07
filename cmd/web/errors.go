package main

import "net/http"

func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		url    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "url", url)
}

// func (app *application) errResponse(w http.ResponseWriter, r *http.Request, status int, message any) {

// }
