package parser

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"reflect"
)

const SchemaTagName = "schema"

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

// ParseCSVFile maps csv file data b to l struct pointer according to s schema.
// Returns an error.
func ParseCSVFile(l *[]interface{}, b *[]byte, s Schema) error {
	fb, err := csvToXlsx(b)
	if err != nil {
		return err
	}

	err = ParseXLSXFile(l, fb, s)
	if err != nil {
		return err
	}

	return nil
}

// csvToXlsx converts csv file data to xslx file byte slice
// Returns byte slice pointer and error.
func csvToXlsx(data *[]byte) (*[]byte, error) {
	reader := bytes.NewReader(*data)
	csvReader := csv.NewReader(reader)
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

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	b := buffer.Bytes()
	return &b, nil
}

// ParseXLSXFile ParseXlsxFile maps xlsx file data b to l struct pointer according to s schema.
// Returns an error.
func ParseXLSXFile(l *[]interface{}, b *[]byte, s Schema) error {
	s = defaultSchema
	f, err := readFileFromBuffer(b)
	if err != nil {
		return err
	}

	dataMapList, err := getDataMapList(f, s)
	if err != nil {
		return err
	}

	for _, fileItem := range dataMapList {
		s := struct{}{}
		err := mapStruct(&s, fileItem)
		*l = append(*l, s)
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

// readFileFromBuffer reads bytes slice b to excelize.File
// Returns excelize.File and an error.
func readFileFromBuffer(b *[]byte) (*excelize.File, error) {
	reader := bytes.NewReader(*b)
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return f, err
	}
	return f, nil
}

// mapStruct maps valueMap to s struct pointer according to struct field tags.
// Returns an error.
func mapStruct(s any, valueMap map[string]string) error {
	t := reflect.ValueOf(s).Elem()
	typeOfT := t.Type()
	for i := 0; i < t.NumField(); i++ {
		typeField := typeOfT.Field(i)
		tag := typeField.Tag.Get(SchemaTagName)
		err := setField(t.FieldByName(typeField.Name), valueMap[tag])
		if err != nil {
			return err
		}
	}
	return nil
}

// setField sets field value.
// Returns an error.
func setField(field reflect.Value, value any) error {
	if !field.IsValid() {
		return fmt.Errorf("No such field: %s in obj", field)
	}

	if !field.CanSet() {
		return fmt.Errorf("Cannot set %s field value", field)
	}

	structFieldType := field.Type()
	val := reflect.ValueOf(value)

	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	field.Set(val)
	return nil
}
