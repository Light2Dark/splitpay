package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"math"
	"math/rand"
	"mime/multipart"
	"net/http"

	"github.com/Light2Dark/splitpay/internal/receipts"
	"github.com/Light2Dark/splitpay/internal/templates"
	"github.com/Light2Dark/splitpay/models"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

func (app *application) scanReceiptHandler(w http.ResponseWriter, r *http.Request) {
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

	var model = openai.GPT4o
	if app.env == "DEV" {
		model = openai.GPT4oMini
	}

	resp, err := app.openai.CreateChatCompletion(
		r.Context(),
		openai.ChatCompletionRequest{
			Model:     model,
			MaxTokens: 1000,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleUser,
					MultiContent: []openai.ChatMessagePart{
						{
							Type: openai.ChatMessagePartTypeText,
							Text: "This is a receipt. I would like you to return a list of items that were ordered, alongside their price and quantity. The quantity might be at the start like '1' or in the middle. There may be times when there are additions to the item, like HOT / + milk underneath Coffee. These rows will not have a corresponding charge. Ignore them. Lastly, return the service charge amount and tax (amount and percent) that was charged. If there is a discount, identify the amount too. Check the image carefully and spend some time to extract accurate information. Do not change the words of what is displayed as it may be in different languages.",
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

	receipt, messages := receipts.ConvertOpenAIReceiptToReceipt(receiptOpenAI)
	for _, msg := range messages {
		app.logger.Info(msg)
	}

	itemsJson, err := json.Marshal(receipt.Items)
	if err != nil {
		app.logger.Error("unable to marshal json", "error", err)
	}
	row := app.db.QueryRow(`
		INSERT INTO receipts (items, subtotal, serviceCharge, serviceChargePercent, taxPercent, taxAmount, discount, discountPercent, totalAmount)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id;
	`, string(itemsJson), receipt.Subtotal, receipt.ServiceCharge, receipt.ServiceChargePercent, receipt.TaxPercent, receipt.TaxAmount, receipt.Discount, receipt.DiscountPercent, receipt.TotalAmount)

	var id int
	err = row.Scan(&id)

	if err != nil {
		app.logger.Error("error executing query", "err", err)
	}
	receipt.ID = id

	key := generateShortKey(9)
	_, err = app.db.Exec(`INSERT INTO splits (link, receipt_id) VALUES(?, ?)`, key, id)
	if err != nil {
		app.logger.Error("Unable to store receipt in splits table", "error", err)
	}
	receipt.Link = key

	templates.ReceiptTable(receipt).Render(r.Context(), w)
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

// TODO: not the right way to store this info, chance of collision
func generateShortKey(keyLength int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}
