package safe

import (
	"testing"
)

func TestEncryptAES_DecryptAES(t *testing.T) {

	d := "just_test"
	key := "12345678abcdefgh"

	t.Logf("src before encrypt: %s", string(d))

	x1, err := EncryptAES(d, key)
	if err != nil {
		t.Fatalf("get encrypt err %v", err)

	}
	t.Logf("src after encrypt: %s", string(x1))

	x2, err := DecryptAES(x1, key)
	if err != nil {
		t.Fatalf("get decrypt err %v", err)
	}
	t.Logf("dst after decrypt: %s", string(x2))

	if d != x2 {
		t.Fatalf("not equal of data before encrypt and after decrypt")
	}
}

func TestEncryptAES_DecryptAES_Base64(t *testing.T) {

	d := "just_test"
	key := "12345678abcdefgh"

	t.Logf("src before encrypt: %s", string(d))

	x1, err := EncryptAESWithBase64(d, key)
	if err != nil {
		t.Fatalf("get encrypt err %v", err)

	}
	t.Logf("src after encrypt: %s", string(x1))

	x2, err := DecryptAESWithBase64(x1, key)
	if err != nil {
		t.Fatalf("get decrypt err %v", err)
	}
	t.Logf("dst after decrypt: %s", string(x2))

	if d != x2 {
		t.Fatalf("not equal of data before encrypt and after decrypt")
	}
}
