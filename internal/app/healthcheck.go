package app

import (
	"net/http"

	"greenlight.dekutyavin.net/internal/data"
)

func (app *Application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := data.HealthCheckResponse{
		Status:      "available",
		Environment: app.Config.Env,
		Version:     app.Config.Version,
	}

	err := app.writeJSON(w, data, http.StatusOK)
	if err != nil {
		app.Logger.Error(err.Error())
		http.Error(w, "The server encountered an error and could not process your request", http.StatusInternalServerError)
		return
	}

}
