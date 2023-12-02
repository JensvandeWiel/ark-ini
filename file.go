package ini

import "path/filepath"

type IniFile struct {
	fileName             string
	FilePath             string
	AllowedDuplicateKeys []string
	Sections             []*IniSection
}

func NewIniFile(path string, allowedDuplicateValues ...string) *IniFile {
	fileName := removeFileExtension(filepath.Base(path))

	return &IniFile{
		fileName:             fileName,
		FilePath:             path,
		AllowedDuplicateKeys: allowedDuplicateValues,
		Sections:             make([]*IniSection, 0),
	}
}
