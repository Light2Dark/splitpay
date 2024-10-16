package models

type Splits struct {
	ID   int
	Link string
}

type ReceiptOpenAI struct {
	Items []struct {
		Name     string  `json:"name"`
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
	} `json:"items"`
	ServiceCharge float64 `json:"serviceCharge"` // this fields can be empty but openai does not support optional fields
	TaxPercent    int     `json:"taxPercent"`
	TaxAmount     float64 `json:"taxAmount"`
	TotalAmount   float64 `json:"totalAmount"`
}

type ReceiptItem struct {
	ID       int     `json:"ID"`
	Name     string  `json:"Name"`
	Quantity int     `json:"Quantity"`
	Price    float64 `json:"Price"`
}

type Receipt struct {
	ID            int           `json:"ID"`
	Link          string        `json:"Link"`
	Items         []ReceiptItem `json:"Items"`
	Subtotal      float64       `json:"Subtotal"`
	ServiceCharge float64       `json:"ServiceCharge"`
	TaxPercent    int           `json:"TaxPercent"`
	TaxAmount     float64       `json:"TaxAmount"`
	TotalAmount   float64       `json:"TotalAmount"`
}

var MockReceipt = Receipt{
	ID:   1, // test data
	Link: "1212klkml",
	Items: []ReceiptItem{
		{ID: 1, Name: "Item1", Quantity: 2, Price: 8.40},
		{ID: 2, Name: "Item2", Quantity: 1, Price: 5.75},
		{ID: 3, Name: "Item3", Quantity: 3, Price: 15.00},
	},
	Subtotal:      29.15,
	ServiceCharge: 2.92,
	TaxPercent:    6,
	TaxAmount:     1.75,
	TotalAmount:   33.82,
}
