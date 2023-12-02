package ini

// IniSection represents a section in an INI file
type IniSection struct {
	SectionName string
	Keys        []*IniKey
}

// NewIniSection returns a new IniSection with the given section name
func NewIniSection(sectionName string) *IniSection {
	return &IniSection{
		SectionName: sectionName,
		Keys:        make([]*IniKey, 0),
	}
}

// AddKey adds a key to the section, no matter if it already exists. This can create duplicate values
func (s *IniSection) AddKey(keyName string, value string) {
	s.Keys = append(s.Keys, NewIniKey(keyName, value))
}

// AddOrReplaceKey adds a key to the section, or replaces the value of an existing key. Use this if you don't want duplicate keys
func (s *IniSection) AddOrReplaceKey(keyName string, value string) {
	for _, key := range s.Keys {
		if key.KeyName == keyName {
			key.Value = value
			return
		}
	}
	s.Keys = append(s.Keys, NewIniKey(keyName, value))
}

// AddParsedKey adds a key to the section, parsing the keyString into a key value pair
func (s *IniSection) AddParsedKey(keyString string) {
	key := NewParsedIniKey(keyString)

	s.Keys = append(s.Keys, key)
}

// GetKey returns the key with the given name and true, or nil and false if it doesn't exist
func (s *IniSection) GetKey(keyName string) (*IniKey, bool) {
	for _, key := range s.Keys {
		if key.KeyName == keyName {
			return key, true
		}
	}
	return nil, false
}

// AllKeysToStringSlice returns all keys in the section as a slice of strings
func (s *IniSection) AllKeysToStringSlice() []string {
	var keys []string
	for _, key := range s.Keys {
		keys = append(keys, key.ToString())
	}
	return keys
}

// RemoveKey removes the key with the given name from the section
func (s *IniSection) RemoveKey(keyName string) {
	for i, key := range s.Keys {
		if key.KeyName == keyName {
			s.Keys = append(s.Keys[:i], s.Keys[i+1:]...)
			return
		}
	}
}

// RemoveAllKeys removes all keys from the section
func (s *IniSection) RemoveAllKeys() {
	s.Keys = make([]*IniKey, 0)
}

// ToString returns the section as a string in ini format
func (s *IniSection) ToString() string {
	var section string = "[" + s.SectionName + "]\n"
	for _, key := range s.Keys {
		section += key.ToString() + "\n"
	}
	return section
}

// SectionNameToString returns only the section name as a string in ini format
func (s *IniSection) SectionNameToString() string {
	return "[" + s.SectionName + "]"
}
