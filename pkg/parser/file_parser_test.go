package parser

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/xuri/excelize/v2"
)

func getFileReader(t *testing.T, sheets []string, data map[string][]interface{}) io.Reader {
	f := excelize.NewFile()

	for _, sheet := range sheets {
		f.NewSheet(sheet)
		for addr, row := range data {
			err := f.SetSheetRow(sheet, addr, &row)
			if err != nil {
				return nil
			}
		}
	}

	buffer, err := f.WriteToBuffer()
	if err != nil {
		t.Fatal(err)
	}

	return bytes.NewReader(buffer.Bytes())
}

func TestParseFile(t *testing.T) {
	type student struct {
		Name    string `mapstructure:"first_name"`
		Surname string `mapstructure:"last_name"`
		Email   string `mapstructure:"email"`
	}
	//var s []student
	type args struct {
		in *[]student
		b  io.Reader
		s  Schema
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantLen int
	}{
		{
			name: "parse xlsx file with headers",
			args: args{
				&[]student{},
				getFileReader(
					t,
					[]string{"Sheet1"},
					map[string][]interface{}{
						"A1": {"First Name", "Last Name", "Email"},
						"A2": {"Anakin", "Skywalker", "anakin.skywalker@deathstar.imp"},
						"A3": {"Obi-Wan", "Kenobi", "obi@jedi.rules"},
					},
				),
				defaultSchema,
			},
			wantErr: false,
			wantLen: 2,
		}, {
			name: "parse xlsx file with multiple sheets",
			args: args{
				&[]student{},
				getFileReader(
					t,
					[]string{"Sheet1", "Sheet2"},
					map[string][]interface{}{
						"A1": {"First Name", "Last Name", "Email"},
						"A2": {"Anakin", "Skywalker", "anakin.skywalker@deathstar.imp"},
						"A3": {"Obi-Wan", "Kenobi", "obi@jedi.rules"},
					},
				),
				defaultSchema,
			},
			wantErr: false,
			wantLen: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ParseXLSXFile(tt.args.in, tt.args.b, tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("ParseXLSXFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(*tt.args.in) != tt.wantLen {
				t.Errorf("ParseXLSXFile() len = %v, wantLen %v", len(*tt.args.in), tt.wantLen)
			}

			if n := *tt.args.in; n[0].Name != "Anakin" {
				t.Errorf("Expected name")
			}
		})
	}
}

func TestParseCSVFile(t *testing.T) {
	type student struct {
		Name    string `mapstructure:"first_name"`
		Surname string `mapstructure:"last_name"`
		Email   string `mapstructure:"email"`
	}

	type args struct {
		in *[]student
		b  io.Reader
		s  Schema
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []student
	}{
		{
			name: "parse csv file with headers",
			args: args{
				&[]student{},
				bytes.NewReader([]byte{
					70, 105, 114, 115, 116, 32, 78, 97, 109, 101, 44, 76, 97, 115, 116, 32, 78, 97, 109, 101, 44, 69, 109, 97, 105, 108, 10, 65,
					110, 97, 107, 105, 110, 44, 83, 107, 121, 119, 97, 108, 107, 101, 114, 44, 97, 110, 97, 107, 105, 110, 46, 115, 107, 121, 119,
					97, 108, 107, 101, 114, 64, 100, 101, 97, 116, 104, 115, 116, 97, 114, 46, 105, 109, 112, 10, 79, 98, 105, 45, 87, 97, 110, 44,
					75, 101, 110, 111, 98, 105, 44, 111, 98, 105, 64, 106, 101, 100, 105, 46, 114, 117, 108, 101, 115,
				}),
				defaultSchema,
			},
			wantErr: false,
			want: []student{
				{"Anakin", "Skywalker", "anakin.skywalker@deathstar.imp"},
				{"Obi-Wan", "Kenobi", "obi@jedi.rules"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ParseCSVFile(tt.args.in, tt.args.b, tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("ParseCSVFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			for i, s := range tt.want {
				got := *tt.args.in
				if !reflect.DeepEqual(got[i], s) {
					t.Errorf("ParseCSVFile() want = %v, got = %v", s, got[i])
				}
			}
		})
	}
}

func TestParseXLSXFile(t *testing.T) {
	type student struct {
		Name        string   `mapstructure:"first_name"`
		Surname     string   `mapstructure:"last_name"`
		Email       string   `mapstructure:"email"`
		Age         int      `mapstructure:"age"`
		IsGraduated bool     `mapstructure:"graduated"`
		Languages   []string `mapstructure:"languages"`
	}

	schema := Schema{
		"1",
		"",
		true,
		[]FieldSchema{
			{Name: "first_name", Col: "A"},
			{Name: "last_name", Col: "B"},
			{Name: "email", Col: "C"},
			{Name: "age", Col: "D"},
			{Name: "graduated", Col: "E"},
			{Name: "languages", Col: "F", IsMultiple: true},
			{Name: "languages", Col: "G", IsMultiple: true},
		},
	}
	type args struct {
		in *[]student
		r  io.Reader
		s  Schema
	}
	tests := []struct {
		name    string
		args    args
		want    []student
		wantLen int
		wantErr bool
	}{
		{
			"parse multiple data type",
			args{
				&[]student{},
				getFileReader(
					t,
					[]string{"Sheet1", "Sheet2"},
					map[string][]interface{}{
						"A1": {"First Name", "Last Name", "Email", "Age", "Graduated", "Language", "Language"},
						"A2": {"Anakin", "Skywalker", "anakin.skywalker@deathstar.imp", 9, false},
						"A3": {"Obi-Wan", "Kenobi", "obi@jedi.rules", 25, true, "Golang", "Python"},
					},
				),
				schema,
			},
			[]student{
				{"Anakin", "Skywalker", "anakin.skywalker@deathstar.imp", 9, false, nil},
				{"Obi-Wan", "Kenobi", "obi@jedi.rules", 25, true, []string{"Golang", "Python"}},
			},
			2,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ParseXLSXFile(tt.args.in, tt.args.r, tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("ParseXLSXFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			for i, s := range tt.want {
				got := *tt.args.in
				if !reflect.DeepEqual(got[i], s) {
					t.Errorf("ParseXLSXFile() got - %v, want - %v", got[i], s)
				}
			}
		})
	}
}
