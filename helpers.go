package ini

import (
	"math"
	"path/filepath"
)

func removeFileExtension(filename string) string {
	return filepath.Base(filename[:len(filename)-len(filepath.Ext(filename))])
}

func hasDecimal(floatNum float64) bool {
	return floatNum != math.Floor(floatNum)
}
