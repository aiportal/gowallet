package address

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/fatih/color"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/ssh/terminal"
)

const hardened = 0x80000000

const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// Generate BIP44 account extended private key and extended public key.
func GenerateAccount(seed []byte, k uint32) (privateKey string, publicKey string, err error) {
	master_key, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return
	}
	purpose, err := master_key.Child(hardened + 44)
	if err != nil {
		return
	}
	coin_type, err := purpose.Child(hardened + 0)
	if err != nil {
		return
	}
	account, err := coin_type.Child(hardened + k)
	if err != nil {
		return
	}

	account_pub, err := account.Neuter()
	if err != nil {
		return
	}

	privateKey = account.String()
	publicKey = account_pub.String()
	return
}

// Generate multiple address
func GenerateWallets(account string, count uint32) (wallets [][]string, err error) {
	account_key, err := hdkeychain.NewKeyFromString(account)
	if err != nil {
		return
	}
	change, err := account_key.Child(0)
	if err != nil {
		return
	}

	wallets = make([][]string, count)
	for i := uint32(0); i < count; i++ {
		child, err1 := change.Child(i)
		if err1 != nil {
			err = err1
			break
		}
		private_key, err1 := child.ECPrivKey()
		if err1 != nil {
			err = err1
			return
		}
		private_wif, err1 := btcutil.NewWIF(private_key, &chaincfg.MainNetParams, false)
		if err1 != nil {
			err = err1
			return
		}
		private := private_wif.String()
		address, err1 := child.Address(&chaincfg.MainNetParams)
		if err1 != nil {
			err = err1
			return
		}
		wallets[i] = []string{private, address.String()}
	}
	return
}

// Find vanity address
func SearchVanities(account string, vanity string, count uint32,
	progress func(i uint32, count uint32, n uint32)) (wallets [][]string, err error) {

	// check vanity
	if len(vanity) > 6 {
		err = errors.New("Vanity maximum 6 characters.")
		return
	}
	for _, c := range vanity {
		if !strings.Contains(alphabet, string(c)) {
			err = errors.New("Invalid vanity character: " + string(c))
			return
		}
	}
	pattern := "1" + vanity

	account_key, err := hdkeychain.NewKeyFromString(account)
	if err != nil {
		return
	}
	change, err := account_key.Child(0)
	if err != nil {
		return
	}

	wallets = [][]string{}
	for i := uint32(0); i < hardened; i++ {
		if i > 0 && i%10000 == 0 {
			progress(i, hardened, uint32(len(wallets)))
		}

		child, err1 := change.Child(i)
		if err1 != nil {
			err = err1
			break
		}
		address, err1 := child.Address(&chaincfg.MainNetParams)
		if err1 != nil {
			err = err1
			return
		}

		address_wif := address.String()
		if strings.HasPrefix(address_wif, pattern) {
			private_key, err1 := child.ECPrivKey()
			if err1 != nil {
				err = err1
				return
			}
			private_wif, err1 := btcutil.NewWIF(private_key, &chaincfg.MainNetParams, false)
			if err1 != nil {
				err = err1
				return
			}
			private := private_wif.String()

			wallet := []string{private, address.String(), fmt.Sprintf("%d", i)}
			wallets = append(wallets, wallet)
			if len(wallets) == int(count) {
				return
			}
		}
	}
	if len(wallets) < 1 {
		err = errors.New("Vanity pattern not found.")
	}
	return
}

// Input secret and salt for brain wallet
func InputBrainWalletSecret(tip string) (secret string, salt string, err error) {
	// TODO: if input error, repeat input.

	errInput := errors.New("Input error")

	// Tip
	color.Yellow(tip)
	println("")

	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return
	}
	t := terminal.NewTerminal(os.Stdin, "")
	defer terminal.Restore(int(os.Stdin.Fd()), oldState)

	// Secret
	print("Brain wallet secret:")
	secret1, err := t.ReadPassword("")
	if err != nil {
		return
	}
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
	print("Enter a salt:")
	salt1, err := t.ReadPassword("")
	if err != nil {
		return
	}
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

	return secret1, salt1, err
}

// Converts a string like "\xF0" or "\x0f" into a byte
func escapeHexString(str string) []byte {

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
func checkSecretStrength(secret string) bool {
	number, _ := regexp.MatchString("[0-9]+", secret)
	lower, _ := regexp.MatchString("[a-z]+", secret)
	upper, _ := regexp.MatchString("[A-Z]+", secret)
	special, _ := regexp.MatchString("[^0-9a-zA-Z ]", secret)
	return number && lower && upper && special
}

// Generate wallet seed from secret and salt
func GenerateBrainWalletSeed(secret string, salt string) (seed []byte, err error) {
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
	secret_bytes := escapeHexString(secret)
	salt_bytes := escapeHexString(salt)

	secret1 := make([]byte, len(secret_bytes))
	secret2 := make([]byte, len(secret_bytes))
	for i, v := range secret_bytes {
		secret1[i] = byte(v | 0x01)
		secret2[i] = byte(v | 0x02)
	}

	salt1 := make([]byte, len(salt_bytes))
	salt2 := make([]byte, len(salt_bytes))
	for i, v := range salt_bytes {
		salt1[i] = byte(v | 0x01)
		salt2[i] = byte(v | 0x02)
	}

	s1, err := scrypt.Key(secret1, salt1, 16384, 8, 1, 32)
	if err != nil {
		return
	}
	s2 := pbkdf2.Key(secret2, salt2, 4096, 32, sha1.New)

	pk1, _ := btcec.PrivKeyFromBytes(btcec.S256(), s1)
	pk2, _ := btcec.PrivKeyFromBytes(btcec.S256(), s2)
	x, y := btcec.S256().Add(pk1.X, pk1.Y, pk2.X, pk2.Y)

	seed = []byte{0x04}
	seed = append(seed, x.Bytes()...)
	seed = append(seed, y.Bytes()...)

	seed_hash := sha256.Sum256(seed)
	seed = seed_hash[:]
	return
}
