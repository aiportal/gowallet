package view

import (
	"encoding/hex"
	"errors"
	"os"
	"regexp"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	createWalletTip = `
  Creating wallet...

`
	createSecretPrompt = "Enter secret phrase:"
	createSecretAgain  = "Enter secret again:"
	createSaltePrompt  = "Enter salt phrase:"
	createSalteAgain   = "Enter salt again:"

	createErrSecretLength    = "Secret at least 16 characters."
	createErrSaltLength      = "Salt at least 16 characters."
	createErrSecretStrength  = "Secret should containing uppercase letters, lowercase letters, numbers, and special characters."
	createErrInconsistencies = "Two input inconsistencies"
)

var ErrInput = errors.New("input error")

type WalletParam struct {
	Secret string
	Salt   string
}

func (wp *WalletParam) SecretBytes() []byte {
	return wp.escapeHexString(wp.Secret)
}

func (wp *WalletParam) SaltBytes() []byte{
	return wp.escapeHexString(wp.Salt)
}

func InputNewParameters(chance uint32) (wp *WalletParam, err error) {

	color.HiYellow(createWalletTip)

	wp = new(WalletParam)
	for i := uint32(0); i < chance; i++ {
		if err = wp.inputSecret(); err == nil {
			break
		}
	}
	if err != nil {
		return
	}
	for i := uint32(0); i < chance; i++ {
		if err = wp.inputSalt(); err == nil {
			break
		}
	}
	return
}

func (wp *WalletParam) inputSecret() (err error) {

	state, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return
	}
	t := terminal.NewTerminal(os.Stdin, "")
	defer terminal.Restore(int(os.Stdin.Fd()), state)

	// Secret
	print(createSecretPrompt)
	secret1, err := t.ReadPassword("")
	if err != nil {
		return
	}
	println("")
	if len(wp.escapeHexString(secret1)) < 16 {
		color.HiRed(createErrSecretLength)
		err = ErrInput
		return
	}
	if !wp.checkSecretStrength(secret1) {
		color.HiRed(createErrSecretStrength)
		err = ErrInput
		return
	}

	// Secret again
	print(createSecretAgain)
	secret2, err := t.ReadPassword("")
	if err != nil {
		return
	}
	println("")
	if secret1 != secret2 {
		color.HiRed(createErrInconsistencies)
		err = ErrInput
		return
	}

	wp.Secret = secret1
	return
}

func (wp *WalletParam) inputSalt() (err error) {

	state, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return
	}
	t := terminal.NewTerminal(os.Stdin, "")
	defer terminal.Restore(int(os.Stdin.Fd()), state)

	// Salt
	print(createSaltePrompt)
	salt1, err := t.ReadPassword("")
	if err != nil {
		return
	}
	println("")
	if len(wp.escapeHexString(salt1)) < 6 {
		color.HiRed(createErrSaltLength)
		err = ErrInput
		return
	}

	// Salt again
	print(createSalteAgain)
	salt2, err := t.ReadPassword("")
	if err != nil {
		return
	}
	println("")
	if salt1 != salt2 {
		color.HiRed(createErrInconsistencies)
		err = ErrInput
		return
	}

	wp.Salt = salt1
	return
}

// Converts a string like "\xF0" or "\x0f" into a byte
func (wp *WalletParam) escapeHexString(str string) []byte {

	r, _ := regexp.Compile("\\\\x[0-9A-Fa-f]{2}")
	exists := r.MatchString(str)
	if !exists {
		return []byte(str)
	}

	key := r.ReplaceAllFunc([]byte(str), func(s []byte) []byte {
		v, _ := hex.DecodeString(string(s[2:]))
		return v
	})

	return key
}

//Check secret strength
func (wp *WalletParam) checkSecretStrength(secret string) bool {
	number, _ := regexp.MatchString("[0-9]+", secret)
	lower, _ := regexp.MatchString("[a-z]+", secret)
	upper, _ := regexp.MatchString("[A-Z]+", secret)
	special, _ := regexp.MatchString("[^0-9a-zA-Z ]", secret)
	return number && lower && upper && special
}
