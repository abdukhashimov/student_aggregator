package ports

import (
	"context"
)

type AggregatorService interface {
	ParseFile(ctx context.Context, FileName string, SchemaID string) error
}
