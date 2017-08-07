package wallet

import (
	"crypto/sha256"
	"crypto/aes"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/scrypt"
)

// BIP38 encrypt implement
func EncryptKey(privateWif string, passKey []byte) (encryptWif string, err error) {
	// TODO: test compressed private WIF

	wif, err := btcutil.DecodeWIF(privateWif)
	if err != nil {
		return
	}
	wif.CompressPubKey = false
	pub, err := btcutil.NewAddressPubKey(wif.SerializePubKey(), &chaincfg.MainNetParams)
	if err != nil {
		return
	}
	sum := checkSum([]byte(pub.EncodeAddress()))
	key, err := scrypt.Key(passKey, sum, 16384, 8, 8, 64)
	if err != nil {
		return
	}

	data := encrypt(wif.PrivKey.Serialize(), key[:32], key[32:])
	buf := append([]byte{0x01, 0x42, 0xC0}, sum...)
	buf = append(buf, data...)
	buf = append(buf, checkSum(buf)...)
	encryptWif = base58.Encode(buf)

	return
}

func checkSum(buf []byte) []byte {
	h := sha256.Sum256(buf)
	h = sha256.Sum256(h[:])
	return h[:4]
}

func encrypt(pk, f1, f2 []byte) []byte {
	c, _ := aes.NewCipher(f2)

	for i, _ := range f1 {
		f1[i] ^= pk[i]
	}

	var dst = make([]byte, 48)
	c.Encrypt(dst, f1[:16])
	c.Encrypt(dst[16:], f1[16:])
	dst = dst[:32]

	return dst
}
