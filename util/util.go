package util

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"net"
)

func B64Encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func B64Decode(src string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(src)
}

func Md5(src []byte) string {
	h := md5.New()
	h.Write(src)
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func ContainsStr(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func ContainsInt(list []int, item int) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func ContainsStrArr(list []string, item []string) bool {
	for _, v := range item {
		if !ContainsStr(list, v) {
			return false
		}
	}
	return true
}

func MapKeyToList(tmap map[string]struct{}) []string {
	keyList := []string{}
	if tmap == nil {
		return keyList
	}
	for key, _ := range tmap {
		keyList = append(keyList, key)
	}
	return keyList
}

func IsPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		case ip4[0] == 100 && ip4[1] >= 64 && ip4[1] <= 127:
			return false
		case ip4[0] == 169 && ip4[1] == 254:
			return false
		default:
			return true
		}
	}
	return false
}
func Substr(s string, start, length int) string {
	r := []rune(s)
	l := len(r)
	start = start % l
	if start < 0 {
		start = l + start
	}
	length = length % l
	if length < 0 {
		length = l + length
	}
	end := start + length
	if end <= l {
		return string(r[start:end])
	}
	return ""
}
