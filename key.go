package ini

import "strings"

// IniKey represents a key in an INI file
type IniKey struct {
	KeyName string
	Value   string
}

// ToString returns the key as a string in ini format
func (k *IniKey) ToString() string {
	return k.KeyName + "=" + k.Value
}

// NewParsedIniKey returns a key, parsing the keyString into a key value pair
func NewParsedIniKey(keyString string) *IniKey {
	//Check if arg is empty or null
	if keyString == "" || strings.TrimSpace(keyString) == "" {
		return nil
	}

	// Split keyString into key and value
	splitKeyString := strings.SplitN(keyString, "=", 2)
	if strings.TrimSpace(splitKeyString[0]) == "" {
		return nil
	}

	key := &IniKey{KeyName: splitKeyString[0]}

	if len(splitKeyString) > 1 {
		key.Value = splitKeyString[1]
	}

	return key
}

// NewIniKey returns a new IniKey with the given key and value
func NewIniKey(keyName string, keyValue string) *IniKey {
	return &IniKey{KeyName: keyName, Value: keyValue}
}
