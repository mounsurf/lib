package util


func GetMapKeys(m map[string]interface{}) []string {
	result := []string{}
	for k, _ := range m {
		result = append(result, k)
	}
	return result
}