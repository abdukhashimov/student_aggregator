package domain

type ParseFileInput struct {
	FileName string `json:"file_name" validate:"required"`
	SchemaID string `json:"schema_id" validate:"required"`
}
