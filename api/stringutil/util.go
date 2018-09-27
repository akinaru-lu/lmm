package stringutil

import (
	"strconv"
	"strings"
)

// ParseInt is a shortcut of strconv.Atoi(s)
func ParseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// ParseInt64 is a shortcut of strconv.ParseInt(s, 10, 64)
func ParseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// ParseUint is a shortcut of strconv.ParseUint(s, 10, 32)
func ParseUint(s string) (uint, error) {
	i, err := strconv.ParseUint(s, 10, 32)
	return uint(i), err
}

// ParseUint64 is a shortcut of strconv.ParseUint(s, 10, 64)
func ParseUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

// ReplaceAll is a shortcut of strings.Replace(s, old, new, -1)
func ReplaceAll(s, old, new string) string {
	return strings.Replace(s, old, new, -1)
}

// Uint64ToStr is a shortcut of strconv.FormatUint(i, 10)
func Uint64ToStr(i uint64) string {
	return strconv.FormatUint(i, 10)
}