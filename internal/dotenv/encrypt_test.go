package dotenv

import (
	"testing"
)

const testKey = "0123456789abcdef" // 16 bytes for AES-128

func TestEncryptDecrypt_Roundtrip(t *testing.T) {
	plaintext := "super-secret-value"

	enc, err := Encrypt(plaintext, testKey)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if enc == plaintext {
		t.Error("encrypted value should not equal plaintext")
	}

	dec, err := Decrypt(enc, testKey)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if dec != plaintext {
		t.Errorf("expected %q, got %q", plaintext, dec)
	}
}

func TestEncrypt_DifferentCiphertexts(t *testing.T) {
	plaintext := "value"

	enc1, _ := Encrypt(plaintext, testKey)
	enc2, _ := Encrypt(plaintext, testKey)

	if enc1 == enc2 {
		t.Error("two encryptions of the same value should differ due to random nonce")
	}
}

func TestDecrypt_InvalidKey(t *testing.T) {
	enc, _ := Encrypt("value", testKey)

	_, err := Decrypt(enc, "wrongkey12345678")
	if err == nil {
		t.Error("expected error when decrypting with wrong key")
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	_, err := Decrypt("not-valid-base64!!!", testKey)
	if err == nil {
		t.Error("expected error for invalid base64 input")
	}
}

func TestEncryptMap_RoundTrip(t *testing.T) {
	secrets := map[string]string{
		"DB_PASS": "hunter2",
		"API_KEY": "abc123",
	}

	encrypted, err := EncryptMap(secrets, testKey)
	if err != nil {
		t.Fatalf("EncryptMap failed: %v", err)
	}

	for k, v := range secrets {
		dec, err := Decrypt(encrypted[k], testKey)
		if err != nil {
			t.Errorf("Decrypt failed for key %q: %v", k, err)
		}
		if dec != v {
			t.Errorf("key %q: expected %q, got %q", k, v, dec)
		}
	}
}
