package helper

import "strings"

func NormalizeEmail(email string) string {

	trimmed := strings.TrimSpace(email)

	return strings.ToLower(trimmed)
}
