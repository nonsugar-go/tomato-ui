package excel

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/mattn/go-runewidth"
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
	if err := f.SetDefaultFont("游ゴシック"); err != nil {
		_ = f.Close()
		slog.Error("Failed to set default font", "error", err)
		os.Exit(1)
	}
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
	e.AutoFitColumns()
}

func (e *Excel) AutoFitColumns() error {
	rows, err := e.f.GetRows(e.sheet)
	if err != nil {
		return err
	}

	if len(rows) == 0 {
		return nil
	}

	maxRow := len(rows)
	if maxRow > 100 {
		maxRow = 100
	}

	maxCol := 0
	for i := 0; i < maxRow; i++ {
		if len(rows[i]) > maxCol {
			maxCol = len(rows[i])
		}
	}

	for colIdx := 0; colIdx < maxCol; colIdx++ {
		maxWidth := 0

		for rowIdx := 0; rowIdx < maxRow; rowIdx++ {
			if colIdx >= len(rows[rowIdx]) {
				continue
			}

			cell := rows[rowIdx][colIdx]
			w := runewidth.StringWidth(cell)

			if w > maxWidth {
				maxWidth = w
			}
		}

		width := float64(maxWidth) + 2

		if width < 4 {
			width = 4
		}
		if width > 30 {
			width = 30
		}

		colName, err := excelize.ColumnNumberToName(colIdx + 1)
		if err != nil {
			return err
		}

		if err := e.f.SetColWidth(e.sheet, colName, colName, width); err != nil {
			return err
		}
	}

	return nil
}
