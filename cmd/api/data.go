package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Light2Dark/splitpay/internal/templates"
	"github.com/Light2Dark/splitpay/models"
)

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

	var serviceChargePercent = getServiceChargePercent(receipt.ServiceCharge, receipt.Subtotal)
	app.logger.Info("service charge percent", "val", serviceChargePercent)

	_, err = app.db.Exec(`
		UPDATE receipts
		SET items = ?,
			subtotal = ?,
			serviceCharge = ?,
			serviceChargePercent = ?,
			taxPercent = ?,
			taxAmount = ?,
			totalAmount = ?
		WHERE id = ?;
	`, string(itemsBytes), receipt.Subtotal, receipt.ServiceCharge, serviceChargePercent, receipt.TaxPercent, receipt.TaxAmount, receipt.TotalAmount, receipt.ID)

	if err != nil {
		app.logger.Error("update receipts table failed", "error", err)
		json.NewEncoder(w).Encode(&saveReceiptResponse{Message: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&saveReceiptResponse{Message: "OK"})
}

func (app application) viewReceiptHandler(w http.ResponseWriter, r *http.Request) {
	receiptLink := r.PathValue("receiptLink")

	var receipt models.Receipt
	row := app.db.QueryRow(`
		SELECT receipts.id, link, items, subtotal, serviceCharge, serviceChargePercent, taxPercent, taxAmount, totalAmount FROM splits
		INNER JOIN receipts
		ON splits.receipt_id = receipts.id
		WHERE splits.link = ?;
	`, receiptLink)

	var itemsStr string
	err := row.Scan(&receipt.ID, &receipt.Link, &itemsStr, &receipt.Subtotal, &receipt.ServiceCharge, &receipt.ServiceChargePercent, &receipt.TaxPercent, &receipt.TaxAmount, &receipt.TotalAmount)
	if err != nil {
		if err == sql.ErrNoRows {
			app.logger.Info("Receipt does not exist", "error", err, "link", receiptLink)
			templates.ErrorLayout("Sorry, this receipt does not exist").Render(r.Context(), w)
			return
		} else {
			app.logError(w, r, "Internal error", err)
			return
		}
	}

	err = json.Unmarshal([]byte(itemsStr), &receipt.Items)
	if err != nil {
		app.logError(w, r, "Internal error", err)
		return
	}

	var newReceipt models.ReceiptView
	var viewItems []models.ReceiptViewItem

	for _, item := range receipt.Items {
		var paidCount = item.PaidCount
		for range item.Quantity {
			var newItem models.ReceiptViewItem

			singleItemPrice := (item.Price / float64(item.Quantity))
			singleItemPrice = singleItemPrice + (singleItemPrice * float64(receipt.ServiceChargePercent) / 100)
			singleItemPrice = singleItemPrice + (singleItemPrice * float64(receipt.TaxPercent) / 100)
			newItem.Price = roundTo2DP(singleItemPrice)
			newItem.ID = item.ID
			newItem.Quantity = 1
			newItem.Name = item.Name

			if paidCount > 0 {
				newItem.Paid = true
				paidCount = paidCount - 1
			} else {
				newItem.Paid = false
			}
			viewItems = append(viewItems, newItem)
		}
	}
	newReceipt.Items = viewItems
	newReceipt.TotalAmount = receipt.TotalAmount
	newReceipt.ID = receipt.ID

	templates.ReceiptLayout(newReceipt).Render(r.Context(), w)
}

func (app application) payReceiptHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var totalAmount = r.FormValue("totalAmount")
	if totalAmount == "" {
		app.logger.Warn("Total Amount value not returned from frontend")
	}

	templates.PaymentView(totalAmount).Render(r.Context(), w)
}

func (app application) markPaidHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var receiptID string
	var totalAmount string
	var itemsIDToPaidCount = make(map[int]int)

	for key, values := range r.Form {
		if key == "receiptID" {
			receiptID = values[0]
			continue
		}

		if key == "totalAmount" {
			totalAmount = values[0]
			continue
		}

		itemID, err := strconv.Atoi(key)
		if err != nil {
			app.logError(w, r, "non-int values returned from form", err)
		}

		numItemsChecked := len(values)
		itemsIDToPaidCount[itemID] += numItemsChecked
	}
	app.logger.Info(fmt.Sprintf("Marking paid receiptID %s with items %v, totalAmount %s", receiptID, itemsIDToPaidCount, totalAmount))

	var itemsDB []models.ReceiptItem
	var itemsStr string

	row := app.db.QueryRow(`
		SELECT items FROM receipts
		WHERE id = ?;
	`, receiptID)

	err := row.Scan(&itemsStr)
	if err != nil {
		app.logError(w, r, "Error scanning DB", err)
		return
	}

	err = json.Unmarshal([]byte(itemsStr), &itemsDB)
	if err != nil {
		app.logError(w, r, "Error unmarshaling itemsStr", err)
		return
	}

	for id, paidCount := range itemsIDToPaidCount {
		for i := range itemsDB {
			if id == itemsDB[i].ID {
				itemsDB[i].PaidCount += paidCount
			}
		}
	}

	itemsBytes, err := json.Marshal(&itemsDB)
	if err != nil {
		app.logError(w, r, "Error marshaling itemsDB", err)
		return
	}

	_, err = app.db.Exec(`
		UPDATE receipts 
		SET items = ?
		WHERE id = ?;
	`, string(itemsBytes), receiptID)
	if err != nil {
		app.logError(w, r, "Error updating receipts table", err)
		return
	}

	templates.PaymentResult().Render(r.Context(), w)
}

func getServiceChargePercent(serviceCharge float64, subtotal float64) int {
	return  int(serviceCharge * 100 / subtotal)
}