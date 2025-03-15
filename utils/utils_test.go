package utils

import (
	"testing"
)

func TestEncrypt(t *testing.T) {
	plainText := "p4sautomation"

	cipherText, err := Encrypt(plainText)
	if err != nil {
		t.Fatalf("Encrypted text %s failed with err %s", plainText, err)
	}
	t.Logf("Encrypted text %s = %s", plainText, cipherText)
}

func TestDecrypt(t *testing.T) {
	plainText := "p4sautomation"

	cipherText, err := Encrypt(plainText)
	if err != nil {
		t.Fatalf("Encryption failed for text %s with error %s", plainText, err)
	}

	decryptedText, err := Decrypt(cipherText)
	if err != nil {
		t.Fatalf("Decryption failed for text %s with error %s", cipherText, err)
	}

	if decryptedText != plainText {
		t.Errorf("Decryption result mismatch: got %s, want %s", decryptedText, plainText)
	}
	t.Logf("Decrypted text %s = %s", cipherText, decryptedText)
}
func TestRunCommand(t *testing.T) {
	cwd := "."
	cmd := "echo"
	args := []string{"hello world"}

	output, err := RunCommand(cwd, cmd, args...)
	if err != nil {
		t.Fatalf("RunCommand failed with error: %s", err)
	}

	if err != nil {
		t.Errorf("RunCommand failed with error: %s", err)
	}
	t.Logf("RunCommand output: %s", output)

	output, err = RunCommand(cwd, "ls", "-l")
	if err != nil {
		t.Fatalf("RunCommand failed with error: %s", err)
	} else {
		t.Logf("RunCommand output: %s", output)
	}
	
}

