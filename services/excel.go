package services

import (
	"archive/zip"
	"bytes"
	"github.com/arod1213/auto_ingestion/models"
)

func WriteShares(shares []models.Share, user models.User) (*[]byte, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	count := 0

	mlcFile, err := models.MLCWrite(shares, user)
	if err == nil {
		f2, err := zipWriter.Create("mlc.xlsx")
		if err != nil {
			return nil, err
		}

		_, err = f2.Write(mlcFile.Bytes())
		if err != nil {
			return nil, err
		}
		count++
	}

	sxFile, err := models.SXWrite(shares)
	if err == nil {
		f1, err := zipWriter.Create("sx.xlsx")
		if err != nil {
			return nil, err
		}

		_, err = f1.Write(sxFile.Bytes())
		if err != nil {
			return nil, err
		}
		count++
	}

	err = zipWriter.Close()
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, err

	}

	data := buf.Bytes()
	return &data, nil
}

