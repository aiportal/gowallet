package wallet

import (
	"fmt"
)

type WalletBalance struct {
	Satoshi uint64
}

func NewBalance(satoshi uint64) *WalletBalance {
	b := new(WalletBalance)
	b.Satoshi = satoshi
	return b
}

func (wb *WalletBalance) Value() uint64 {
	return wb.Satoshi
}

func (wb *WalletBalance) FmtBtc() string {
	satoshi := wb.Satoshi

	if satoshi == 0 {
		return "0.000 BTC"
	}
	btc := float64(satoshi) / (100 * 1000 * 1000)
	if satoshi%(100*1000) == 0 {
		return fmt.Sprintf("%.3f BTC", btc)
	}
	if satoshi%100 == 0 {
		return fmt.Sprintf("%.6f BTC", btc)
	}
	return fmt.Sprintf("%.8f BTC", btc)
}
