package domain

type Schema struct {
	ID         string        `json:"id" bson:"_id,omitempty"`
	Name       string        `json:"name" bson:"name"`
	Slug       string        `json:"slug" bson:"slug"`
	Version    string        `json:"version" bson:"version"`
	SchemaType string        `json:"schema_type" bson:"schema_type"`
	Headers    bool          `json:"headers" bson:"headers"`
	Fields     []FieldSchema `json:"fields"  bson:"fields"`
}

type FieldSchema struct {
	Col  string `json:"col" bson:"col"`
	Name string `json:"name" bson:"name"`
}

type NewSchemaInput struct {
	Name       string        `json:"name"`
	Version    string        `json:"version"`
	SchemaType string        `json:"schema_type"`
	Headers    bool          `json:"headers"`
	Fields     []FieldSchema `json:"fields"`
}

type UpdateSchemaInput struct {
	Name       *string        `json:"name"`
	Slug       *string        `json:"slug"`
	Version    *string        `json:"version"`
	SchemaType *string        `json:"schema_type"`
	Headers    *bool          `json:"headers"`
	Fields     *[]FieldSchema `json:"fields"`
}
