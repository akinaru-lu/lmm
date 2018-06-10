package strings

import (
	"strconv"
)

func StrToUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}
