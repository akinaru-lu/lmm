package strings

import (
	"strconv"
	"strings"
)

func Join(sep string, target []string) string {
	return strings.Join(target, sep)
}

func Uint64ToStr(i uint64) string {
	return strconv.FormatUint(i, 10)
}

func StrToInt(s string) (int, error) {
	n, err := StrToUint64(s)
	return int(n), err
}

func StrToUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

func ReplaceAll(s, old, new string) string {
	return strings.Replace(s, old, new, -1)
}

func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}
