package KeyStore

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
)

// 每次启动生成一个对密钥
var (
	PublicKey  rsa.PublicKey
	PrivateKey rsa.PrivateKey
)

func init() {
	var err error
	PrivateKey, PublicKey, err = GenerateKey(2048)
	if err != nil {
		log.Fatalf("Error generating keys: %v", err)
	}
}

// GenerateKey generates a new RSA key pair.
func GenerateKey(bits int) (rsa.PrivateKey, rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return rsa.PrivateKey{}, rsa.PublicKey{}, err
	}
	return *privateKey, privateKey.PublicKey, nil
}

// Decrypt decrypts the given base64-encoded string using the private key.
func Decrypt(data string) ([]byte, error) {
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 data: %v", err)
	}
	decrypted, err := PrivateKey.Decrypt(rand.Reader, decodedData, &rsa.OAEPOptions{
		Hash:    crypto.SHA256,
		MGFHash: crypto.SHA1, // Should match the hash function used for encryption
	})
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %v", err)
	}
	return decrypted, nil
}

// ExportPublicKey returns the public key in PEM format.
func ExportPublicKey() string {
	pubKeyDER, err := x509.MarshalPKIXPublicKey(&PublicKey)
	if err != nil {
		log.Fatalf("Error marshalling public key: %v", err)
	}
	pemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyDER,
	}
	pub := base64.StdEncoding.EncodeToString(pemBlock.Bytes)
	return pub
}
