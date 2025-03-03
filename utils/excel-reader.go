package utils

import (
	"github.com/xuri/excelize/v2"
)

type DataExcel struct {
	Rows [][]string
	Cols []string
}

func ReadExcelFile(filename string) (error, *DataExcel) {

	d := new(DataExcel)

	file, err := excelize.OpenFile(filename)
	if err != nil {
		return err, nil
	}

	rows, err := file.GetRows("Form Responses 1")
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
	}

	d.Cols = cols
	d.Rows = rows

	return nil, d
}
