package parser

import (
	"github.com/xuri/excelize/v2"
	"reflect"
	"testing"
)

func getFileBytes(t *testing.T, sheets []string, data map[string][]string) *[]byte {
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
	b := buffer.Bytes()
	return &b
}

func TestParseFile(t *testing.T) {
	var s []interface{}
	type args struct {
		l *[]interface{}
		b *[]byte
		s Schema
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
				&s,
				getFileBytes(
					t,
					[]string{"Sheet1"},
					map[string][]string{
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
				&s,
				getFileBytes(
					t,
					[]string{"Sheet1", "Sheet2"},
					map[string][]string{
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
			if err := ParseXLSXFile(tt.args.l, tt.args.b, tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("ParseXLSXFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(s) != tt.wantLen {
				t.Errorf("ParseXLSXFile() len = %v, wantLen %v", len(s), tt.wantLen)
			}

			t.Cleanup(func() {
				s = s[:0]
			})
		})
	}
}

func Test_mapStruct(t *testing.T) {
	type args struct {
		s        any
		valueMap map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    any
	}{
		{
			name: "map tagged struct",
			args: args{
				&struct {
					FirstName string `schema:"first_name"`
					LastName  string `schema:"last_name"`
					Email     string `schema:"email"`
				}{},
				map[string]string{
					"first_name": "Anakin",
					"last_name":  "Skywalker",
					"email":      "anakin.skywalker@deathstar.imp",
				},
			},
			wantErr: false,
			want: &struct {
				FirstName string `schema:"first_name"`
				LastName  string `schema:"last_name"`
				Email     string `schema:"email"`
			}{"Anakin", "Skywalker", "anakin.skywalker@deathstar.imp"},
		},
		{
			name: "cant map unexported field",
			args: args{
				&struct {
					FirstName string `schema:"first_name"`
					LastName  string `schema:"last_name"`
					email     string `schema:"email"`
				}{},
				map[string]string{
					"first_name": "Anakin",
					"last_name":  "Skywalker",
					"email":      "anakin.skywalker@deathstar.imp",
				},
			},
			wantErr: true,
			want: &struct {
				FirstName string `schema:"first_name"`
				LastName  string `schema:"last_name"`
				email     string `schema:"email"`
			}{"Anakin", "Skywalker", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mapStruct(tt.args.s, tt.args.valueMap); (err != nil) != tt.wantErr {
				t.Errorf("mapStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.s, tt.want) {
				t.Errorf("mapStruct() sturct = %v, want %v", tt.args.s, tt.want)
			}
		})
	}
}

func TestParseCSVFile(t *testing.T) {
	var s []interface{}
	type args struct {
		l *[]interface{}
		b *[]byte
		s Schema
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantLen int
	}{
		{
			name: "parse csv file with headers",
			args: args{
				&s,
				&[]byte{
					70, 105, 114, 115, 116, 32, 78, 97, 109, 101, 44, 76, 97, 115, 116, 32, 78, 97, 109, 101, 44,
					69, 109, 97, 105, 108, 10, 65, 110, 97, 107, 105, 110, 44, 32, 83, 107, 121, 119, 97, 108, 107, 101,
					114, 44, 32, 97, 110, 97, 107, 105, 110, 46, 115, 107, 121, 119, 97, 108, 107, 101, 114, 64, 100,
					101, 97, 116, 104, 115, 116, 97, 114, 46, 105, 109, 112, 10, 79, 98, 105, 45, 87, 97, 110, 44, 75,
					101, 110, 111, 98, 105, 44, 111, 98, 105, 64, 106, 101, 100, 105, 46, 114, 117, 108, 101, 115},
				defaultSchema,
			},
			wantErr: false,
			wantLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ParseCSVFile(tt.args.l, tt.args.b, tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("ParseCSVFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(s) != tt.wantLen {
				t.Errorf("ParseXLSXFile() len = %v, wantLen %v", len(s), tt.wantLen)
			}

			t.Cleanup(func() {
				s = s[:0]
			})
		})
	}
}
