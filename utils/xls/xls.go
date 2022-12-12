package xls

import (
	"bytes"
	"github.com/shakinm/xlsReader/xls"
	"log"
)

// ReadXls 读取 XLS 数据
func ReadXls(data []byte, sheetIndex int, skipRow int) [][]string {
	wb, err := xls.OpenReader(bytes.NewReader(data))
	if err != nil {
		log.Fatalf("读取 Xls 出错 : %s", err)
	}

	sheet, err := wb.GetSheet(sheetIndex)
	if err != nil {
		log.Fatalf("读取 Sheet 出错 : %s", err)
	}

	sheetData := make([][]string, 0)
	for i := 0; i < sheet.GetNumberRows(); i++ {
		if i == skipRow {
			continue
		}

		rowData := make([]string, 0)
		row, err := sheet.GetRow(i)
		if err != nil {
			log.Fatalf("读取 Row 出错 : %s", err)
		}

		for _, cellData := range row.GetCols() {
			rowData = append(rowData, cellData.GetString())
		}
		sheetData = append(sheetData, rowData)
	}

	return sheetData
}
