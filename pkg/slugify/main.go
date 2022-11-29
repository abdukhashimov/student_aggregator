package slugify

import (
	"fmt"

	"github.com/gosimple/slug"
	"github.com/teris-io/shortid"
)

const (
	maxFilenameLength = 10
)

func GenSlug(filename string) string {
	return slug.Make(fileNameWithMaxLength((filename)))
}

func GenSlugWithID(filename string) string {
	id := shortid.MustGenerate()

	return fmt.Sprintf("%s-%s", fileNameWithMaxLength(filename), id)
}

func fileNameWithMaxLength(filename string) string {
	if len(filename) > maxFilenameLength {
		return filename[:maxFilenameLength]
	}

	return filename
}
