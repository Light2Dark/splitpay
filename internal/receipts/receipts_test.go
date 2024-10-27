package receipts

import (
	"testing"

	"github.com/Light2Dark/splitpay/models"
)

func TestX(t *testing.T) {
	res := 2 * 2
	if res != 4 {
		t.Fatalf("expected %d, received %d", 4, res)
	}
}

var receiptConversionTest = []struct{
	OpenAIReceipt models.ReceiptOpenAI
	Receipt models.Receipt
}{
	{
		models.ReceiptOpenAI{},
		models.Receipt{},
	},
}