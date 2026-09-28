package Song

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"github.com/h2non/filetype"
	"hwyy/KeyStore"
	"os"
)

// IsEncryptedFile 判断是否是加密文件
func IsEncryptedFile(file string) (bool, error) {
	inputFile, err := os.Open(file)
	if err != nil {
		return false, err
	}
	buffer := make([]byte, 512)
	n, err := inputFile.Read(buffer)
	if err != nil {
		return false, err
	}
	if n < 512 {
		return false, nil
	}
	kind, _ := filetype.Match(buffer)
	if kind == filetype.Unknown {
		return true, nil
	}
	return false, nil
}

// Decrypt11 decrypts a file with AES OFB mode
func Decrypt11(file, iv, SecretKey string) error {
	isEncrypted, err := IsEncryptedFile(file)
	if err != nil {
		return err
	}
	if !isEncrypted {
		return nil
	}
	key, err := KeyStore.Decrypt(SecretKey)
	if err != nil {
		return err
	}
	ivByte, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	inputFile, err := os.Open(file)
	if err != nil {
		return err
	}

	defer func(inputFile *os.File) {
		err = inputFile.Close()
		if err != nil {
			return
		}
	}(inputFile)

	outputFile, err := os.Create(file + ".dec")
	if err != nil {
		return err
	}
	defer func(outputFile *os.File) {
		err := outputFile.Close()
		if err != nil {
			return
		}
	}(outputFile)

	chunkSize := 2048
	buffer := make([]byte, chunkSize)
	for {
		n, err := inputFile.Read(buffer)
		if err != nil && err.Error() != "EOF" {
			return err
		}
		if n == 0 {
			break
		}
		mode := cipher.NewOFB(block, ivByte)
		mode.XORKeyStream(buffer, buffer)
		_, err = outputFile.Write(buffer[:n])
		if err != nil {
			return err
		}
	}
	err = os.Remove(file)
	if err != nil {
		return err
	}
	err = os.Rename(file+".dec", file)
	return err
}
