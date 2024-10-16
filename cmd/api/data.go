package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Light2Dark/splitpay/internal/templates"
	"github.com/Light2Dark/splitpay/models"
)

func (app application) dataHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := app.db.Query("SELECT * FROM splits")
	if err != nil {
		app.logger.Error("failed to execute query", "error", err)
		templates.Error(fmt.Sprintf("Error: %s", err)).Render(r.Context(), w)
		return
	}
	defer rows.Close()

	var splits []models.Splits
	for rows.Next() {
		var split models.Splits
		if err := rows.Scan(&split.ID, &split.Link); err != nil {
			app.logger.Error("failed to scan row", "error", err)
			templates.Error(fmt.Sprintf("Error: %s", err)).Render(r.Context(), w)
			return
		}
		splits = append(splits, split)
	}

	app.logger.Info("success query", "splits", splits)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&healthResponse{Status: "200 received bro"})
}