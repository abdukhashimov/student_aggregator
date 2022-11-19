package repository

import (
	"fmt"
	"strconv"
)

func IncrementId(in string) string {
	value, _ := strconv.ParseInt(in[15:], 16, 64)
	value += 1

	return fmt.Sprintf("%s%09x", in[:15], value)
}
