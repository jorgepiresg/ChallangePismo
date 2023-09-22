package utils

import "strings"

func CleanDocument(document string) string {
	document = strings.ReplaceAll(document, "-", "")
	document = strings.ReplaceAll(document, ".", "")
	return document
}
