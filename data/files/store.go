package files

import (
	"runtime"
	"strings"
)

func GetLocation() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "./data/files"
	}

	filename = strings.Trim(filename, "store.go")
	return filename
}
