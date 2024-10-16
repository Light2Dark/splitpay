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

type saveReceiptResponse struct {
	Message string `json:"message"`
}

func (app application) saveReceiptHandler(w http.ResponseWriter, r *http.Request) {
	var receipt models.Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		app.logger.Error("unable to save receipt", "error", err)
		json.NewEncoder(w).Encode(&saveReceiptResponse{Message: err.Error()})
		return
	}

	itemsBytes, err := json.Marshal(&receipt.Items)
	if err != nil {
		app.logger.Error("unable to marshal receipt items", "error", err)
		json.NewEncoder(w).Encode(&saveReceiptResponse{Message: err.Error()})
		return
	}

	_, err = app.db.Exec(`
		UPDATE receipts
		SET items = ?,
			subtotal = ?,
			serviceCharge = ?,
			taxPercent = ?,
			taxAmount = ?,
			totalAmount = ?
		WHERE id = ?;
	`, string(itemsBytes), receipt.Subtotal, receipt.ServiceCharge, receipt.TaxPercent, receipt.TaxAmount, receipt.TotalAmount, receipt.ID)

	if err != nil {
		app.logger.Error("update receipts table failed", "error", err)
		json.NewEncoder(w).Encode(&saveReceiptResponse{Message: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&saveReceiptResponse{Message: "OK"})
}
