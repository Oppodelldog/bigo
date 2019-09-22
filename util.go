package bigo

import (
	"fmt"
	"regexp"
)

func normalizeFileName(name, extension string) string {
	var re = regexp.MustCompile(`[!"§$%&/()=?\\:,'*+~;]`)
	normalizedName := re.ReplaceAllString(name, "")
	normalizedExtension := re.ReplaceAllString(extension, "")
	return fmt.Sprintf("%s.%s", normalizedName, normalizedExtension)
}
