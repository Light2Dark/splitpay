package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

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

	var receipt models.Receipt
	schema, err := jsonschema.GenerateSchemaForType(receipt)
	if err != nil {
		app.logError(w, r, "GenerateSchemaForType error", err)
		return
	}

	resp, err := app.openai.CreateChatCompletion(
		r.Context(),
		openai.ChatCompletionRequest{
			Model:     openai.GPT4oMini,
			MaxTokens: 1000,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleUser,
					MultiContent: []openai.ChatMessagePart{
						{
							Type: openai.ChatMessagePartTypeText,
							Text: "This is a receipt for a meal that a group of people had. I would like you to return a list of items that were ordered, alongside their price and quantity. Sometimes, an item will have descriptions or additions, ignore these, just take the main item and any additional charges if there is any. Also, return the service charge and sales tax percentage that was charged. If there is nothing, return an empty string.",
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

	err = schema.Unmarshal(resp.Choices[0].Message.Content, &receipt)
	if err != nil {
		app.logError(w, r, "Unmarschal schema error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&healthResponse{Status: fmt.Sprintf("response: %v", receipt)})
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
