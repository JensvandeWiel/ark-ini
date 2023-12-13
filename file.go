package ini

import (
	"errors"
)

type IniFile struct {
	AllowedDuplicateKeys []string
	Sections             []*IniSection
}

func NewIniFile(allowedDuplicateKeys ...string) *IniFile {
	return &IniFile{
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
		section = NewIniSection(sectionName, &f.AllowedDuplicateKeys)
		f.Sections = append(f.Sections, section)
	}
	return section
}

// GetKeyFromSection returns the key with the given name from the section with the given name
func (f *IniFile) GetKeyFromSection(sectionName string, keyName string) (*IniKey, error) {
	section, exists := f.GetSection(sectionName)
	if exists {
		key, exists := section.FindKey(keyName)
		if exists {
			return key, nil
		}
		return nil, errors.New("key not found")
	}
	return nil, errors.New("section not found")
}

// GetKeyFromSectionWithMultipleValues returns all the keys with the given name from the section with the given name
func (f *IniFile) GetKeyFromSectionWithMultipleValues(sectionName string, keyName string) ([]*IniKey, error) {
	section, exists := f.GetSection(sectionName)
	if exists {
		return section.FindKeys(keyName)
	}
	return nil, errors.New("section not found")
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

// SafelyAddKeyToSection same as the others but will automatically check if duplicates are allowed, if so it will add the key, if not it will replace it.
func (f *IniFile) SafelyAddKeyToSection(sectionName string, keyName string, value interface{}) {
	if f.duplicateAllowed(keyName) {
		f.AddKeyToSection(sectionName, keyName, value)
	} else {
		f.UpdateOrCreateKeyInSection(sectionName, keyName, value)

	}
}

// AddKeyToSection adds a key to the section with the given name, if the section does not exist it is created
func (f *IniFile) AddKeyToSection(sectionName string, keyName string, value interface{}) {
	section := f.GetOrCreateSection(sectionName)
	section.AddKey(keyName, value)
}

// UpdateOrCreateKeyInSection updates the value of the key with the given name in the section with the given name, or creates a new key with the given name and value in the section. If the section does not exist it is created.
func (f *IniFile) UpdateOrCreateKeyInSection(sectionName string, keyName string, value interface{}) {
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

func (f *IniFile) duplicateAllowed(key string) bool {
	for _, allowedKey := range f.AllowedDuplicateKeys {
		if key == allowedKey {
			return true
		}
	}
	return false
}
