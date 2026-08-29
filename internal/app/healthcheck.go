package app

import (
	"net/http"
)

func (app *Application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := map[string]any{
		"status":      "available",
		"environment": app.Config.Env,
		"version":     app.Config.Version,
	}

	err := app.writeJSON(w, data, http.StatusOK)
	if err != nil {
		app.Logger.Error(err.Error())
		http.Error(w, "The server encountered an error and could not process your request", http.StatusInternalServerError)
		return
	}
}
