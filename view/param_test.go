package view

import (
	"testing"
	"bytes"
)

func TestWalletParam(t *testing.T) {
	ts_data := [] struct {
		Secret string
		SecretBytes []byte
		Salt string
		SaltBytes []byte
	}{
		{
			Secret: `https://github.com/\x0aiport\x0al`,
			SecretBytes: append([]byte(`https://github.com/`), 0x0a, 'i', 'p', 'o', 'r', 't', 0x0a, 'l'),
			Salt: `gowall\x0et`,
			SaltBytes: append([]byte("gowall"), 0x0e, 't'),
		},
		{
			Secret: `www.aiportal.\x0n\x0et`,
			SecretBytes: append([]byte(`www.aiportal.\x0n`), 0x0e, 't'),
			Salt: `gow\x0all\x0et`,
			SaltBytes: []byte{'g', 'o', 'w', 0x0a, 'l', 'l', 0x0e, 't'},
		},
	}

	for _, v := range ts_data {
		wp := WalletParam { Secret:v.Secret, Salt:v.Salt }
		if !bytes.Equal(wp.SecretBytes(), v.SecretBytes) {
			t.Fatal("hex escape fail")
		}
		if !bytes.Equal(wp.SaltBytes(), v.SaltBytes) {
			t.Fatal("hex escape fail")
		}
	}
}

