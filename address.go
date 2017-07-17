package main

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"crypto/rand"
	"os"
	"regexp"
	"syscall"
	"./secp256k1/bitecdsa"
	"./secp256k1/bitelliptic"
	"github.com/btcsuite/btcutil/base58"
	"github.com/fatih/color"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/ssh/terminal"
)

// WarpWallet encryption:
// 1. s1 ← scrypt(key=passphrase||0x1, salt=salt||0x1, N=218, r=8, p=1, dkLen=32)
// 2. s2 ← PBKDF2(key=passphrase||0x2, salt=salt||0x2, c=216, dkLen=32)
// 3. private_key ← s1 ⊕ s2
// 4. Generate public_key from private_key using standard Bitcoin EC crypto
// 5. Output (private_key, public_key)

//脑钱包使用一个秘密短语和一个盐短语生成私钥。
//秘密短语至少16个字符，包含大写字母，小写字母，数字和特殊字符。
//盐短语至少6个字符。
//建议使用较为复杂的秘密短语并将秘密短语记在纸上。
//同时建议将盐短语记在脑中。

const brainWalletTip = `
The brain wallet uses a secret phrase and a salt phrase to generate the private key.

Secret phrase at least 16 characters, containing uppercase letters, lowercase letters, numbers, and special characters.
Salt phrase at least 6 characters.

It is advisable to use more complex secret phrases and to write secret phrases on paper.
It is also recommended that salt phrases be memorized in the brain.

`

//Parse command line parameters
func parseCommandParams() (private string, brain bool, output string) {
	flag.StringVar(&private, "private", "", "Private key wif string for test.")

	flag.BoolVar(&brain, "brain", false, "Brain wallet mode.")
	flag.BoolVar(&brain, "b", false, "Brain wallet mode(shorthand).")

	flag.StringVar(&output, "output", "", "Output file name. (optional)")
	flag.StringVar(&output, "o", "", "Output file name(shorthand). (optional)")

	flag.Parse()
	return
}

func main() {
	private, brain, output := parseCommandParams()

	var private_key [32]byte
	if private == "" {
		if brain == true {
			// Brain wallet
			secret, salt, err := inputBrainWalletSecret()
			if err != nil {
				println(err.Error())
				return
			}
			//secret, salt = []byte("www.aiportal.net"), []byte("aiportal")
			private_key, err = generateBrainWalletKey(secret, salt)
			if err != nil {
				println(err.Error())
				return
			}
		} else {
			// Random private key.
			private_key_bytes, err := generateRandomBytes(32)
			if err == nil {
				copy(private_key[:], private_key_bytes)
			} else {
				println(err)
				return
			}
		}
	} else {
		// Private key from WIF string.
		private_key_bytes, _, _ := base58.CheckDecode(private)
		copy(private_key[:], private_key_bytes)
	}

	private_wif := base58.CheckEncode(private_key[:], 0x80)
	println("private: " + private_wif)

	public_key := computePublicKey(private_key)
	public_wif := base58.CheckEncode(public_key[:], 0x00)
	println("address: " + public_wif)

	if output != "" {
		err := ioutil.WriteFile(output, []byte(public_wif), os.ModeAppend)
		if err != nil {
			fmt.Printf("Failed to write to file. %s", err)
		}
	}
}

// Compute the public key from private key.
func computePublicKey(privateKey [32]byte) []byte {

	reader := bytes.NewReader(privateKey[:])
	key, err := bitecdsa.GenerateKey(bitelliptic.S256(), reader)
	if err != nil {
		println(err)
		return []byte{}
	}

	var public_key = [65]byte{0x04}
	x_bytes := key.X.Bytes()
	y_bytes := key.Y.Bytes()
	copy(public_key[33-len(x_bytes):], x_bytes)
	copy(public_key[65-len(y_bytes):], y_bytes)

	public_key_sha := sha256.Sum256(public_key[:])

	ripeHash := ripemd160.New()
	ripeHash.Write(public_key_sha[:])
	public_key_ripe := ripeHash.Sum(nil)

	return public_key_ripe[:]
}

//Generate secure random private key seed.
func generateRandomBytes(n int) ([]byte, error) {
	key := make([]byte, n)
	_, err := rand.Read(key[:])
	if err != nil {
		return nil, err
	}
	return key, nil
}

// Input secret and salt for brain wallet
func inputBrainWalletSecret() (secret []byte, salt []byte, err error) {

	// Tip
	color.Yellow(brainWalletTip)

	stdin := int(syscall.Stdin) //int(os.Stdin.Fd())
	errInput := errors.New("Input error")

	// Secret
	print("Brain wallet secret:")
	secret, err = terminal.ReadPassword(stdin)
	if err != nil {
		return
	}
	println("")
	if len(secret) < 16 {
		color.HiRed("Secret at least 16 characters")
		err = errInput
		return
	}
	if !checkSecretStrength(string(secret[:])) {
		color.HiRed("Secret should containing uppercase letters, lowercase letters, numbers, and special characters.")
		err = errInput
		return
	}

	// Secret again
	print("Enter secret again:")
	secret_again, err := terminal.ReadPassword(stdin)
	if err != nil {
		return
	}
	println("")
	if !bytes.Equal(secret, secret_again) {
		color.HiRed("Two input secret is different.")
		err = errInput
		return
	}

	// Salt
	print("Enter a salt phrase:")
	salt, err = terminal.ReadPassword(stdin)
	if err != nil {
		return
	}
	println("")
	if len(salt) < 6 {
		salt = []byte{}
		color.HiRed("Salt at least 6 characters.")
		err = errInput
		return
	}

	// Salt again
	print("Enter salt again:")
	salt_again, err := terminal.ReadPassword(stdin)
	if err != nil {
		return
	}
	println("")
	if !bytes.Equal(salt, salt_again) {
		color.HiRed("Two input salt is different.")
		err = errInput
		return
	}

	return
}

//Check secret strength
func checkSecretStrength(secret string) (valid bool) {
	number, _ := regexp.MatchString("[0-9]+", secret)
	lower, _ := regexp.MatchString("[a-z]+", secret)
	upper, _ := regexp.MatchString("[A-Z]+", secret)
	special, _ := regexp.MatchString("[^0-9a-zA-Z ]", secret)
	valid = number && lower && upper && special
	return
}

// Generate private key from secret and salt
func generateBrainWalletKey(secret []byte, salt []byte) (key [32]byte, err error) {

	if len(secret) == 0 || len(salt) < 0 {
		err = errors.New("empty secret or salt")
		return
	}

	secret1 := make([]byte, len(secret))
	secret2 := make([]byte, len(secret))
	for i, v := range secret {
		secret1[i] = v | 0x1
		secret2[i] = v | 0x2
	}
	salt1 := make([]byte, len(salt))
	salt2 := make([]byte, len(salt))
	for i, v := range salt {
		salt1[i] = v | 0x1
		salt2[i] = v | 0x2
	}
	key1, err := scrypt.Key(secret1, salt1, 16384, 8, 1, 32)
	if err != nil {
		return
	}
	key2 := pbkdf2.Key(secret2, salt2, 4096, 32, sha1.New)

	pk1, err := bitecdsa.GenerateKey(bitelliptic.S256(), bytes.NewReader(key1))
	if err != nil {
		return
	}
	pk2, err := bitecdsa.GenerateKey(bitelliptic.S256(), bytes.NewReader(key2))
	if err != nil {
		return
	}
	x, y := bitelliptic.S256().Add(pk1.X, pk1.Y, pk2.X, pk2.Y)

	key_bytes := []byte{0x04}
	key_bytes = append(key_bytes, x.Bytes()...)
	key_bytes = append(key_bytes, y.Bytes()...)
	key = sha256.Sum256(key_bytes[:])

	return
}
