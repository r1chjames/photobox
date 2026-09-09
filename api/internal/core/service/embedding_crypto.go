package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

// EmbeddingCipher encrypts/decrypts face embeddings at rest using AES-256-GCM.
// The key is derived from the app TOKEN secret so face data can never be read
// without the deployment secret, and rotating TOKEN orphans (but never leaks)
// stored embeddings.
type EmbeddingCipher struct {
	aead cipher.AEAD
}

// NewEmbeddingCipher derives an AES-256 key from the app secret (TOKEN env).
func NewEmbeddingCipher(tokenSecret string) (*EmbeddingCipher, error) {
	sum := sha256.Sum256([]byte(tokenSecret))
	block, err := aes.NewCipher(sum[:])
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &EmbeddingCipher{aead: aead}, nil
}

// Encrypt returns nonce || ciphertext.
func (c *EmbeddingCipher) Encrypt(plain []byte) ([]byte, error) {
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return c.aead.Seal(nonce, nonce, plain, nil), nil
}

// Decrypt reverses Encrypt.
func (c *EmbeddingCipher) Decrypt(data []byte) ([]byte, error) {
	if len(data) < c.aead.NonceSize() {
		return nil, errors.New("embedding ciphertext too short")
	}
	nonce, ciphertext := data[:c.aead.NonceSize()], data[c.aead.NonceSize():]
	plain, err := c.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New("embedding decryption failed")
	}
	return plain, nil
}
