package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aiportal/gowallet/view"
	"github.com/aiportal/gowallet/wallet"
)

func main() {
	number, vanity, export := parseParams()
	if number > 0 {
		err := generateWallets(uint32(number), vanity, export)
		if err != nil {
			println(err.Error())
			return
		}
	} else {
		view.ShowSplashView(view.SplashStartView)

		var ws []*wallet.Wallet
		if !wallet.IsFileExists() {
			var err error
			ws, err = createWallets(1, 10)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			// save wallets
			wf := wallet.NewWalletFile(ws)
			wf.Save()
		} else {
			wf, err := wallet.LoadWalletFile()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			ws = wf.Wallets
		}

		showUI(ws)
	}
}

func showUI(ws []*wallet.Wallet) {

	accountView := view.NewAccountView(ws)
	accountView.Show()

	for accountView.Data != nil {
		cmd := accountView.Data.(string)
		if cmd == "quit" {
			break
		}
		tipView := view.NewTipView(cmd)
		if tipView != nil {
			tipView.Show()
		}
		accountView.Show()
	}
}

// create wallets by secret and salt
func createWallets(start, count uint32) (ws []*wallet.Wallet, err error) {
	view.ShowSplashView(view.SplashCreateView)

	// create wallets
	wp, err := view.InputNewParameters(3)
	if err != nil {
		return
	}
	//wp := view.WalletParam{Secret:"https://github.com/aiportal", Salt:"gowallet"}

	wa, err := wallet.NewWalletAccount(wp.SecretBytes(), wp.SaltBytes())
	if err != nil {
		return
	}
	ws, err = wa.GenerateWallets(start, count)
	if err != nil {
		return
	}
	return
}

//Parse command line parameters
func parseParams() (number uint, vanity, export string) {

	flag.UintVar(&number, "number", 0, "Number of wallets to generate.")
	flag.UintVar(&number, "n", 0, "Number of wallets to generate.")

	flag.StringVar(&vanity, "vanity", "", "Find vanity wallet address matching. (prefix)")
	flag.StringVar(&vanity, "v", "", "Find vanity wallet address matching. (prefix)")

	flag.StringVar(&export, "export", "", "Export wallets in WIF format.")
	flag.StringVar(&export, "e", "", "Export wallets in WIF format.")

	flag.Parse()
	return
}

func generateWallets(number uint32, vanity, export string) (err error) {

	view.ShowSplashView(view.SplashStartView)
	view.ShowSplashView(view.SplashCreateView)
	wp, err := view.InputNewParameters(3)
	if err != nil {
		return
	}
	wa, err := wallet.NewWalletAccount(wp.SecretBytes(), wp.SaltBytes())
	if err != nil {
		return
	}
	var ws []*wallet.Wallet
	if vanity == "" {
		ws, err = wa.GenerateWallets(0, uint32(number))
		if err != nil {
			return
		}
	} else {
		var patterns []string
		patterns, err = wa.NormalizeVanities([]string{vanity})
		if err != nil {
			return
		}
		ws, err = wa.FindVanities(patterns, func(i, c, n uint32) bool {
			fmt.Printf("progress: %d, %d, %d\n", i, c, n)
			return (n >= number)
		})
	}
	if export == "" {
		for _, w := range ws {
			fmt.Printf("wallet (%d): \n", w.No)
			fmt.Println("  " + w.Private)
			fmt.Println("  " + w.Address)
		}
	} else {
		var f *os.File
		f, err = os.Create(export)
		if err != nil {
			return
		}
		defer f.Close()
		for _, w := range ws {
			f.WriteString(fmt.Sprintf("wallet(%d): \r\n", w.No))
			f.WriteString(fmt.Sprintf("   private: %s\r\n", w.Private))
			f.WriteString(fmt.Sprintf("   address: %s\r\n", w.Address))
		}
	}
	return
}
