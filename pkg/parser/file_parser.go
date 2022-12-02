package parser

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/mitchellh/mapstructure"
	"github.com/xuri/excelize/v2"
)

type FieldSchema struct {
	Col  string `json:"col"`
	Name string `json:"name"`
}

type Schema struct {
	Version    string        `json:"version"`
	SchemaType string        `json:"schema_type"`
	Headers    bool          `json:"headers"`
	Fields     []FieldSchema `json:"fields"`
}

// ParseCSVFile maps csv file data b to in struct pointer according to s schema.
// Returns an error.
func ParseCSVFile[T any](in *[]T, r io.Reader, s Schema) error {
	fb, err := csvToXlsx(r)
	if err != nil {
		return err
	}

	err = ParseXLSXFile(in, fb, s)
	if err != nil {
		return err
	}

	return nil
}

// csvToXlsx converts csv file data to xslx file byte slice
// Returns byte slice pointer and error.
func csvToXlsx(r io.Reader) (io.Reader, error) {
	csvReader := csv.NewReader(r)
	f := excelize.NewFile()
	f.NewSheet("default")
	worksheet, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	for i, row := range worksheet {
		start := fmt.Sprintf("A%v", i+1)
		err = f.SetSheetRow("default", start, &row)
		if err != nil {
			return nil, err
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}

// ParseXLSXFile Reads file Reader r, maps file data to T type according to s schema and append it to slice in.
// Returns an error.
func ParseXLSXFile[T any](in *[]T, r io.Reader, s Schema) error {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}

	dataMapList, err := getDataMapList(f, s)
	if err != nil {
		return err
	}

	for _, fileItem := range dataMapList {
		var item T
		err := mapstructure.WeakDecode(fileItem, &item)
		*in = append(*in, item)
		if err != nil {
			return err
		}
	}

	return nil
}

// getDataMapList maps excelize.File f to value map according to schema s.
// Returns data map and error.
func getDataMapList(f *excelize.File, s Schema) ([]map[string]string, error) {
	var dataMapList []map[string]string
	sl := f.GetSheetList()
	for _, sheetName := range sl {
		currRowIndex := 1
		rows, err := f.Rows(sheetName)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			if s.Headers && currRowIndex == 1 {
				currRowIndex++
				continue
			}
			fim, err := mapRow(currRowIndex, sheetName, f, s)
			if err != nil {
				return nil, err
			}
			dataMapList = append(dataMapList, fim)
			currRowIndex++
		}

		if err = rows.Close(); err != nil {
			return nil, err
		}
	}
	return dataMapList, nil
}

// mapRow maps excelize.File f row data to value map using schema by index and sheetName.
// Returns row data values map and an error.
func mapRow(index int, sheetName string, f *excelize.File, s Schema) (map[string]string, error) {
	fim := make(map[string]string)
	for _, fs := range s.Fields {
		value, err := f.GetCellValue(sheetName, fmt.Sprintf("%s%v", fs.Col, index))
		if err != nil {
			return nil, err
		}
		fim[fs.Name] = value
	}
	return fim, nil
}
