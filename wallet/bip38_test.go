package wallet

import "testing"

func TestEncryptKey(t *testing.T) {

	testEncrypt := []struct {
		label string
		pass string
		wif string
		secret string
	}{
		{
			label: "normal",
			pass : "TestingOneTwoThree",
			wif : "5KN7MzqK5wt2TP1fQCYyHBtDrXdJuXbUzm4A9rKAteGu3Qi5CVR",
			secret : "6PRVWUbkzzsbcVac2qwfssoUJAN1Xhrg6bNk8J7Nzm5H7kxEbn2Nh2ZoGg",
		},
		{
			label: "normal",
			pass: "Satoshi",
			wif: "5HtasZ6ofTHP6HCwTqTkLDuLQisYPah7aUnSKfC7h4hMUVw2gi5",
			secret: "6PRNFFkZc2NZ6dJqFfhRoFNMR9Lnyj7dYGrzdgXXVMXcxoKTePPX1dWByq",
		},
		//{
		//	label: "compressed",
		//	pass: "TestingOneTwoThree",
		//	wif: "L44B5gGEpqEDRS9vVPz7QT35jcBG2r3CZwSwQ4fCewXAhAhqGVpP",
		//	secret: "6PYNKZ1EAgYgmQfmNVamxyXVWHzK5s6DGhwP4J5o44cvXdoY7sRzhtpUeo",
		//},
		//{
		//	label: "compressed",
		//	pass: "Satoshi",
		//	wif: "KwYgW8gcxj1JWJXhPSu4Fqwzfhp5Yfi42mdYmMa4XqK7NJxXUSK7",
		//	secret: "6PYLtMnXvfG3oJde97zRyLYFZCYizPU5T3LwgdYJz1fRhh16bU7u6PPmY7",
		//},
	}

	for _, v := range testEncrypt {
		encrypt, err := EncryptKey(v.wif, []byte(v.pass))
		if err != nil {
			t.Fatal(err)
		}
		if (encrypt != v.secret) {
			t.Fatal(v.label)
		}
	}
}
