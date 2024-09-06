package commonservices

import (
	"strings"
)

func GetFileExtensionFromFileName(file_name string) string {
	file_name_chunks := strings.Split(file_name, ".")
	if len(file_name_chunks) < 1 {
		return ""
	}

	return file_name_chunks[len(file_name_chunks)-1]
}
