package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
)

// RunCommand executes a command in the specified directory.
func RunCommand(cwd string, cmd string, args ...string) (string, error) {
	command := exec.Command(cmd, args...)
	command.Dir = cwd
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	fmt.Println("Running command: ", command.String())
	if err := command.Run(); err != nil {
		return "", err
	}
	output, err := command.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil

}
var encryptKey string = "WB6H!5BTX9-bA6@N-dsKLvWk6AYp#BUr"
// Encrypt Encrypt an string using the provided key.
func Encrypt(text string) (string, error) {
    block, err := aes.NewCipher([]byte(encryptKey))
    if err != nil {
        return "", err
    }

    plaintext := []byte(text)
    cfb := cipher.NewCFBEncrypter(block, []byte(encryptKey)[:aes.BlockSize])
    ciphertext := make([]byte, len(plaintext))
    cfb.XORKeyStream(ciphertext, plaintext)
    return hex.EncodeToString(ciphertext), nil
}

// Decrypt decrypts an encrypted string using the provided key.
func Decrypt(encryptedText string) (string, error) {
    ciphertext, err := hex.DecodeString(encryptedText)
    if err != nil {
        return "", err
    }

    block, err := aes.NewCipher([]byte(encryptKey))
    if err != nil {
        return "", err
    }

    // if len(ciphertext) < aes.BlockSize {
    //     return "", errors.New("ciphertext too short")
    // }

    cfb := cipher.NewCFBDecrypter(block, []byte(encryptKey)[:aes.BlockSize])
    plaintext := make([]byte, len(ciphertext))
    cfb.XORKeyStream(plaintext, ciphertext)
    return string(plaintext), nil
}
