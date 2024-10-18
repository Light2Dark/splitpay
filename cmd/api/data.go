package main

import (
	"encoding/json"
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

func (app application) viewReceiptHandler(w http.ResponseWriter, r *http.Request) {
	receiptLink := r.PathValue("receiptLink")
	app.logger.Info("receipt link", "receiptLink", receiptLink)

	// templates.PaymentView("1", 0, map[int]int{1: 2, 2: 3}).Render(r.Context(), w)
	// return

	// TODO: get receipt from DB
	var receipt = models.MockReceipt
	var newReceipt models.ReceiptView

	var serviceChargePercent = int(receipt.ServiceCharge * 100 / receipt.Subtotal)
	app.logger.Info("service charge percent", "val", serviceChargePercent)

	var viewItems []models.ReceiptViewItem

	for _, item := range receipt.Items {
		var paidCount = item.PaidCount
		for range item.Quantity {
			var newItem models.ReceiptViewItem

			singleItemPrice := (item.Price / float64(item.Quantity))
			singleItemPrice = singleItemPrice + (singleItemPrice * float64(serviceChargePercent) / 100)
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

// todo: they will pass back item ID of the thing they clicked?
// we will save that in DB that the item is paid?
// paid counter

func (app application) payReceiptHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// TODO: get paidCount from DB

	var receiptID string
	var itemsIDToPaidCount = make(map[int]int)
	for key, values := range r.Form {
		if key == "receiptID" {
			receiptID = values[0]
			continue
		}

		itemID, err := strconv.Atoi(key)
		if err != nil {
			app.logError(w, r, "non-int values returned from form", err)
		}

		numItemsChecked := len(values)
		itemsIDToPaidCount[itemID] += numItemsChecked
	}

	templates.PaymentView(receiptID, 0, itemsIDToPaidCount).Render(r.Context(), w)
}
