package main

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"flag"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/ssh/terminal"
	"crypto/rand"
	"regexp"

	"./secp256k1/bitecdsa"
	"./secp256k1/bitelliptic"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"encoding/hex"
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
It is also recommended that salt phrases be memorized in the brain.`

const debug = true


//Parse command line parameters
func parseCommandParams() (private string, brain bool, output string) {
	flag.StringVar(&private, "private", "", "Private key wif string for test.")

	flag.BoolVar(&brain, "brain", false, "Brain wallet mode.")
	flag.BoolVar(&brain, "b", false, "...")

	flag.StringVar(&output, "output", "", "Output file name. (optional)")
	flag.StringVar(&output, "o", "", "...")

	flag.Parse()
	return
}

func main() {
	private, brain, output := parseCommandParams()

	var private_key [32]byte
	if private == "" {
		if brain == true {
			// Brain wallet
			secret, salt, err := inputBrainWalletSecret(brainWalletTip)
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
			random_bytes := generateRandomBytes(32)
			private_key = sha256.Sum256(random_bytes)
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

// Generate random private key seed.
func generateRandomBytes(n int) []byte {
	buf := make([]byte, n)
	rand.Read(buf)
	return buf
}

// Input secret and salt for brain wallet
func inputBrainWalletSecret(tip string) (secret []byte, salt []byte, err error) {

	errInput := errors.New("Input error")

	// Tip
	color.Yellow(tip)
	println("")

	t := terminal.NewTerminal(os.Stdin, "")

	// Secret
	print("Brain wallet secret:")
	secret1, err := t.ReadPassword("")
	if err != nil {
		return
	}
	if debug { print(secret1) }
	println("")
	if len(escapeHexString(secret1)) < 16 {
		color.HiRed("  Secret at least 16 characters")
		err = errInput
		return
	}
	if !checkSecretStrength(secret1) {
		color.HiRed("  Secret should containing uppercase letters, lowercase letters, numbers, and special characters.")
		err = errInput
		return
	}

	// Secret again
	print("Enter secret again:")
	secret2, err := t.ReadPassword("")
	if err != nil {
		return
	}
	println("")
	if secret1 != secret2 {
		color.HiRed("  Two input secret is different.")
		err = errInput
		return
	}

	// Salt
	print("Enter a salt phrase:")
	salt1, err := t.ReadPassword("")
	if err != nil {
		return
	}
	if debug { print(salt1) }
	println("")
	if len(salt1) < 6 {
		color.HiRed("  Salt at least 6 characters.")
		err = errInput
		return
	}

	// Salt again
	print("Enter salt again:")
	salt2, err := t.ReadPassword("")
	if err != nil {
		return
	}
	println("")
	if salt1 != salt2 {
		color.HiRed("  Two input salt is different.")
		err = errInput
		return
	}

	secret = escapeHexString(secret1)
	salt = escapeHexString(salt1)
	return
}

// Converts a string like "\xF0" or "\x0f" into a byte
func escapeHexString(str string) []byte {

	r, _ := regexp.Compile("\\\\x[0-9A-Fa-f]{2}")
	exists := r.MatchString(str)
	if !exists {
		return []byte(str)
	}

	key := r.ReplaceAllFunc([]byte(str), func(s []byte) []byte{
		v, _ := hex.DecodeString(string(s[2:]))
		return v
	})

	return key
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
