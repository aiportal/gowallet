package wallet

import (
	"testing"
	"fmt"
)

func TestBalanceFmt(t *testing.T) {
	ts_data := [] struct {
		Satoshi uint64
		FmtBtc  string
	}{
		{Satoshi: 0, FmtBtc: "0.000 BTC"},
		{Satoshi: 6*100*1000*1000, FmtBtc: "6.000 BTC"},
		{Satoshi: 1*100*1000*1000 + 18*100*1000, FmtBtc: "1.018 BTC"},
		{Satoshi: 3*100*1000*1000 + 88*100, FmtBtc: "3.000088 BTC"},
		{Satoshi: 9*100*1000*1000 + 88, FmtBtc: "9.00000088 BTC"},
		{Satoshi: 18*100*1000, FmtBtc: "0.018 BTC"},
		{Satoshi: 88*100, FmtBtc: "0.000088 BTC"},
		{Satoshi: 88, FmtBtc: "0.00000088 BTC"},
		{Satoshi: 69259423, FmtBtc: "0.69259423 BTC"},
	}

	for _, v := range ts_data {
		wb := NewBalance(v.Satoshi)
		if wb.FmtBtc() != v.FmtBtc {
			fmt.Printf("satoshi: %d, %q", v.Satoshi, wb.FmtBtc())
			t.Fatal("fomat error")
		}
	}
}
