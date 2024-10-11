package models

type Splits struct {
	ID   int
	Link string
}

type Receipt struct {
	Items []struct {
		Name     string  `json:"name"`
		Quantity string  `json:"quantity"`
		Price    float32 `json:"price"`
	} `json:"items"`
	ServiceCharge string  `json:"serviceCharge"` // this fields can be empty but openai does not support optional fields
	Tax           string  `json:"tax"`
	TotalAmount   float32 `json:"totalAmount"`
}

var MockReceipt = Receipt{
	Items: []struct {
		Name     string  `json:"name"`
		Quantity string  `json:"quantity"`
		Price    float32 `json:"price"`
	}{
		{Name: "Item1", Quantity: "2", Price: 10.50},
		{Name: "Item2", Quantity: "1", Price: 5.75},
		{Name: "Item3", Quantity: "3", Price: 7.20},
	},
	ServiceCharge: "2.50",
	Tax:           "1.75",
	TotalAmount:   45.90,
}
