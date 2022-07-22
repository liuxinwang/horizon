package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func EncryptAES(key []byte, plaintext string) string {
	origData := []byte(plaintext)
	c, err := aes.NewCipher(key)
	CheckError(err)
	blockSize := c.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(c, key[:blockSize])
	out := make([]byte, len(origData))
	blockMode.CryptBlocks(out, origData)
	return base64.StdEncoding.EncodeToString(out)
}
func DecryptAES(key []byte, ct string) string {
	cipherByte, _ := base64.StdEncoding.DecodeString(ct)

	c, err := aes.NewCipher(key)
	CheckError(err)
	blockSize := c.BlockSize()
	blockMode := cipher.NewCBCDecrypter(c, key[:blockSize])

	pt := make([]byte, len(cipherByte))
	blockMode.CryptBlocks(pt, cipherByte)
	pt = PKCS7UnPadding(pt)

	s := string(pt[:])
	return s
}

func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
