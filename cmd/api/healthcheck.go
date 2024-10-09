package main

import (
	"encoding/json"
	"net/http"
)

type healthResponse struct {
	Status string `json:"status"`
}

func (app application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Healthcheck endpoint hit", "request", r)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&healthResponse{Status: "OK"})
}
