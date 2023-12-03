package ini

import (
	"errors"
	"strconv"
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
	KeyName string
	Value   string
}

// ToString returns the key as a string in ini format
func (k *IniKey) ToString() string {
	return k.KeyName + "=" + k.Value
}

// ToContainerString returns the key value as a container string e.g. "OverrideNamedEngramEntries=(EngramClassName="EngramEntry_CryoGun_Mod_C",EngramHidden=True,EngramPointsCost=0,EngramLevelRequirement=90,RemoveEngramPreReq=False)"
func (k *IniKey) ToContainerString() (string, error) {
	if strings.HasPrefix(k.Value, "(") && strings.HasSuffix(k.Value, ")") {

		cont, err := NewIniContainerFromString(k.ToString())

		if err != nil {
			return "", err
		}

		return cont.ToString(), nil
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

//region Key Conversions

// AsString returns the key value as a string
func (k *IniKey) AsString() string {
	return k.Value
}

// AsInt returns the key value as an int
func (k *IniKey) AsInt() (int, error) {
	return strconv.Atoi(k.Value)
}

// AsFloat64 returns the key value as an int64
func (k *IniKey) AsFloat64() (float64, error) {
	return strconv.ParseFloat(k.Value, 64)
}

func (k *IniKey) AsBool() (bool, error) {
	return strconv.ParseBool(k.Value)
}

// AsContainer returns the key value as a container
func (k *IniKey) AsContainer() (IniContainer, error) {
	if strings.HasPrefix(k.Value, "(") && strings.HasSuffix(k.Value, ")") {
		return NewIniContainerFromString(k.Value)
	} else {
		return IniContainer{}, errors.New("key value is not a container")
	}
}

// AsGuessedValue returns the key value as a guessed value and the value type
func (k *IniKey) AsGuessedValue() (interface{}, KeyType) {
	if strings.HasPrefix(k.Value, "(") && strings.HasSuffix(k.Value, ")") {
		cont, err := NewIniContainerFromString(k.Value)
		if err != nil {
			return nil, Fail
		}
		return cont, Container
	} else if val, err := strconv.ParseFloat(k.Value, 64); err == nil && hasDecimal(val) {
		return val, Float64
	} else if val, err := strconv.Atoi(k.Value); err == nil {
		return val, Int
	} else if val, err := strconv.ParseBool(k.Value); err == nil {
		return val, Boolean
	} else {
		return k.Value, String
	}

}

//endregion
