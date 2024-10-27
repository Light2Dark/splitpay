package receipts

import (
	"fmt"
	"math"

	"github.com/Light2Dark/splitpay/models"
)

func ConvertOpenAIReceiptToReceipt(AIReceipt models.ReceiptOpenAI) (models.Receipt, []string) {
	var dataMessages []string

	// convert OpenAIReceipt to actual receipt
	var receipt models.Receipt
	receipt.ServiceCharge = AIReceipt.ServiceCharge
	receipt.TotalAmount = AIReceipt.TotalAmount
	receipt.TaxAmount = AIReceipt.TaxAmount
	receipt.TaxPercent = AIReceipt.TaxPercent
	receipt.Discount = AIReceipt.Discount

	var itemCount int = 1
	for _, item := range AIReceipt.Items {
		var itemAI = models.ReceiptItem{
			ReceiptItemBase: models.ReceiptItemBase{
				ID:       itemCount,
				Name:     item.Name,
				Quantity: item.Quantity,
				Price:    item.Price,
			},
			PaidCount: 0,
		}
		receipt.Items = append(receipt.Items, itemAI)
		itemCount = itemCount + 1
	}

	// totalAmount is almost always true, but the indiv items will be not true, so total will be false.
	// we need to get a single item price and store that in receipt, to improve accuracy
	var subtotal float64
	for _, item := range receipt.Items {
		totalPrice := item.Price
		subtotal += totalPrice
		qty := item.Quantity

		singleItemPrice := totalPrice / float64(qty)
		singleItemPrice = roundTo2DP(singleItemPrice)
		item.Price = singleItemPrice
	}

	// subtotal = total of items.
	// - overall discount = amount to be service charged
	// + service charge = pre tax total.
	// + tax = total

	if subtotal != receipt.Subtotal {
		dataMessages = append(dataMessages, fmt.Sprintf("incorrect subtotal by openai, openai: %f, expected: %f", receipt.Subtotal, subtotal))
		receipt.Subtotal = roundTo2DP(subtotal)
	}

	// change discount to negative
	if receipt.Discount > 0 {
		receipt.Discount = roundTo2DP(-1 * receipt.Discount)
	}
	receipt.DiscountPercent = int(receipt.Discount * 100 / subtotal)

	var toBeServiceCharged = subtotal + receipt.Discount
	receipt.ServiceChargePercent = GetServiceChargePercent(receipt.ServiceCharge, toBeServiceCharged)
	dataMessages = append(dataMessages, fmt.Sprintf("service charge percent, val: %d", receipt.ServiceChargePercent))

	var preTax = toBeServiceCharged + receipt.ServiceCharge

	var taxAmount = roundTo2DP(preTax * 0.06)
	receipt.TaxPercent = 6 // hardcode for now
	if taxAmount != receipt.TaxAmount {
		dataMessages = append(dataMessages, fmt.Sprintf("incorrect tax amount by openai, openai: %f, expected: %f", receipt.TaxAmount, taxAmount))
	}

	var totalAmountExpected = roundTo2DP(preTax + taxAmount)
	if totalAmountExpected != receipt.TotalAmount {
		receipt.TotalAmount = totalAmountExpected
		dataMessages = append(dataMessages, fmt.Sprintf("incorrect total amount by openai, openai: %f, expected: %f", receipt.TotalAmount, totalAmountExpected))
	}

	return receipt, dataMessages
}

// TODO: Change to truncate(?)
// This func is also used in data.go
func roundTo2DP(val float64) float64 {
	return math.Round(val*100) / 100
}

func GetServiceChargePercent(serviceCharge float64, total float64) int {
	return int(serviceCharge * 100 / total)
}
