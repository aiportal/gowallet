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
const trace = true


func main() {
	number, _ := parseParams()

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
			seed, err = address.GenerateBrainWalletSeed("www.aiportal.net", "gowallet")
			if err != nil {
				println(err.Error())
				return
			}
		}

		accountKey, accountPub, err := address.GenerateHDAccount(seed[:], 0)
		if err != nil {
			println(err.Error())
			return
		}

		fmt.Printf("account key: %s\n", accountKey)
		fmt.Printf("account pub: %s\n", accountPub)

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
	} else {
		// Open wallets file.
	}
}

//Parse command line parameters
func parseParams() (number uint, export string) {

	//flag.StringVar(&vanity, "vanity", "", "Find vanity wallet address matching. (prefix or regular)")
	//flag.StringVar(&vanity, "v", "", "Find vanity wallet address matching. (prefix or regular)")

	flag.UintVar(&number, "number", 1, "Number of wallets to generate. (default 1)")
	flag.UintVar(&number, "n", 1, "Number of wallets to generate. (default 1)")

	flag.StringVar(&export, "export", "", "Export wallets in WIF format.")
	flag.StringVar(&export, "e", "", "Export wallets in WIF format.")

	flag.Parse()
	return
}
