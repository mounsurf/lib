package util

import (
	"encoding/hex"
	"fmt"
	"github.com/alwindoss/morse"
	"strconv"
	"strings"
)

/**
解码hex数据为string
如: 776562 -> web

*/
func HexToString(data string) (string, error) {
	src := []byte(data)
	dst := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(dst, src)
	if err != nil {
		return "", err
	}
	return string(dst), nil
}

/**
解码hex数据为int
如: 10 -> 16

*/
func HexToUInt32(data string) (uint32, error) {
	n, err := strconv.ParseUint(data, 16, 32)
	if err != nil {
		return 0, err
	}
	return uint32(n), nil
}

/**
int to hex
*/
func Int64ToHex(n int64) string {
	if n < 0 {
		return ""
	}
	if n == 0 {
		return "0"
	}
	hex := map[int64]int64{10: 65, 11: 66, 12: 67, 13: 68, 14: 69, 15: 70}
	s := ""
	for q := n; q > 0; q = q / 16 {
		m := q % 16
		if m > 9 && m < 16 {
			m = hex[m]
			s = fmt.Sprintf("%v%v", string(m), s)
			continue
		}
		s = fmt.Sprintf("%v%v", m, s)
	}
	return s
}

func getRot13Step(x byte) int {
	if x >= 'a' && x < 'n' || x >= 'A' && x < 'N' {
		return 13
	}
	if x >= 'n' && x <= 'z' || x >= 'N' && x <= 'Z' {
		return -13
	}
	return 0
}

/**
 * rot 13解码
 */
func Rot13Decode(strData string) string {
	byteData := []byte(strData)
	result := ""
	for i := 0; i < len(byteData); i += 1 {
		result += string(int(byteData[i]) + getRot13Step(byteData[i]))
	}
	return result
}

/**
 * morse编码
 */
func MorseEncode(data string) (string, error) {
	h := morse.NewHacker()
	text, err := h.Encode(strings.NewReader(data))
	if err != nil {
		return "", err
	}
	return string(text), nil
}

/**
 * morse解码
 */
func MorseDecode(data string) (string, error) {
	h := morse.NewHacker()
	text, err := h.Decode(strings.NewReader(data))
	if err != nil {
		return "", err
	}
	return string(text), nil
}

func BinaryStringToAscii(binary string) byte {
	if len(binary) != 8 {
		return 0
	}
	result := 0
	uint8arr := [8]int{128, 64, 32, 16, 8, 4, 2, 1}
	for i := 0; i < 8; i++ {
		if binary[i] == '1' {
			result += uint8arr[i]
		}
	}
	return byte(result)
}
