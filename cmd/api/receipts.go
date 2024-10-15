package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/Light2Dark/splitpay/internal/templates"
	"github.com/Light2Dark/splitpay/models"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

func (app application) scanReceiptHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("receipt")
	if err != nil {
		app.logError(w, r, "Error when scanning file", err)
		return
	}
	defer file.Close()

	base64Img, err := imageFileToBase64(file)
	if err != nil {
		app.logError(w, r, "Error when converting image to base64", err)
		return
	}

	var receiptOpenAI models.ReceiptOpenAI

	schema, err := jsonschema.GenerateSchemaForType(receiptOpenAI)
	if err != nil {
		app.logError(w, r, "GenerateSchemaForType error", err)
		return
	}

	resp, err := app.openai.CreateChatCompletion(
		r.Context(),
		openai.ChatCompletionRequest{
			Model:     openai.GPT4o,
			MaxTokens: 1000,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleUser,
					MultiContent: []openai.ChatMessagePart{
						{
							Type: openai.ChatMessagePartTypeText,
							Text: "This is a receipt. I would like you to return a list of items that were ordered, alongside their price and quantity. The quantity might be at the start like '1' or in the middle. There may be times when there are additions to the item, like HOT / + milk underneath Coffee. These rows will not have a corresponding charge. Ignore them. Lastly, return the service charge amount and tax (amount and percent) that was charged. Check the image carefully and spend some time to extract accurate information. Do not change the words of what is displayed as it may be in different languages.",
						},
						{
							Type: openai.ChatMessagePartTypeImageURL,
							ImageURL: &openai.ChatMessageImageURL{
								URL:    base64Img,
								Detail: openai.ImageURLDetailLow,
							},
						},
					}},
			},
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
				JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
					Name:   "receipt",
					Schema: schema,
					Strict: true,
				},
			},
		},
	)

	if err != nil {
		app.logError(w, r, "Error calling openai", err)
		return
	}

	err = schema.Unmarshal(resp.Choices[0].Message.Content, &receiptOpenAI)
	if err != nil {
		app.logError(w, r, "Unmarschal schema error", err)
		return
	}

	// convert OpenAIReceipt to actual receipt
	var receipt models.Receipt
	receipt.ServiceCharge = receiptOpenAI.ServiceCharge
	receipt.TotalAmount = receiptOpenAI.TotalAmount
	receipt.TaxAmount = receiptOpenAI.TaxAmount
	receipt.TaxPercent = receiptOpenAI.TaxPercent

	var itemCount int = 1
	for _, item := range receiptOpenAI.Items {
		var itemAI = struct {
			ID       int
			Name     string
			Quantity int
			Price    float64
		}{
			ID:       itemCount,
			Name:     item.Name,
			Quantity: item.Quantity,
			Price:    item.Price,
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

	if subtotal != receipt.Subtotal {
		app.logger.Info(fmt.Sprintf("incorrect subtotal by openai, openai: %f, expected: %f", receipt.Subtotal, subtotal))
		receipt.Subtotal = roundTo2DP(subtotal)
	}

	var taxAmount = roundTo2DP(subtotal * 0.06)
	if taxAmount != receipt.TaxAmount {
		app.logger.Info(fmt.Sprintf("incorrect tax amount by openai, openai: %f, expected: %f", receipt.TaxAmount, taxAmount))
		receipt.TaxAmount = roundTo2DP(taxAmount)
	}

	var totalAmountExpected = roundTo2DP(subtotal + taxAmount + receipt.ServiceCharge)
	if totalAmountExpected != receipt.TotalAmount {
		app.logger.Info(fmt.Sprintf("incorrect total amount by openai, openai: %f, expected: %f", receipt.TotalAmount, totalAmountExpected))
	}

	app.receipt = receipt
	templates.ReceiptTable(receipt).Render(r.Context(), w)
}

func (app application) deleteItemHandler(w http.ResponseWriter, r *http.Request) {
	itemNumStr := r.PathValue("itemNum")
	itemNum, err := strconv.Atoi(itemNumStr)
	if err != nil {
		app.logger.Error("non-integer passed to delete item handler", "itemNum", itemNumStr)
		return
	}

	for ind, item := range app.receipt.Items {
		if item.ID == itemNum {
			app.receipt.Items = append(app.receipt.Items[:ind], app.receipt.Items[ind+1:]...)
			app.receipt.Subtotal = app.receipt.Subtotal - (item.Price * float64(item.Quantity))
			app.receipt.TaxAmount = app.receipt.Subtotal * float64(app.receipt.TaxPercent)
			app.receipt.TotalAmount = app.receipt.Subtotal + app.receipt.TaxAmount + app.receipt.ServiceCharge
			break
		}
	}

	app.logger.Info("deleted", "receipt", app.receipt)
}

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func imageFileToBase64(file multipart.File) (string, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return "", err
	}
	imgBytes := buf.Bytes()
	var base64Encoding string

	// Determine the content type of the image file
	mimeType := http.DetectContentType(imgBytes)

	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}
	base64Encoding += toBase64(imgBytes)

	return base64Encoding, nil
}

// TODO: Change to truncate(?)
func roundTo2DP(val float64) float64 {
	return math.Round(val*100) / 100
}
