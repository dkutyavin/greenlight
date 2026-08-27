package main

import (
	"encoding/json"
	"net/http"
)

type HealthCheckResponse struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := HealthCheckResponse{
		Status:      "available",
		Environment: app.config.env,
		Version:     version,
	}

	json, err := json.Marshal(data)
	// add an extra line for nice looking response via curl
	json = append(json, '\n')

	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered an error and could not process your request", http.StatusInternalServerError)
		return
	}

	w.Write(json)
}
