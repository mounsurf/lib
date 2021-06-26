package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"github.com/mounsurf/lib/util"
)

//fake aes
// 填充数据
func padding(src []byte, blockSize int) []byte {
	padNum := blockSize - len(src)%blockSize
	pad := bytes.Repeat([]byte{byte(padNum)}, padNum)
	return append(src, pad...)
}

// 去掉填充数据
func unpadding(src []byte) []byte {
	n := len(src)
	unPadNum := int(src[n-1])
	return src[:n-unPadNum]
}

// 加密
func Encrypt(src []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	src = padding(src, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	blockMode.CryptBlocks(src, src)
	return src, nil
}

// 解密
func Decrypt(src []byte, key []byte) ([]byte, error) {
	if len(key) < 16 {
		for {
			key = append(key, 0x00)
			if len(key) == 16 {
				break
			}
		}
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	blockMode.CryptBlocks(src, src)
	src = unpadding(src)
	return src, nil
}

// 解密
func DecryptB64(data string, key []byte) ([]byte, error) {
	src, err := util.B64Decode(data)
	if err != nil {
		return nil, err
	}
	return Decrypt(src, key)
}
