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

	var newReceipt = CalculateReceipt(receipt)
	if newReceipt.Subtotal != receipt.Subtotal {
		dataMessages = append(dataMessages, fmt.Sprintf("incorrect subtotal by openai, openai: %f, expected: %f", receipt.Subtotal, newReceipt.Subtotal))
	}
	dataMessages = append(dataMessages, fmt.Sprintf("service charge percent, val: %d", receipt.ServiceChargePercent))

	if newReceipt.TaxAmount != receipt.TaxAmount {
		dataMessages = append(dataMessages, fmt.Sprintf("incorrect tax amount by openai, openai: %f, expected: %f", receipt.TaxAmount, newReceipt.TaxAmount))
	}

	if newReceipt.TotalAmount != receipt.TotalAmount {
		dataMessages = append(dataMessages, fmt.Sprintf("incorrect total amount by openai, openai: %f, expected: %f", receipt.TotalAmount, newReceipt.TotalAmount))
	}

	return receipt, dataMessages
}

// function will correctify the receipt by calculating tax, subtotal, total, service charge
func CalculateReceipt(receipt models.Receipt) models.Receipt {
	var subtotal float64
	for _, item := range receipt.Items {
		totalPrice := item.Price
		subtotal += totalPrice
	}
	receipt.Subtotal = RoundTo2DP(subtotal)

	// subtotal = total of items.
	// - overall discount = amount to be service charged
	// + service charge = pre tax total.
	// + tax = total

	// change discount to negative
	if receipt.Discount > 0 {
		receipt.Discount = RoundTo2DP(-1 * receipt.Discount)
	}
	receipt.DiscountPercent = int(receipt.Discount * 100 / subtotal)

	var toBeServiceCharged = subtotal + receipt.Discount
	receipt.ServiceChargePercent = GetServiceChargePercent(receipt.ServiceCharge, toBeServiceCharged)

	var preTax = toBeServiceCharged + receipt.ServiceCharge

	var taxAmount = RoundTo2DP(preTax * 0.06)
	receipt.TaxPercent = 6 // hardcode for now
	receipt.TaxAmount = taxAmount
	receipt.TotalAmount = RoundTo2DP(preTax + taxAmount)

	return receipt
}

// TODO: Change to truncate(?)
// This func is also used in data.go
func RoundTo2DP(val float64) float64 {
	return math.Round(val*100) / 100
}

func GetServiceChargePercent(serviceCharge float64, total float64) int {
	return int(serviceCharge * 100 / total)
}

func CalculateSingleItemPrice(quantity int, price float64) float64 {
	singleItemPrice := price / float64(quantity)
	return RoundTo2DP(singleItemPrice)
}
