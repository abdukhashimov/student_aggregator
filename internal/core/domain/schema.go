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
	Name       string        `json:"name" validate:"required,min=3"`
	Version    string        `json:"version" validate:"required"`
	SchemaType string        `json:"schema_type" validate:"required"`
	Headers    bool          `json:"headers" validate:"required"`
	Fields     []FieldSchema `json:"fields" validate:"required"`
}

type UpdateSchemaInput struct {
	Name       *string        `json:"name" bson:"name,omitempty" validate:"omitempty,required,min=3"`
	Slug       *string        `json:"-" bson:"slug,omitempty" validate:"omitempty,required"`
	Version    *string        `json:"version" bson:"version,omitempty" validate:"omitempty,required"`
	SchemaType *string        `json:"schema_type" bson:"schema_type,omitempty" validate:"omitempty,required"`
	Headers    *bool          `json:"headers" bson:"headers,omitempty" validate:"omitempty,required"`
	Fields     *[]FieldSchema `json:"fields" bson:"fields,omitempty" validate:"omitempty,required"`
}
