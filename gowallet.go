package main

import (
	"flag"
	"fmt"
	"os"
	"gowallet/address"
)

const goWalletTip = `
GoWallet uses a secret phrase and a password phrase to generate your safe wallets.
Project location: https://github.com/aiportal/gowallet

Secret at least 16 characters, containing uppercase letters, lowercase letters, numbers, and special characters.
Password at least 8 characters.
Secret and password allow the use of hexadecimal notation similar to '\xff' or '\xFF' to represent a character.

It is advisable to use more complex secret and to write secret on paper.
It is also recommended that password be memorized in the brain.`

const debug = false
const trace = false


func main() {
	vanity, number, export := parseParams()

	if _, err := os.Stat("./gowallet.wlt"); os.IsNotExist(err) {
		// New wallets.
		var seed []byte
		if !debug {
			secret, pwd, err := address.InputBrainWalletSecret(goWalletTip)
			if err != nil {
				println(err.Error())
				return
			}
			if trace {
				println("your secret is: " + secret)
				println("your password is: " + pwd)
			}
			seed, err = address.GenerateBrainWalletSeed(secret, pwd)
			if err != nil {
				println(err.Error())
				return
			}
		} else {
			seed, err = address.GenerateBrainWalletSeed("https://github.com/aiportal", "gowallet")
			if err != nil {
				println(err.Error())
				return
			}
		}

		accountKey, accountPub, err := address.GenerateAccount(seed[:], 0)
		if err != nil {
			println(err.Error())
			return
		}
		fmt.Printf("account extended key: %s\n", accountKey)
		fmt.Printf("account extended pub: %s\n", accountPub)

		if vanity == "" {
			wallets, err := address.GenerateWallets(accountKey, uint32(number))
			if err != nil {
				println(err.Error())
				return
			}
			for i, w := range wallets {
				fmt.Printf("wallet(%d): \n", i)
				fmt.Printf("	private: %s\n", w[0])
				fmt.Printf("	address: %s\n", w[1])
			}
			if export != "" {
				err := exportWallets(export, wallets)
				if err != nil {
					println(err.Error())
					return
				}
			}
		} else {
			wallets, err := address.SearchVanities(accountKey, vanity, uint32(number),
				func(i uint32, count uint32, n uint32) {
					fmt.Printf("processed：%d / %d, found: %d \n", i, count, n)
			})
			if err != nil {
				println(err.Error())
				return
			}
			for i, w := range wallets {
				fmt.Printf("wallet(%d): \n", i)
				fmt.Printf("	private: %s\n", w[0])
				fmt.Printf("	address: %s\n", w[1])
			}
			if export != "" {
				err := exportWallets(export, wallets)
				if err != nil {
					println(err.Error())
					return
				}
			}
		}
	} else {
		// Open wallets file.
	}
}

//Parse command line parameters
func parseParams() (vanity string, number uint, export string) {

	flag.StringVar(&vanity, "vanity", "", "Find vanity wallet address matching. (prefix or regular)")
	flag.StringVar(&vanity, "v", "", "Find vanity wallet address matching. (prefix or regular)")

	flag.UintVar(&number, "number", 1, "Number of wallets to generate. (default 1)")
	flag.UintVar(&number, "n", 1, "Number of wallets to generate. (default 1)")

	flag.StringVar(&export, "export", "", "Export wallets in WIF format.")
	flag.StringVar(&export, "e", "", "Export wallets in WIF format.")

	flag.Parse()
	return
}

// Export wallets
func exportWallets(filename string, wallets [][]string) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return
	}
	defer f.Close()
	for i, w := range wallets {
		f.WriteString(fmt.Sprintf("wallet(%d): \n", i))
		f.WriteString(fmt.Sprintf("   private: %s", w[0]))
		f.WriteString(fmt.Sprintf("   address: %s", w[1]))
	}
	return
}
