package models

type Splits struct {
	ID   int
	Link string
}

type ReceiptBase struct {
	ID                   int     `json:"ID"`
	Link                 string  `json:"Link"`
	Subtotal             float64 `json:"Subtotal"`
	ServiceCharge        float64 `json:"ServiceCharge"`
	ServiceChargePercent int     `json:"ServiceChargePercent"`
	TaxPercent           int     `json:"TaxPercent"`
	TaxAmount            float64 `json:"TaxAmount"`
	Discount             float64 `json:"Discount"`
	DiscountPercent      int     `json:"DiscountPercent"`
	TotalAmount          float64 `json:"TotalAmount"`
}

type ReceiptItemBase struct {
	ID       int     `json:"ID"`
	Name     string  `json:"Name"`
	Quantity int     `json:"Quantity"`
	Price    float64 `json:"Price"`
}

type ReceiptItem struct {
	ReceiptItemBase
	PaidCount int `json:"PaidCount"`
}

type ReceiptViewItem struct {
	ReceiptItemBase
	FinalPrice    float64 `json:"FinalPrice"`
	TaxAmount     float64 `json:"TaxAmount"`
	ServiceCharge float64 `json:"ServiceCharge"`
	Paid          bool    `json:"Paid"`
}

type Receipt struct {
	ReceiptBase
	Items []ReceiptItem `json:"Items"`
}

type ReceiptView struct {
	ReceiptBase
	Items []ReceiptViewItem `json:"Items"`
}

type ReceiptOpenAI struct {
	Items []struct {
		Name     string  `json:"name"`
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
	} `json:"items"`
	ServiceCharge   float64 `json:"serviceCharge"` // this fields can be empty but openai does not support optional fields
	TaxPercent      int     `json:"taxPercent"`
	TaxAmount       float64 `json:"taxAmount"`
	Discount        float64 `json:"discount"`
	DiscountPercent int     `json:"discountPercent"`
	TotalAmount     float64 `json:"totalAmount"`
}

var MockReceipt = Receipt{
	ReceiptBase: ReceiptBase{
		ID:                   1,
		Link:                 "12121212kl",
		Subtotal:             29.15,
		ServiceCharge:        2.92,
		ServiceChargePercent: 10,
		TaxPercent:           6,
		TaxAmount:            1.92,
		TotalAmount:          31.99,
		Discount:             2,
	},
	Items: []ReceiptItem{
		{ReceiptItemBase: ReceiptItemBase{ID: 1, Name: "Item1", Quantity: 2, Price: 8.40}, PaidCount: 0},
		{ReceiptItemBase: ReceiptItemBase{ID: 2, Name: "Item2", Quantity: 1, Price: 5.75}, PaidCount: 1},
		{ReceiptItemBase: ReceiptItemBase{ID: 3, Name: "Item3", Quantity: 3, Price: 15.00}, PaidCount: 2},
	},
}

var MockReceiptView = ReceiptView{
	ReceiptBase: ReceiptBase{
		ID:                   1,
		Link:                 "12121212kl",
		Subtotal:             29.15,
		ServiceCharge:        2.92,
		ServiceChargePercent: 10,
		TaxPercent:           6,
		TaxAmount:            1.92,
		TotalAmount:          31.99,
		Discount:             2,
	},
	Items: []ReceiptViewItem{
		{ReceiptItemBase: ReceiptItemBase{ID: 1, Name: "Item1", Quantity: 1, Price: 8.40}, FinalPrice: 9.00, TaxAmount: 0.50, ServiceCharge: 0.84, Paid: false},
		{ReceiptItemBase: ReceiptItemBase{ID: 1, Name: "Item1", Quantity: 1, Price: 8.40}, FinalPrice: 9.00, TaxAmount: 0.50, ServiceCharge: 0.84, Paid: true},
		{ReceiptItemBase: ReceiptItemBase{ID: 2, Name: "Item2", Quantity: 1, Price: 5.75}, FinalPrice: 6.10, TaxAmount: 0.35, ServiceCharge: 0.58, Paid: true},
		{ReceiptItemBase: ReceiptItemBase{ID: 3, Name: "Item3", Quantity: 1, Price: 5.00}, FinalPrice: 5.30, TaxAmount: 0.30, ServiceCharge: 0.50, Paid: false},
	},
}
