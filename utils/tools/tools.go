package tools

// IndexOf 查询指定元素在数组中的位置
func IndexOf[T string](source []T, target T) int {
	for index, item := range source {
		if item == target {
			return index
		}
	}
	return -1
}

// FetchColData 获取指定列的数据
func FetchColData(table [][]string, colIndex int) []string {
	data := make([]string, 0)
	for _, row := range table {
		for i := 0; i < len(row); i++ {
			if i == colIndex {
				data = append(data, row[i])
			}
		}
	}
	return data
}
