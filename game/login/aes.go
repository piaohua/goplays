package login

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
)

func Decrypt(ciphertext, key []byte) ([]byte, error) {
	pkey := PaddingLeft(key, '0', 16)
	block, err := aes.NewCipher(pkey)
	if err != nil {
		return nil, err
	}
	blockModel := cipher.NewCBCDecrypter(block, pkey)
	plantText := make([]byte, len(ciphertext))
	blockModel.CryptBlocks(plantText, []byte(ciphertext))
	plantText, err = PKCS7UnPadding(plantText, block.BlockSize())
	if err != nil {
		return nil, err
	}
	return plantText, nil
}

func PKCS7UnPadding(plantText []byte, blockSize int) ([]byte, error) {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	if unpadding > length {
		return nil, errors.New("faild")
	}
	return plantText[:(length - unpadding)], nil
}

func PaddingLeft(ori []byte, pad byte, length int) []byte {
	if len(ori) >= length {
		return ori[:length]
	}
	pads := bytes.Repeat([]byte{pad}, length-len(ori))
	return append(pads, ori...)
}

func TouristAccount(account, key string) (string, error) {
	ciphertext, err1 := hex.DecodeString(account)
	if err1 != nil {
		return "", err1
	}
	b, err2 := Decrypt([]byte(ciphertext), []byte(key))
	if err2 != nil {
		return "", err2
	}
	return string(b), nil
}
