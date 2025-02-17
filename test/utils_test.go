package test

import (
	"testing"
	"ansible-runner/utils"
)

func TestEncrypt(t *testing.T) {
	plainText := "Hello, World!"

	cipherText, err := utils.Encrypt(plainText)
	if err != nil {
		t.Fatalf("Encrypted text %s failed with err %s", plainText, err)
	}
	t.Logf("Encrypted text %s = %s", plainText, cipherText)
}

func TestDecrypt(t *testing.T) {
	plainText := "Hello, World!"

	cipherText, err := utils.Encrypt(plainText)
	if err != nil {
		t.Fatalf("Encryption failed for text %s with error %s", plainText, err)
	}

	decryptedText, err := utils.Decrypt(cipherText)
	if err != nil {
		t.Fatalf("Decryption failed for text %s with error %s", cipherText, err)
	}

	if decryptedText != plainText {
		t.Errorf("Decryption result mismatch: got %s, want %s", decryptedText, plainText)
	}
	t.Logf("Decrypted text %s = %s", cipherText, decryptedText)
}
