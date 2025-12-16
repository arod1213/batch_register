package models

import "github.com/xuri/excelize/v2"

type pos struct {
	key   string
	value any
}

func valOrDefault[T any](x *T) T {
	var y T
	if x != nil {
		y = *x
	}
	return y
}

func mergeCells(f *excelize.File, sheet string, span int, row int, col int, horiz bool) error {
	x, err := excelize.CoordinatesToCellName(col, row)
	if err != nil {
		return err
	}

	var y string
	if horiz {
		y, err = excelize.CoordinatesToCellName(col+span-1, row)
		if err != nil {
			return err
		}
	} else {
		y, err = excelize.CoordinatesToCellName(col, row+span-1)
		if err != nil {
			return err
		}
	}

	return f.MergeCell(sheet, x, y)
}
