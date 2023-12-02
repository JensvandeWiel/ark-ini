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

//region Key & Value adding

// AddKey adds a key no matter if it already exists. (May result in duplicate keys) (it will take the first key found if there are more)
func (s *IniSection) AddKey(keyName string, value string) {
	s.Keys = append(s.Keys, NewIniKey(keyName, value))
}

// AddOrReplaceKey adds a key if it not exists otherwise it will replace it (it will take the first key found if there are more) (Use this to avoid duplicate keys)
func (s *IniSection) AddOrReplaceKey(keyName string, value string) {
	for _, key := range s.Keys {
		if key.KeyName == keyName {
			key.Value = value
			return
		}
	}
	s.Keys = append(s.Keys, NewIniKey(keyName, value))
}

// AddParsedKey adds a key from a string like “key=value”, no matter if it already exists (it will take the first key found if there are more)
func (s *IniSection) AddParsedKey(keyString string) {
	key := NewParsedIniKey(keyString)
	s.AddKey(key.KeyName, key.Value)
}

// AddOrReplaceParsedKey adds a key from a string like “key=value” if it does not exist otherwise it will replace it (it will take the first key found if there are more) (Use this to avoid duplicate keys)
func (s *IniSection) AddOrReplaceParsedKey(keyString string) {
	key := NewParsedIniKey(keyString)
	s.AddOrReplaceKey(key.KeyName, key.Value)
}

//endregion

//region Getting keys

// GetKey returns the key with the given name and true, or nil and false if it doesn't exist
func (s *IniSection) GetKey(keyName string) (*IniKey, bool) {
	for _, key := range s.Keys {
		if key.KeyName == keyName {
			return key, true
		}
	}
	return nil, false
}

// GetMultipleKeys gets all the keys with the same name
func (s *IniSection) GetMultipleKeys(keyName string) []*IniKey {
	var keys []*IniKey
	for _, key := range s.Keys {
		if key.KeyName == keyName {
			keys = append(keys, key)
		}
	}
	return keys
}

//endregion

//region Removing keys

// RemoveKey removes the key with the given name from the section
func (s *IniSection) RemoveKey(keyName string) {
	for i, key := range s.Keys {
		if key.KeyName == keyName {
			s.Keys = append(s.Keys[:i], s.Keys[i+1:]...)
			return
		}
	}
}

// RemoveMultipleKey removes all the keys with the same KeyName
func (s *IniSection) RemoveMultipleKey(keyName string) {
	for i, key := range s.Keys {
		if key.KeyName == keyName {
			s.Keys = append(s.Keys[:i], s.Keys[i+1:]...)
		}
	}
}

// RemoveAllKeys removes all keys from the section
func (s *IniSection) RemoveAllKeys() {
	s.Keys = make([]*IniKey, 0)
}

//endregion

//region Helpers

// AllKeysToStringSlice returns all keys in the section as a slice of strings
func (s *IniSection) AllKeysToStringSlice() []string {
	var keys []string
	for _, key := range s.Keys {
		keys = append(keys, key.ToString())
	}
	return keys
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

// CheckForMultipleKeys returns the number of keys with the given name in the section.
func (s *IniSection) CheckForMultipleKeys(keyName string) int {
	var count = 0
	for _, key := range s.Keys {
		if key.KeyName == keyName {
			count++
		}
	}
	return count
}

//endregion
