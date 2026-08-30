package app

import (
	"net/http"
)

func (app *Application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := envelope{
		"status": "available",
		"system_info": map[string]any{
			"environment": app.Config.Env,
			"version":     app.Config.Version,
		},
	}

	err := app.writeJSON(w, data, http.StatusOK)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
