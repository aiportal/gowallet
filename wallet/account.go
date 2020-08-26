package wallet

import (
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

const hardened = 0x80000000

const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

var AddressNetParams = chaincfg.MainNetParams

type WalletAccount struct {
	PrivateKey string
	PublicKey  string
}

func NewWalletAccount(secret, salt []byte) (wa *WalletAccount, err error) {
	seed, err := generateSeed(secret, salt)
	if err != nil {
		return wa, err
	}
	return NewWalletAccountFromSeed(seed)
}

func NewWalletAccountFromSeed(seed []byte) (wa *WalletAccount, err error) {

	wa = new(WalletAccount)
	master_key, err := hdkeychain.NewMaster(seed, &AddressNetParams)
	if err != nil {
		return
	}
	err = wa.generateAccount(master_key.String(), 0)
	return
}

// Generate wallet seed from secret and salt
func generateSeed(secret, salt []byte) (seed []byte, err error) {
	// WarpWallet encryption:
	// 1. s1 ← scrypt(key=passphrase||0x1, salt=salt||0x1, N=218, r=8, p=1, dkLen=32)
	// 2. s2 ← PBKDF2(key=passphrase||0x2, salt=salt||0x2, c=216, dkLen=32)
	// 3. private_key ← s1 ⊕ s2
	// 4. Generate public_key from private_key using standard Bitcoin EC crypto
	// 5. Output (private_key, public_key)

	if len(secret) == 0 {
		err = errors.New("empty secret")
		return
	}
	if len(salt) == 0 {
		err = errors.New("empty salt")
		return
	}

	secret1 := make([]byte, len(secret))
	secret2 := make([]byte, len(secret))
	for i, v := range secret {
		secret1[i] = byte(v | 0x01)
		secret2[i] = byte(v | 0x02)
	}

	salt1 := make([]byte, len(salt))
	salt2 := make([]byte, len(salt))
	for i, v := range salt {
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

// Generate BIP44 account extended private key and extended public key.
func (wa *WalletAccount) generateAccount(masterKey string, k uint32) (err error) {
	master_key, err := hdkeychain.NewKeyFromString(masterKey)
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

	wa.PrivateKey = account.String()
	wa.PublicKey = account_pub.String()
	return
}

// Generate multiple address
func (wa *WalletAccount) GenerateWallets(start, count uint32) (wallets []*Wallet, err error) {

	var account, change *hdkeychain.ExtendedKey
	account, err = hdkeychain.NewKeyFromString(wa.PrivateKey)
	if err != nil {
		return
	}
	change, err = account.Child(0)
	if err != nil {
		return
	}

	wallets = make([]*Wallet, count)
	for i := uint32(0); i < count; i++ {
		n := start + i
		child, err := change.Child(n)
		if err != nil {
			break
		}
		private_key, err := child.ECPrivKey()
		if err != nil {
			break
		}
		private_wif, err := btcutil.NewWIF(private_key, &AddressNetParams, true)
		if err != nil {
			break
		}
		address_key, err := child.Address(&AddressNetParams)
		if err != nil {
			break
		}

		private_str := private_wif.String()
		address_str := address_key.String()
		wallets[i] = NewWallet(n, private_str, address_str)
	}
	return
}

func (wa *WalletAccount) NormalizeVanities(vanities []string) (patterns []string, err error) {
	// check vanity
	for _, v := range vanities {
		if len(v) > 6 {
			err = errors.New("Vanity maximum 6 characters.")
			return
		}
		for _, c := range v {
			if !strings.Contains(alphabet, string(c)) {
				err = errors.New("Invalid vanity character: " + string(c))
				return
			}
		}
	}

	patterns = make([]string, len(vanities))
	for i, v := range vanities {
		patterns[i] = "1" + v
	}
	return
}

type FindProgress func(progress, count, found uint32) (stop bool)

// Find vanity address
func (wa *WalletAccount) FindVanities(patterns []string, progress FindProgress) (ws []*Wallet, err error) {

	account, err := hdkeychain.NewKeyFromString(wa.PrivateKey)
	if err != nil {
		return
	}
	change, err := account.Child(0)
	if err != nil {
		return
	}

	for i := uint32(0); i < hardened; i++ {

		if i%100000 == 0 {
			if progress != nil {
				stop := progress(i, hardened, uint32(len(ws)))
				if stop {
					break
				}
			}
		}

		child, err := change.Child(i)
		if err != nil {
			break
		}
		address_key, err := child.Address(&AddressNetParams)
		if err != nil {
			break
		}
		address_str := address_key.String()

		// check patterns
		match := false
		for _, p := range patterns {
			match = strings.EqualFold(p, address_str[:len(p)])
			if match {
				break
			}
		}

		if match {
			var w *Wallet
			w, err = wa.createWallet(child, i)
			if err != nil {
				break
			}
			ws = append(ws, w)
		}
	}
	if len(ws) == 0 {
		err = errors.New("vanities not found.")
	}
	return
}

func (wa *WalletAccount) createWallet(child *hdkeychain.ExtendedKey, seqNum uint32) (w *Wallet, err error) {

	private_key, err := child.ECPrivKey()
	if err != nil {
		return
	}
	private_wif, err := btcutil.NewWIF(private_key, &AddressNetParams, true)
	if err != nil {
		err = err
		return
	}
	private_str := private_wif.String()

	address_key, err := child.Address(&AddressNetParams)
	if err != nil {
		return
	}
	address_str := address_key.String()

	w = NewWallet(seqNum, private_str, address_str)
	return
}
