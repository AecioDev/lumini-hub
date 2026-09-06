package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
)

// EncryptAES criptografa plaintext com AES-256-GCM usando hexKey (chave em
// hexadecimal, 64 caracteres = 32 bytes — ver config.SecurityConfig.CertificateEncryptionKey).
// Usado pra dados sensíveis que precisam ser recuperados depois (ex.: senha
// de certificado digital A1) — nunca pra senha de usuário, que continua
// hash irreversível (bcrypt, ver password.go). O nonce aleatório é
// prependado ao ciphertext e tudo é retornado em base64, pronto pra gravar
// direto numa coluna de texto.
func EncryptAES(plaintext string, hexKey string) (string, error) {
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return "", errors.New("chave de criptografia inválida: não é hexadecimal")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES reverte EncryptAES.
func DecryptAES(ciphertextB64 string, hexKey string) (string, error) {
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return "", errors.New("chave de criptografia inválida: não é hexadecimal")
	}

	data, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("dado criptografado inválido")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
