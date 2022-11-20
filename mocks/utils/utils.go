package utils

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
)

type contextKey string

const (
	withErrorKey contextKey = "withError"
)

func WithError(ctx context.Context) bool {
	withError, ok := ctx.Value(withErrorKey).(bool)

	if !ok {
		return false
	}

	return withError
}

func SetWithError(r *http.Request, withError bool) *http.Request {
	ctx := context.WithValue(r.Context(), withErrorKey, withError)

	return r.WithContext(ctx)
}

func CopySchema(schema *domain.Schema) *domain.Schema {
	schemaCopy := *schema
	schemaCopy.Fields = make([]domain.FieldSchema, len(schema.Fields))
	copy(schemaCopy.Fields, schema.Fields)

	return &schemaCopy
}

func GetSlug(in string) string {
	space := regexp.MustCompile(`\s+`)
	result := space.ReplaceAllString(in, " ")
	result = strings.TrimSpace(result)
	result = space.ReplaceAllString(result, "-")
	result = strings.ToLower(result)

	return result
}

func IncrementMongoId(in string) string {
	value, _ := strconv.ParseInt(in[15:], 16, 64)
	value += 1

	return fmt.Sprintf("%s%09x", in[:15], value)
}
