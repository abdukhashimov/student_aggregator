package parser

var defaultSchema = Schema{
	Version:    "1.0",
	SchemaType: "coords",
	Headers:    true,
	Fields: []FieldSchema{
		{Name: "first_name", Col: "A"},
		{Name: "last_name", Col: "B"},
		{Name: "name", Col: "C"},
	},
}
