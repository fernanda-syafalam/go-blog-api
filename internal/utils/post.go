package utils

import (
	"strings"
)

func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.Replace(slug, " ", "-", -1)
	return slug
}
