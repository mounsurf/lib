package fuzz

func combine(list []string, start, end int, result *[][]string) {
	if start == end {
		temp := []string{}
		for i := 0; i <= end; i++ {
			temp = append(temp, list[i])
		}
		*result = append(*result, temp)
	} else {
		for i := start; i <= end; i++ {
			if start != i {
				temp := list[start]
				list[start] = list[i]
				list[i] = temp
			}
			combine(list, start+1, end, result)
			if start != i {
				temp := list[start]
				list[start] = list[i]
				list[i] = temp
			}
		}
	}
}

/**
 * 获取所有的排列组合方式
 */
func GetCombineList(buf []string) [][]string {
	if buf == nil || len(buf) == 0 {
		return nil
	}
	result := [][]string{}
	combine(buf, 0, len(buf)-1, &result)
	return result
}
