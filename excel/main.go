package excel

import (
	"reflect"
	"time"

	"github.com/xuri/excelize/v2"
)

func unwrap(v any) any {
	for {
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Pointer {
			return v
		}
		if rv.IsNil() {
			return nil
		}
		v = rv.Elem().Interface()
	}
}

func WriteTypeAgno(file *excelize.File, sheet string, cell string, val any) error {
	x := unwrap(val)
	if x == nil {
		return nil
	}

	i := info{file: file, sheet: sheet}
	switch v := x.(type) {
	case string:
		return file.SetCellStr(sheet, cell, v)
	case time.Time:
		// TODO: allow custom date formats
		return writeDate(i, cell, v, "mm/dd/yyyy")
	case uint64:
		styleStr := "000000000000"

		styleID, err := file.NewStyle(&excelize.Style{
			CustomNumFmt: &styleStr,
		})
		if err != nil {
			return err
		}

		err = file.SetCellStyle(sheet, cell, cell, styleID)
		if err != nil {
			return err
		}

		return file.SetCellUint(sheet, cell, v)
	case int:
		return file.SetCellInt(sheet, cell, int64(v))
	case float32:
		return file.SetCellFloat(sheet, cell, float64(v), 2, 32)
	case float64:
		return file.SetCellFloat(sheet, cell, v, 2, 32)
	default:
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.String {
			return file.SetCellStr(sheet, cell, rv.String())
		}
	}
	return nil
}

type info struct {
	file  *excelize.File
	sheet string
}

func writeDate(i info, cell string, d time.Time, format string) error {
	style, err := i.file.NewStyle(&excelize.Style{
		CustomNumFmt: &format,
	})

	if err != nil {
		return err
	}

	i.file.SetCellStyle(i.sheet, cell, cell, style)
	return i.file.SetCellValue(i.sheet, cell, d)
}
