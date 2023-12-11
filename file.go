package ini

import (
	"path/filepath"
)

//TODO add parsing of values and save them as the type and not a string

type IniFile struct {
	fileName             string
	FilePath             string
	AllowedDuplicateKeys []string
	Sections             []*IniSection
}

func NewIniFile(path string, allowedDuplicateKeys ...string) IniFile {
	fileName := removeFileExtension(filepath.Base(path))
	return IniFile{
		fileName:             fileName,
		FilePath:             path,
		AllowedDuplicateKeys: allowedDuplicateKeys,
		Sections:             make([]*IniSection, 0),
	}
}

// GetSection returns the section with the given name and true, or nil and false if it doesn't exist
//
// Returns:
//
//	*IniSection - The gotten section if it exists else it is nil.
//
//	bool - True if the section exists else false.
func (f *IniFile) GetSection(sectionName string) (*IniSection, bool) {
	for _, section := range f.Sections {
		if section.SectionName == sectionName {
			return section, true
		}
	}
	return nil, false
}

// GetOrCreateSection returns the section with the given name if it exists, or creates a new section with the given name and returns it
func (f *IniFile) GetOrCreateSection(sectionName string) *IniSection {
	section, exists := f.GetSection(sectionName)
	if !exists {
		section = NewIniSection(sectionName)
		f.Sections = append(f.Sections, section)
	}
	return section
}

// RemoveSection removes the section with the given name from the file
func (f *IniFile) RemoveSection(sectionName string) {
	for i, section := range f.Sections {
		if section.SectionName == sectionName {
			f.Sections = append(f.Sections[:i], f.Sections[i+1:]...)
			return
		}
	}
}

// RemoveAllSections removes all sections from the file
func (f *IniFile) RemoveAllSections() {
	f.Sections = make([]*IniSection, 0)
}

// AddKeyToSection adds a key to the section with the given name, if the section does not exist it is created
func (f *IniFile) AddKeyToSection(sectionName string, keyName string, value string) {
	section := f.GetOrCreateSection(sectionName)
	section.AddKey(keyName, value)
}

// UpdateOrCreateKeyInSection updates the value of the key with the given name in the section with the given name, or creates a new key with the given name and value in the section. If the section does not exist it is created.
func (f *IniFile) UpdateOrCreateKeyInSection(sectionName string, keyName string, value string) {
	section := f.GetOrCreateSection(sectionName)
	section.AddOrReplaceKey(keyName, value)
}

// RemoveKeyFromSection removes the key with the given name from the section with the given name
func (f *IniFile) RemoveKeyFromSection(sectionName string, keyName string) {
	section, exists := f.GetSection(sectionName)
	if exists {
		section.RemoveKey(keyName)
	}
}

// RemoveMultipleKeysFromSection removes all keys with the given name from the section with the given name
func (f *IniFile) RemoveMultipleKeysFromSection(sectionName string, keyName string) {
	section, exists := f.GetSection(sectionName)
	if exists {
		section.RemoveMultipleKey(keyName)
	}
}

// RemoveAllKeysFromSection removes all keys from the section with the given name
func (f *IniFile) RemoveAllKeysFromSection(sectionName string) {
	section, exists := f.GetSection(sectionName)
	if exists {
		section.RemoveAllKeys()
	}
}

// ToString returns the ini file as a string
func (f *IniFile) ToString() string {
	var file string
	for _, section := range f.Sections {
		file += section.ToString()
	}
	return file
}
