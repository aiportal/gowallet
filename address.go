package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"github.com/btcsuite/btcutil/base58"
	"github.com/fatih/color"
	"github.com/njones/bitcoin-crypto/bitelliptic"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/ssh/terminal"
)

const brainWalletTip = `
The brain wallet uses a secret phrase and a salt phrase to generate the private key.

Secret phrase at least 16 characters, containing uppercase letters, lowercase letters, numbers, and special characters.
Salt phrase at least 6 characters.
Secret phrases and salt phrases allow the use of hexadecimal notation similar to ' \xff ' to represent a character.

It is advisable to use more complex secret phrases and to write secret phrases on paper.
It is also recommended that salt phrases be memorized in the brain.`

const debug = false


//Parse command line parameters
func parseCommandParams() (brain bool, output string) {

	flag.BoolVar(&brain, "brain", false, "Brain wallet mode.")
	flag.BoolVar(&brain, "b", false, "Brain wallet mode.")

	flag.StringVar(&output, "output", "", "Output file name.")
	flag.StringVar(&output, "o", "", "Output file name.")

	flag.Parse()
	return
}

func main() {
	brain, output := parseCommandParams()

	var seed []byte
	if brain == true {
		// Brain wallet
		secret, salt, err := inputBrainWalletSecret(brainWalletTip)
		if err != nil {
			return
		}
		seed, err = generateBrainWalletSeed(secret, salt)
		if err != nil {
			println(err.Error())
			return
		}
	} else {
		// Random wallet.
		var err error
		seed, err = generateRandomBytes(32)
		if err != nil {
			println(err.Error())
			return
		}
	}

	private_wif, address_wif, err := GenerateWalletWif(seed)
	if err != nil {
		println(err.Error())
		return
	}
	println("")
	println("private: " + private_wif)
	println("address: " + address_wif)

	if output != "" {
		ln := fmt.Sprintf("private: %s\naddress: %s", private_wif, address_wif)
		err := ioutil.WriteFile(output, []byte(ln), os.ModeAppend)
		if err != nil {
			println(err.Error())
		}
	}
}

// Generate wallet private key and address
func GenerateWalletWif(seed []byte) (privateWif string, addressWif string, err error) {
	reader := bytes.NewReader(seed)
	private_bytes, x, y, err := bitelliptic.S256().GenerateKey(reader)
	if err != nil {
		return
	}
	privateWif = base58.CheckEncode(private_bytes, 0x80)

	var public_bytes = [65]byte{0x04}
	copy(public_bytes[33 - len(x.Bytes()):], x.Bytes())
	copy(public_bytes[65 - len(y.Bytes()):], y.Bytes())

	public_sha := sha256.Sum256(public_bytes[:])

	ripeHash := ripemd160.New()
	ripeHash.Write(public_sha[:])
	public_ripe := ripeHash.Sum(nil)

	addressWif = base58.CheckEncode(public_ripe, 0x00)
	return
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
func inputBrainWalletSecret(tip string) (secret []byte, salt []byte, err error) {

	errInput := errors.New("Input error")

	// Tip
	color.Yellow(tip)
	println("")

	terminal.MakeRaw(int(os.Stdin.Fd()))
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
	if len(escapeHexString(salt1)) < 6 {
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
	if debug {
		fmt.Printf("secret: %X\n", secret)
		fmt.Printf("salt: %X\n", salt)
	}
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
func checkSecretStrength(secret string) bool {
	number, _ := regexp.MatchString("[0-9]+", secret)
	lower, _ := regexp.MatchString("[a-z]+", secret)
	upper, _ := regexp.MatchString("[A-Z]+", secret)
	special, _ := regexp.MatchString("[^0-9a-zA-Z ]", secret)
	return number && lower && upper && special
}

// Generate wallet seed from secret and salt
func generateBrainWalletSeed(secret []byte, salt []byte) (seed []byte, err error) {
	// WarpWallet encryption:
	// 1. s1 ← scrypt(key=passphrase||0x1, salt=salt||0x1, N=218, r=8, p=1, dkLen=32)
	// 2. s2 ← PBKDF2(key=passphrase||0x2, salt=salt||0x2, c=216, dkLen=32)
	// 3. private_key ← s1 ⊕ s2
	// 4. Generate public_key from private_key using standard Bitcoin EC crypto
	// 5. Output (private_key, public_key)

	if len(secret) == 0 {
		err = errors.New("Empty secret")
		return
	}
	if len(salt) < 0 {
		err = errors.New("Empty salt")
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

	s1, err := scrypt.Key(secret1, salt1, 16384, 8, 1, 32)
	if err != nil {
		return
	}
	s2 := pbkdf2.Key(secret2, salt2, 4096, 32, sha1.New)

	_, x1, y1, err := bitelliptic.S256().GenerateKey(bytes.NewReader(s1))
	if err != nil {
		return
	}
	_, x2, y2, err := bitelliptic.S256().GenerateKey(bytes.NewReader(s2))
	if err != nil {
		return
	}

	x, y := bitelliptic.S256().Add(x1, y1, x2, y2)

	seed = []byte{0x04}
	seed = append(seed, x.Bytes()...)
	seed = append(seed, y.Bytes()...)

	return
}
