package utils

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

type DataExcel struct {
	rows [][]string
	cols []string
}

func ReadExcelFile(filename string) (error, *DataExcel) {

	d := new(DataExcel)

	file, err := excelize.OpenFile(filename)
	if err != nil {
		return err, nil
	}

	rows, err := file.GetRows("Sheet1")
	if err != nil {
		return err, nil
	}

	cols := make([]string, len(rows[0]))

	for i, row := range rows {

		for c, col := range row {
			if i == 0 {
				cols[c] = col
			}
		}
		fmt.Println()
	}

	d.cols = cols
	d.rows = rows

	return nil, d
}
