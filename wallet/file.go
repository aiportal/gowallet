package wallet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
)

const WalletFileName = ".gowallet.w02" // version 0.2

// check if user wallet file exists
func IsFileExists() bool {
	u, err := user.Current()
	if err != nil {
		return false
	}
	filename := fmt.Sprintf("%s/%s", u.HomeDir, WalletFileName)
	f, err := os.Stat(filename)
	if err != nil {
		return false
	}
	if f.IsDir() {
		return false
	}
	return true
}

type WalletFile struct {
	Wallets []*Wallet
}

func NewWalletFile(ws []*Wallet) *WalletFile {
	wf := new(WalletFile)
	wf.Wallets = make([]*Wallet, len(ws))
	for i, v := range ws {
		w := NewWallet(v.No, "", v.Address)
		wf.Wallets[i] = w
	}
	return wf
}

func LoadWalletFile() (wf *WalletFile, err error) {
	u, err := user.Current()
	if err != nil {
		return
	}
	path := fmt.Sprintf("%s/%s", u.HomeDir, WalletFileName)
	_, err = os.Stat(path)
	if err != nil {
		return
	}
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	wf = new(WalletFile)
	err = json.Unmarshal(bs, wf)
	if err != nil {
		return
	}
	return
}

func (wf *WalletFile) Save() (err error) {
	u, err := user.Current()
	if err != nil {
		return
	}
	path := fmt.Sprintf("%s/%s", u.HomeDir, WalletFileName)

	bs, err := json.Marshal(wf)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path, bs, os.ModeExclusive)
	if err != nil {
		return
	}
	return
}
//
//func (wf *WalletFile) Wallets() (ws []*Wallet) {
//	ws = make([]*Wallet, len(wf.wallets))
//	for i, v := range wf.wallets {
//		ws[i] = v
//	}
//	return
//}