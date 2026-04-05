package excel

import (
	"fmt"
	"log"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Excel struct {
	f        *excelize.File
	filename string
	sheet    string
	c, r     int // Col, Row
}

func NewExcel(filename string) *Excel {
	f := excelize.NewFile()
	return &Excel{f: f, filename: filename, sheet: "Sheet1", c: 1, r: 1}
}

func (e *Excel) Close() {
	if err := e.f.DeleteSheet("Sheet1"); err != nil {
		log.Print(err)
	}
	if err := e.f.SaveAs(e.filename); err != nil {
		log.Print(err)
	}
	if err := e.f.Close(); err != nil {
		log.Print(err)
	}
}

func (e *Excel) NewSheet(sheet string) {
	index, err := e.f.NewSheet(sheet)
	if err != nil {
		log.Fatal(err)
	}
	e.sheet = sheet
	e.c, e.r = 1, 1
	e.f.SetActiveSheet(index)
}

func (e *Excel) Println(ss ...string) {
	cell, err := excelize.CoordinatesToCellName(e.c, e.r)
	if err != nil {
		log.Fatal(err)
	}
	e.f.SetSheetRow(e.sheet, cell, &ss)
	e.r++
}

func (e *Excel) Printf(format string, args ...any) {
	str := fmt.Sprintf(format, args...)
	str = strings.TrimRight(str, "\n")
	e.Println(str)
}

func (e *Excel) AddTable() {
	rows, err := e.f.GetRows(e.sheet)
	if err != nil {
		log.Fatal(err)
	}
	lastRow := len(rows)
	if lastRow == 0 {
		log.Print("データがありません")
		return
	}
	lastColNum := 0
	for _, row := range rows {
		if len(row) > lastColNum {
			lastColNum = len(row)
		}
	}
	lastColName, err := excelize.ColumnNumberToName(lastColNum)
	if err != nil {
		log.Fatal(err)
	}
	dimension := fmt.Sprintf("A1:%s%d", lastColName, lastRow)
	err = e.f.AddTable(e.sheet, &excelize.Table{
		Range:     dimension,
		StyleName: "TableStyleMedium9",
	})
	if err != nil {
		log.Print(err)
	}
}
