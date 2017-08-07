package wallet

import (
	"fmt"
	"testing"
)

func TestAccountCreate(t *testing.T) {
	ts_data := []struct {
		Secret     string
		Salt       string
		PrivateKey string
		PublicKey  string
	}{
		{
			Secret:     "https://github.com/aiportal",
			Salt:       "gowallet",
			PrivateKey: "xprv9yADpbnN26FEeubBCZfUbF261DY8AtVZLxsDnzJCSZkcwP31Vc7ipHRozKLjHLp8LE2Pbs5K6LbS7AmnC1Mr5GsMDHBKmaabRXuq8J4ruVJ",
			PublicKey:  "xpub6C9aE7KFrToXsPfeJbCUxNxpZFNcaMDQiBnpbNhozuHbpBNA39RyN5kHqc4WPuc5Cx7WWnXLQ2WU99AdbZfUXbe3udiEwzLsDER85NK9Pa5",
		},
		{
			Secret:     "http://www.aiportal.net",
			Salt:       "gowallet",
			PrivateKey: "xprv9yF7QXAxyFqnrsxU7reHvRW1wGYSopJdyNSoar1FkTPJUrxAWiXP2GQz9HAgEJVJh5sC5Q662aycnbmqgeT3Ygyct954m3aVntV76WPNTYF",
			PublicKey:  "xpub6CETp2hrodQ65N2wDtBJHZSkVJNwDH2VLbNQPEQsJnvHMfHK4Fqda4jTzXykBD4Fonwf7RtAXUdZdpxCrTtc6zRmUAyMPYtNfFTBcYQYEdX",
		},
	}

	for _, v := range ts_data {
		wa, err := NewWalletAccount([]byte(v.Secret), []byte(v.Salt))
		if err != nil {
			t.Fatal(err)
		}
		if (wa.PrivateKey != v.PrivateKey) || (wa.PublicKey != v.PublicKey) {
			t.Fatal("account create error.")
		}
	}
}

// BIP32
// data from: https://bip32jp.github.io/english/
// uncompressed data from: https://www.bitaddress.org/bitaddress.org-v2.4-SHA1-1d5951f6a04dd5a287ac925da4e626870ee58d60.html
func TestAccountWallets(t *testing.T) {
	ts_data := []struct {
		MasterKey string
		Wallets   []Wallet
	}{
		{
			MasterKey: "xprv9yADpbnN26FEeubBCZfUbF261DY8AtVZLxsDnzJCSZkcwP31Vc7ipHRozKLjHLp8LE2Pbs5K6LbS7AmnC1Mr5GsMDHBKmaabRXuq8J4ruVJ",
			Wallets: []Wallet{
				{
					No:      0,
					Private: "Kx66Sn5RYLs5YkZGBhHjEdP9wAUEb9cFDGazrNNiiko2gJzsd6c4",
					Address: "1AS3QxQrNWSE35tqdR3QRcCzXRistNE2Za",
				},
				{
					No:      1,
					Private: "KzfnUwJiqNYrjVkCYTvWAtTjxM2jGeN25cw8SLdpZdcX2hTLfm93",
					Address: "14tPBNhoPV1N8FTyPJT4nyPcMKiUaVK1JL",
				},
				{
					No:      2,
					Private: "Kzn1RBUx5QtrgZsmpBxASAWioGpo4SGd8j3Ujy3wZe1r1WXBGiga",
					Address: "1EWJJViJfEAPECgzdXEUC3433gy6XBbQb8",
				},
				{
					No:      123769097,
					Private: "KwMsFJU9VugNPpnuZ9fyLAFCqR58e65k1jWMSsTT8tZNy2s5ohvd",
					Address: "1EtyNou5TPKbqwEkLBbsxEwuTyHk36TYZ2",
				},
				{
					No:      76909722,
					Private: "KxFQpy5gGg7fukpdYzDf8UtCRi9Vnm9aHDha5mn9xtf16TLPmSfA",
					Address: "1MZ5X7CXVUhUF6FYTe1svqbeXcA5jvCx6u",
				},
				//{
				// "KxvQRoZ1fkurEPsxnbGd8gSHQcZbKQVdMkdgw6pGNgzDAUZxwEtg"
				// "1Mh3kzqUd7Sp6MGckeDGGJNuSQTY91WWkT"
				//No: 999,
				//Private: "5JCeUphnTuYTZVHz3dK3J9Uw1pUmw1nueKwz5Z8Y8BRHoeiXVoh",
				//Address: "1DswXkvXkEcCsP1vwDS8B32VL2JQXtSZjZ",
				//},
			},
		},
		{
			MasterKey: "xprv9yF7QXAxyFqnrsxU7reHvRW1wGYSopJdyNSoar1FkTPJUrxAWiXP2GQz9HAgEJVJh5sC5Q662aycnbmqgeT3Ygyct954m3aVntV76WPNTYF",
			Wallets: []Wallet{
				{
					No:      0,
					Private: "L56WbZfTm5ayLzRLKii7KQNpXBfSGBPFAYTmaZ6t3rWxVS7s4XBB",
					Address: "18Qva5nNFuybzJMWU8si4fjznyADzRrKeX",
				},
				{
					No:      1,
					Private: "L2dfjPScTQ8uibT6Uu9NFMga6sD3uSZ1uWZDPem3JgvrH6srj9U4",
					Address: "1K4ogXK9jSo37kQkYzQARxpFuyPGhqiz9Q",
				},
				{
					No:      2,
					Private: "Kyy4yvABd1JasfSC7VzRtcqVU2jXV4Yq8j9aW9dUc5MSjhZ8AaUg",
					Address: "1FBvnG5igTvdADukFVYeNaiuWD6Wo2HEdq",
				},
				{
					No:      2376909,
					Private: "KxRXpZW5MYmqCJxpwizUUc5VKKRidY62oXbExy8kUFXzjMDxLMQn",
					Address: "1Dd88sLQiE9mPwuZTyvici8q26C64Vpcoj",
				},
				{
					No:      237690921,
					Private: "KwLx7zWwb7uLY9w8BrH35gpXnprYZi1KVHCuwfshMFr7dQm5ojSd",
					Address: "1LeVgU7oSkrPF53p9N9Zs7iHUfM94gewFb",
				},
			},
		},
	}

	for _, v := range ts_data {
		wa := new(WalletAccount)
		err := wa.generateAccount(v.MasterKey, 0)
		if err != nil {
			t.Fatal(err)
		}

		for _, tw := range v.Wallets {
			ws, err := wa.GenerateWallets(uint32(tw.No), 1)
			if err != nil {
				t.Fatal(err)
			}

			w := ws[0]
			if w.Private != tw.Private {
				fmt.Printf("wallet (%d): \n", w.No)
				fmt.Println("  " + w.Private)
				fmt.Println("  " + w.Address)
				t.Fatal("wallet private error")
			}
			if w.Address != tw.Address {
				fmt.Printf("wallet (%d): \n", w.No)
				fmt.Println("  " + w.Private)
				fmt.Println("  " + w.Address)
				t.Fatal("wallet address error")
			}
		}
	}
}
