package ini

import (
	"errors"
	"fmt"
	"strings"
)

type KeyType string

const (
	Container KeyType = "container"
	String    KeyType = "string"
	Int       KeyType = "int"
	Float64   KeyType = "float64"
	Boolean   KeyType = "bool"
	Fail      KeyType = "fail"
)

// IniKey represents a key in an INI file
type IniKey struct {
	Key   string
	Value interface{}
}

// ToString returns the key as a string in ini format
func (k *IniKey) ToString() string {
	return fmt.Sprintf("%s=%s", k.Key, k.ToValueString())
}

// ToValueString returns the key's value as a string
func (k *IniKey) ToValueString() string {
	if container, ok := k.Value.(IniContainer); ok {
		return container.ToString()
	} else {
		return fmt.Sprintf("%v", k.Value)
	}
}

// ToContainerString returns the key value as a container string e.g. "OverrideNamedEngramEntries=(EngramClassName="EngramEntry_CryoGun_Mod_C",EngramHidden=True,EngramPointsCost=0,EngramLevelRequirement=90,RemoveEngramPreReq=False)"
func (k *IniKey) ToContainerString() (string, error) {
	if container, ok := k.Value.(IniContainer); ok {
		return container.ToString(), nil
	} else {
		return "", errors.New("key value is not a container")
	}
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

	key := NewIniKey(splitKeyString[0], "")

	if len(splitKeyString) > 1 {
		key.Value = toGuessedType(splitKeyString[1])
	}

	return key
}

// NewIniKey returns a new IniKey with the given key and value
func NewIniKey(keyName string, keyValue interface{}) *IniKey {
	return &IniKey{Key: keyName, Value: keyValue}
}

//region Key Conversions

// AsString returns the key value as a string
func (k *IniKey) AsString() (string, error) {
	if s, ok := k.Value.(string); ok {
		return s, nil
	} else {
		return "", errors.New("key value is not a string")
	}
}

// AsInt returns the key value as an int
func (k *IniKey) AsInt() (int, error) {
	if i, ok := k.Value.(int); ok {
		return i, nil
	} else {
		return -1, errors.New("key value is not an int")
	}
}

// AsFloat64 returns the key value as an int64
func (k *IniKey) AsFloat64() (float64, error) {
	if f, ok := k.Value.(float64); ok {
		return f, nil
	} else {
		return -1, errors.New("key value is not a float64")
	}
}

// AsBool returns the key value as a bool
func (k *IniKey) AsBool() (bool, error) {
	if b, ok := k.Value.(bool); ok {
		return b, nil
	} else {
		return false, errors.New("key value is not a bool")
	}
}

// AsContainer returns the key value as a container
func (k *IniKey) AsContainer() (IniContainer, error) {
	if container, ok := k.Value.(IniContainer); ok {
		return container, nil
	} else {
		return IniContainer{}, errors.New("key value is not a container")
	}
}

// AsGuessedValue returns the key value as a guessed value and the value type
func (k *IniKey) AsGuessedValue() (interface{}, KeyType) {
	switch k.Value.(type) {
	case string:
		return k.Value.(string), String
	case int:
		return k.Value.(int), Int
	case float64:
		return k.Value.(float64), Float64
	case bool:
		return k.Value.(bool), Boolean
	case IniContainer:
		return k.Value.(IniContainer), Container
	default:
		return nil, Fail
	}
}

//endregion
