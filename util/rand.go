package util

import (
	"math/rand"
	"time"
)

func RandByte(n int, letters string) []byte {
	if letters == "" {
		letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	b := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return b
}

func RandStr(n int, letters string) string {
	return string(RandByte(n, letters))
}

//生成随机字符串，第二个参数可选，为字符集
func RandString(n int, letters string) string {
	if letters == "" {
		letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	b := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
