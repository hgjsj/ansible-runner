package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os/exec"
	"strings"
)

// RunCommand executes a command in the specified directory.
func RunCommand(cwd string, env []string, cmd string, args ...string) (string, error) {
	command := exec.Command(cmd, args...)
	Log.Info(fmt.Sprintf("Command directory: %s", cwd))
	command.Dir = cwd

	Log.Info(fmt.Sprintf("Command Env: %s", strings.Join(env, ",")))
	command.Env = append(command.Env, env...)

	Log.Info(fmt.Sprintf("Running command: %s", command.String()))
	output, err := command.Output()
	if err != nil {
		Log.Info(fmt.Sprintf("Command output: %s", string(output)))
		Log.Error(err.Error())
		
		return "", err
	}
	return string(output), nil

}
var encryptKey string = "WB6H!5BTX9-bA6@N-dsKLvWk6AYp#BUr"

// Encrypt Encrypt an string using the provided key.
func Encrypt(text string) (string, error) {
    block, err := aes.NewCipher([]byte(encryptKey))
    if err != nil {
		Log.Error(err.Error())
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
		Log.Error(err.Error())
        return "", err
    }

    block, err := aes.NewCipher([]byte(encryptKey))
    if err != nil {
		Log.Error(err.Error())
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

// ParseBasicAuth parses the Authorization header string for Basic authentication.
func ParseBasicAuth(authHeader string) (authType string, credentials string, ok bool) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}

// DecodeBasicAuth decodes the base64 encoded credentials from Basic authentication.
func DecodeBasicAuth(credentials string) (username, password string, err error) {
	decoded, err := base64.StdEncoding.DecodeString(credentials)
	if err != nil {
		return "", "", fmt.Errorf("invalid base64 encoding: %w", err)
	}

	creds := strings.SplitN(string(decoded), ":", 2)
	if len(creds) != 2 {
		return "", "", fmt.Errorf("invalid credentials format")
	}

	return creds[0], creds[1], nil
}
