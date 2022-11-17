package domain

type Schema struct {
	ID         string        `json:"id" bson:"_id,omitempty"`
	Name       string        `json:"name" bson:"name"`
	Slug       string        `json:"-" bson:"slug"`
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
	Name       *string        `json:"name" bson:"name,omitempty"`
	Slug       *string        `json:"-" bson:"slug,omitempty"`
	Version    *string        `json:"version" bson:"version,omitempty"`
	SchemaType *string        `json:"schema_type" bson:"schema_type,omitempty"`
	Headers    *bool          `json:"headers" bson:"headers,omitempty"`
	Fields     *[]FieldSchema `json:"fields" bson:"fields,omitempty"`
}
