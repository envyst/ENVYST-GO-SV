package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"os"
	"strings"
)

// DeriveKey generates a 32-byte key from the password and salt.
func DeriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)
}

// EncryptData encrypts the given data using the password and returns a base64-encoded string.
func EncryptData(data, password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	key := DeriveKey(password, salt)
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(data))
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext, []byte(data))

	// Combine salt, IV, and ciphertext
	combined := append(salt, append(iv, ciphertext...)...)
	return base64.URLEncoding.EncodeToString(combined), nil
}

// DecryptData decrypts the given base64-encoded token using the password and returns the original data.
func DecryptData(token, password string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}

	if len(data) < 32 {
		return "", errors.New("invalid data length")
	}

	salt, iv, ciphertext := data[:16], data[16:32], data[32:]
	key := DeriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintext := make([]byte, len(ciphertext))
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(plaintext, ciphertext)

	return string(plaintext), nil
}

// FileNameEncrypt encrypts the given data for filenames and replaces unsafe characters.
func FileNameEncrypt(data, password string) (string, error) {
	encrypted, err := EncryptData(data, password)
	if err != nil {
		return "", err
	}

	// Replace characters to make it filename safe
	return strings.ReplaceAll(strings.ReplaceAll(encrypted, "/", "_"), "+", "-"), nil
}

// FileNameDecrypt decrypts the given token for filenames and returns the original data.
func FileNameDecrypt(token, password string) (string, error) {
	// Replace characters back to their original form
	token = strings.ReplaceAll(strings.ReplaceAll(token, "_", "/"), "-", "+")
	decrypted, err := DecryptData(token, password)
	if err != nil {
		return "", err
	}

	return decrypted, nil
}
