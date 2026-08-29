package main

import (
	"net/http"

	"greenlight.dekutyavin.net/internal/data"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := data.HealthCheckResponse{
		Status:      "available",
		Environment: app.config.Env,
		Version:     version,
	}

	err := app.writeJSON(w, data, http.StatusOK)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered an error and could not process your request", http.StatusInternalServerError)
		return
	}

}
