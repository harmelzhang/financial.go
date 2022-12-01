package xls

import (
	"bytes"
	"github.com/extrame/xls"
	"log"
)

// ReadXls 读取 XLS 数据
func ReadXls(data []byte, sheetIndex int, skipRow int) [][]string {
	wb, err := xls.OpenReader(bytes.NewReader(data), "UTF-8")
	if err != nil {
		log.Fatalf("执行出错 : %s", err)
	}

	sheet := wb.GetSheet(sheetIndex)
	sheetData := make([][]string, 0)
	for i := 0; i < int(sheet.MaxRow); i++ {
		if i == skipRow {
			continue
		}
		row := sheet.Row(i)
		rowData := make([]string, 0)
		for j := 0; j < row.LastCol(); j++ {
			colData := row.Col(j)
			rowData = append(rowData, colData)
		}
		sheetData = append(sheetData, rowData)
	}

	return sheetData
}
