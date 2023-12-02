package ini

import "path/filepath"

func removeFileExtension(filename string) string {
	return filepath.Base(filename[:len(filename)-len(filepath.Ext(filename))])
}
