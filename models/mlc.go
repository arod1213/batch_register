package models

import (
	"bytes"
	"fmt"

	"github.com/arod1213/auto_ingestion/excel"
	"github.com/xuri/excelize/v2"
)

type admin struct {
	mlcNum *string
	name   string
	IpiNum uint64
}

type MLC struct {
	title               string
	songCode            *string
	memberSongId        *string
	iswc                string
	titleAka            *string
	titleAkaCode        *string
	writerLastName      string
	writerFirstName     string
	writerIpiNum        uint64
	writerRoleCode      string
	mlcPublisherNum     *string
	publisherName       string
	publisherIpiNum     uint64
	admin               *admin
	collectionShare     float32
	recordingTitle      string
	recordingArtistName string
	recordingISRC       string
	recordingLabel      string
}

func fromInfo(info Info) MLC {
	// publisher := info.Share.Person.PubEntity.AdminOrPub()

	return MLC{
		iswc:            info.Song.Iswc,
		title:           info.Song.Title,
		collectionShare: info.Share.PubPercent,

		recordingTitle:      info.Song.Title,
		recordingArtistName: info.Song.Artist,
		recordingISRC:       info.Song.Isrc,
		recordingLabel:      info.Song.Label,

		admin: nil,

		writerFirstName: info.Share.Person.FirstName,
		writerLastName:  info.Share.Person.LastName,
		writerIpiNum:    info.Share.Person.WriterIpiNum,

		// publisherName:   publisher.Name,
		// publisherIpiNum: publisher.IpiNum,
		writerRoleCode: "CA", // change this
	}
}

func mlcWriteOrder(mlc MLC) []pos {
	admin := valOrDefault(mlc.admin)

	return []pos{
		{key: "PRIMARY TITLE *", value: mlc.title},
		{key: "MLC SONG CODE", value: mlc.songCode},
		{key: "MEMBERS SONG ID", value: mlc.memberSongId},
		{key: "ISWC          ", value: mlc.iswc},
		{key: "AKA TITLE †", value: mlc.titleAka},
		{key: "AKA TITLE TYPE CODE †", value: mlc.titleAkaCode},
		{key: "WRITER LAST NAME *", value: mlc.writerLastName},
		{key: "WRITER FIRST NAME ", value: mlc.writerFirstName},
		{key: "WRITER IPI NUMBER", value: mlc.writerIpiNum},
		{key: "WRITER ROLE CODE *", value: mlc.writerRoleCode},
		{key: "MLC PUBLISHER NUMBER", value: mlc.mlcPublisherNum},
		{key: "PUBLISHER NAME *", value: mlc.publisherName},
		{key: "PUBLISHER IPI NUMBER *", value: mlc.publisherIpiNum},
		{key: "ADMINISTRATOR MLC PUBLISHER NUMBER", value: admin.mlcNum},
		{key: "ADMINISTRATOR NAME †", value: admin.name},
		{key: "ADMINISTRATOR IPI NUMBER †", value: admin.IpiNum},
		{key: "COLLECTION SHARE *", value: mlc.collectionShare},
		{key: "RECORDING TITLE †", value: mlc.recordingTitle},
		{key: "RECORDING ARTIST NAME †", value: mlc.recordingArtistName},
		{key: "RECORDING ISRC", value: mlc.recordingISRC},
		{key: "RECORDING LABEL", value: mlc.recordingLabel},
	}
}

func (m MLC) writeMLC(file *excelize.File, sheet string, row int) error {
	for i, x := range mlcWriteOrder(m) {
		cell, err := excelize.CoordinatesToCellName(i+1, row)
		if err != nil {
			return err
		}
		err = excel.WriteTypeAgno(file, sheet, cell, x.value)
		if err != nil {
			return err
		}
	}
	return nil
}

func MLCWrite(info []Info) (*bytes.Buffer, error) {
	file := excelize.NewFile()
	defer func() {
		err := file.Close()
		if err != nil {
			println("err is ", err)
		}
	}()

	sheet := "Sheet 1"
	idx, err := file.NewSheet(sheet)
	if err != nil {
		println("FAILED TO SAVE ", err)
		return nil, err
	}

	file.DeleteSheet("Sheet1") // delete default sheet
	file.SetActiveSheet(idx)

	header(file, sheet)

	// write songs
	row := 2 // offset for header (1 indexed)
	for _, x := range info {
		if x.Share.PubPercent == 0 {
			continue
		}

		m := fromInfo(x)
		if err := m.writeMLC(file, sheet, row); err != nil {
			fmt.Println("err is ", err)
			continue
		}
		row++
	}

	return file.WriteToBuffer()
}

func header(file *excelize.File, sheet string) {
	order := mlcWriteOrder(MLC{})

	for i, k := range order {
		cell, err := excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			fmt.Println("err is ", err)
			continue
		}
		i += 1
		file.SetCellStr(sheet, cell, k.key)
	}
}
