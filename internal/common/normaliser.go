package common

import "strings"

func NormaliseStringUpper(target string) string {
	return strings.TrimSpace(strings.ToUpper(target))
}

func NormaliseStringLower(target string) string {
	return strings.TrimSpace(strings.ToLower(target))
}

func NormaliseNoSpaces(target string) string {
	return strings.TrimSpace(target)
}

func MapIsNulOrEmpty(target map[string]string) bool {
	return target == nil || len(target) == 0
}
