package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"log"
	"os"
)

func Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(getAESKey())
	if err != nil {
		log.Fatalln("failed to create block", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalln("failed to create gcm", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalln("failed to read nonce", err)
	}

	cipherText := gcm.Seal(nonce, nonce, data, nil)

	return cipherText, nil
}

func Decrypt(cipherText []byte) ([]byte, error) {
	block, err := aes.NewCipher(getAESKey())
	if err != nil {
		log.Fatalln("failed to create block", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalln("failed to create gcm", err)
	}

	nonce := cipherText[:gcm.NonceSize()]
	cipherText = cipherText[gcm.NonceSize():]

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		log.Fatalln("failed to decrypt", err)
	}

	return plainText, nil
}

func getAESKey() []byte {
	rawKey := os.Getenv("AES_KEY")
	hash := sha256.Sum256([]byte(rawKey))
	return hash[:]
}
