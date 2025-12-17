package models

import (
	"bytes"
	"fmt"
	"time"

	"github.com/arod1213/auto_ingestion/excel"
	"github.com/xuri/excelize/v2"
)

type claimBasis string

const (
	CO claimBasis = "Copyright Owner"
	RA claimBasis = "Rights Administrator"
)

type SX struct {
	artist                    string
	title                     string
	isrc                      string
	claimBasis                claimBasis
	claimPercent              float32
	collectionBeganDate       time.Time
	collectionEndDate         time.Time
	nonUsTerritories          []string
	recordingVersion          *string
	duration                  string
	genre                     *string
	recordingDate             *time.Time
	countryOfRecording        *string
	countryOfMastering        *string
	copyrightOwnerNationality *string
	dateOfFirstRelease        *time.Time
	pLine                     *string
	iswc                      string
	composers                 []string
	publishers                []string
	releaseArtist             string
	releaseTitle              string
	releaseVersion            *string
	upc                       uint64
	catalogNum                *string
	releaseDate               time.Time
	countryOfRelease          *string
	releaseLabel              string
}

func sxFromInfo(s Song) SX {
	hours := int(s.Duration.Hours())
	minutes := int(s.Duration.Minutes())
	seconds := int(s.Duration.Seconds()) % 60
	duration := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

	var iswc string
	if s.Iswc != nil {
		iswc = *s.Iswc
	}

	return SX{
		duration:     duration,
		claimBasis:   CO,
		claimPercent: s.Share.MasterPercent,

		collectionBeganDate: s.ReleaseDate,
		dateOfFirstRelease:  &s.ReleaseDate, // this may need to change
		releaseDate:         s.ReleaseDate,
		collectionEndDate:   time.Date(9999, time.Month(12), 30, 0, 0, 0, 0, time.Now().Location()),

		title:         s.Title,
		releaseTitle:  s.Title,
		artist:        s.Artist,
		releaseArtist: s.Artist,
		iswc:          iswc,
		isrc:          s.Isrc,
		upc:           s.Upc,
		releaseLabel:  s.Label,
	}
}

func sxWriteOrder(sx SX) []pos {
	return []pos{
		{key: "Artist \n(*1)", value: sx.artist},
		{key: "Recording Title \n(*1)", value: sx.title},
		{key: "ISRC \n(*1)", value: sx.isrc},
		{key: "What is the basis of your claim?\n(Copyright Owner or Collections Designee)\n(*1)", value: sx.claimBasis},
		{key: "Percentage Claimed \n(*1)", value: sx.claimPercent},
		{key: "Collection Rights Begin Date\n(MM/DD/YYYY) \n(*1)", value: sx.collectionBeganDate},
		{key: "Collection Rights End Date\n(MM/DD/YYYY)\n(*1)", value: sx.collectionEndDate},
		{key: "Non-US Territories of Collection Rights\n (America is entered by default)\n(2)", value: sx.nonUsTerritories},
		{key: "Recording Version, if applicable \n (Ex., \"live\", \"dance remix\", etc)", value: sx.recordingVersion},
		{key: "Duration (HH:MM:SS) \n(2)", value: sx.duration},
		{key: "Genre", value: sx.genre},
		{key: "Recording Date\n(MM/DD/YYYY)", value: sx.recordingDate},
		{key: "Country of Recording/Fixation \n(2)", value: sx.countryOfRecording},
		{key: "Country of Mastering", value: sx.countryOfMastering},
		{key: "Copyright Owner Country of Nationality  \n(2)", value: sx.copyrightOwnerNationality},
		{key: "Date of First Release\n(MM/DD/YYYY)", value: sx.dateOfFirstRelease},
		{key: "Country/Countries of First Release/Publication \n(2)", value: sx.countryOfRelease},
		{key: "(P) Line", value: sx.pLine},
		{key: "ISWC", value: sx.iswc},
		{key: "Composer(s) \n(2)", value: sx.composers},
		{key: "Publisher(s)", value: sx.publishers},
		{key: "Release Artist  \n(*1)", value: sx.releaseArtist},
		{key: "Release Title (Album Title)\n(*1)", value: sx.releaseTitle},
		{key: "Release Version", value: sx.releaseVersion},
		{key: "UPC \n(*1)", value: sx.upc},
		{key: "Catalog # ", value: sx.catalogNum},
		{key: "Release Date\n(MM/DD/YYYY)", value: sx.releaseDate},
		{key: "Country of Release \n(2)", value: sx.countryOfRelease},
		{key: "Release Label \n(*1)", value: sx.releaseLabel},
	}
}

func (s SX) writeSX(file *excelize.File, sheet string, row int) error {
	for i, x := range sxWriteOrder(s) {
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

func SXWrite(songs []Song) (*bytes.Buffer, error) {
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

	err = miscHeader(file, sheet)
	if err != nil {
		fmt.Println("err is ", err)
		return nil, err
	}

	sxHeader(file, sheet)

	row := 11 // offset (1 indexed)
	for _, song := range songs {
		if song.Share.MasterPercent == 0 {
			continue
		}

		s := sxFromInfo(song)
		if err := s.writeSX(file, sheet, row); err != nil {
			fmt.Println("err is ", err)
			continue
		}
		row++
	}

	return file.WriteToBuffer()
}

func sxHeader(file *excelize.File, sheet string) {
	order := sxWriteOrder(SX{})

	for i, k := range order {
		offset := 10
		cell, err := excelize.CoordinatesToCellName(i+1, offset)
		if err != nil {
			fmt.Println("err is ", err)
			continue
		}
		i += 1
		file.SetCellStr(sheet, cell, k.key)
	}
}

func miscHeader(f *excelize.File, sheet string) error {
	date := time.Date(2019, time.Month(10), 30, 0, 0, 0, 0, time.Now().Location())

	excel.WriteTypeAgno(f, sheet, "B1", "ISRC Ingest File")
	excel.WriteTypeAgno(f, sheet, "B2", date)

	excel.WriteTypeAgno(f, sheet, "D1", "Required Fields Key")
	excel.WriteTypeAgno(f, sheet, "D2", "(*1) - All Sound Recording Copyright Owners")
	excel.WriteTypeAgno(f, sheet, "D3", "Required Fields Key")
	excel.WriteTypeAgno(f, sheet, "D4", "(2) - Sound Recording Copyright Owners with International Mandates")

	mergeCells(f, sheet, 2, 9, 1, true)
	excel.WriteTypeAgno(f, sheet, "A9", "Minimum Recording Information")
	mergeCells(f, sheet, 4, 9, 4, true)
	excel.WriteTypeAgno(f, sheet, "D9", "Sound Recording Copyright Owner Claim")
	mergeCells(f, sheet, 12, 9, 9, true)
	excel.WriteTypeAgno(f, sheet, "I9", "Additional Recording Information")
	mergeCells(f, sheet, 7, 9, 22, true)
	excel.WriteTypeAgno(f, sheet, "V9", "Release Information ")
	return nil
}

// doc.writeSpan(.{ .x = 0, .y = 8 }, .{ .horiz = 2 }, "Minimum Recording Information");
// doc.writeSpan(.{ .x = 3, .y = 8 }, .{ .horiz = 4 }, "Sound Recording Copyright Owner Claim");
// doc.writeSpan(.{ .x = 8, .y = 8 }, .{ .horiz = 12 }, "Additional Recording Information");
// doc.writeSpan(.{ .x = 21, .y = 8 }, .{ .horiz = 7 }, "Release Information ");
