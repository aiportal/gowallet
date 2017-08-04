package wallet

import (
	"testing"
	"fmt"
)

func TestWallet(t *testing.T) {
	ts_data := []struct {
		Address string
		Balance uint64
	}{
		{
			Address: "1DogeKd9JrUNzFaLEyWAVxCVXSvWxe6sAm",
			Balance: 69259423,
		},
		{
			Address: "1KtWutb75LqXrAd4BkcW2hqG7SWab2xJeB",
			Balance: 0,
		},
	}

	for _, v := range ts_data {
		w := Wallet{Address: v.Address}
		b, err := w.Balance()
		if err != nil {
			t.Fatal(err)
		}
		if b != v.Balance {
			t.Fatal(fmt.Sprintf("Balance: %d != %d", b, v.Balance))
		}
	}
}
