package util

/**
反转字符串
*/
func Reverse(data string) string {
	var result []byte
	tmp := []byte(data)
	length := len(data)
	for i := 0; i < length; i++ {
		result = append(result, tmp[length-i-1])
	}
	return string(result)

}

func CountCharacter(data string, target byte) int {
	count := 0
	for i := 0; i < len(data); i++ {
		if data[i] == target {
			count++
		}
	}
	return count
}
