package receipts

import "testing"

func TestX(t *testing.T) {
	res := 2 * 2
	if res != 4 {
		t.Fatalf("expected %d, received %d", 4, res)
	}
}