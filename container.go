package ini

import (
	"errors"
	"fmt"
	"strings"
)

//region IniContainer

type IniContainer struct {
	KeyValues []ContainerKey
}

// NewIniContainerFromString returns a new IniContainer. The input string must be '
func NewIniContainerFromString(inputString string) (IniContainer, error) {
	keyValues, err := deserializeToContainerKv(inputString)
	if err != nil {
		return IniContainer{}, err
	}

	return IniContainer{KeyValues: keyValues}, nil
}

func NewIniContainerFromSlice(inputSlice []ContainerKey) IniContainer {
	return IniContainer{KeyValues: inputSlice}
}

// FindKey returns the key with the given name and true, or nil and false if it doesn't exist
func (c *IniContainer) FindKey(keyName string) (*ContainerKey, bool) {
	for _, key := range c.KeyValues {
		if key.Key == keyName {
			return &key, true
		}
	}
	return nil, false
}

//endregion

//region ContainerKey

type ContainerKey struct {
	Key   string
	Value interface{}
}

func (c *ContainerKey) ToString() string {
	return fmt.Sprintf("%s=%v", c.Key, c.Value)
}

// ToValueString returns the key's value as a string
func (c *ContainerKey) ToValueString() string {
	if container, ok := c.Value.(IniContainer); ok {
		return container.ToString()
	} else {
		return fmt.Sprintf("%v", c.Value)
	}
}

//region Key Conversions

// AsString returns the key value as a string
func (c *ContainerKey) AsString() (string, error) {

	if value, ok := c.Value.(string); ok {
		return value, nil

	} else {
		return "", errors.New("key value is not a string")
	}
}

// AsInt returns the key value as an int
func (c *ContainerKey) AsInt() (int, error) {
	if value, ok := c.Value.(int); ok {
		return value, nil
	} else {
		return -1, errors.New("key value is not an int")
	}
}

// AsFloat64 returns the key value as a float64
func (c *ContainerKey) AsFloat64() (float64, error) {
	if value, ok := c.Value.(float64); ok {
		return value, nil
	} else {
		return -1, errors.New("key value is not a float64")
	}
}

// AsBool returns the key value as a bool
func (c *ContainerKey) AsBool() (bool, error) {
	if value, ok := c.Value.(bool); ok {
		return value, nil
	} else {
		return false, errors.New("key value is not a bool")
	}
}

// AsContainer returns the key value as a container
func (c *ContainerKey) AsContainer() (IniContainer, error) {
	if container, ok := c.Value.(IniContainer); ok {
		return container, nil
	} else if container, ok := c.Value.([]ContainerKey); ok {
		return IniContainer{KeyValues: container}, nil
	} else {
		return IniContainer{}, errors.New("key value is not a container")
	}
}

// AsGuessedValue returns the key value as a guessed value and the value type
func (c *ContainerKey) AsGuessedValue() (interface{}, KeyType, error) {

	switch c.Value.(type) {
	case string:
		return c.Value, String, nil
	case int:
		return c.Value, Int, nil
	case float64:
		return c.Value, Float64, nil
	case bool:
		return c.Value, Boolean, nil
	case IniContainer:
		return c.Value, Container, nil
	case []ContainerKey:
		return IniContainer{KeyValues: c.Value.([]ContainerKey)}, Container, nil
	default:
		return nil, Fail, errors.New("unknown key type")
	}
}

//endregion

//endregion

//region helpers

// ValueToString returns the values of the container as a string e.g. "EngramClassName="EngramEntry_CryoGun_Mod_C",EngramHidden=True,EngramPointsCost=0,EngramLevelRequirement=90,RemoveEngramPreReq=False"
func (c *IniContainer) ValueToString() string {
	return serializeToContainerKV(c.KeyValues)
}

// ToString returns the values of the container as a string e.g. "EngramClassName="EngramEntry_CryoGun_Mod_C",EngramHidden=True,EngramPointsCost=0,EngramLevelRequirement=90,RemoveEngramPreReq=False"
func (c *IniContainer) ToString() string {
	return "(" + serializeToContainerKV(c.KeyValues) + ")"
}

// serializeToContainerKV serializes a slice of key-value pairs to a string
func serializeToContainerKV(inputSlice []ContainerKey) string {
	var parts []string

	for _, kv := range inputSlice {
		if nestedSlice, ok := kv.Value.([]ContainerKey); ok {
			// Recursively serialize nested slices
			nestedStr := serializeToContainerKV(nestedSlice)
			parts = append(parts, fmt.Sprintf("%s=(%s)", kv.Key, nestedStr))
		} else {
			// serializeToContainerKV simple values
			parts = append(parts, fmt.Sprintf("%s=%v", kv.Key, kv.Value))
		}
	}

	return strings.Join(parts, ",")
}

// deserializeToContainerKv deserializes a string to a slice of key-value pairs
func deserializeToContainerKv(inputString string) ([]ContainerKey, error) {
	var result []ContainerKey

	if inputString == "" || strings.TrimSpace(inputString) == "" {
		return nil, errors.New("inputString is empty")
	}

	var parts []string

	if strings.HasPrefix(inputString, "(") && strings.HasSuffix(inputString, ")") {
		//it is only the value of a "container" key
		inputString = strings.Trim(inputString, "()")

		parts = splitInputString(inputString)

	} else {
		parts = splitInputString(inputString)
	}

	// Process each part to build the result
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		key := strings.TrimSpace(kv[0])

		if strings.Contains(kv[1], "(") {
			// If the value contains '(', it means it's a nested structure
			nestedStr := strings.Trim(kv[1], "()")
			nestedSlice, err := deserializeToContainerKv(nestedStr)
			if err != nil {
				return nil, err
			}
			result = append(result, ContainerKey{Key: key, Value: nestedSlice})
		} else {
			// If the value is a simple value, store it directly
			result = append(result, ContainerKey{Key: key, Value: toGuessedType(strings.TrimSpace(kv[1]))})
		}
	}

	return result, nil
}

// splitInputString splits the input string based on parentheses and commas
func splitInputString(inputString string) []string {
	var parts []string
	var currentPart strings.Builder
	openParentheses := 0

	for _, char := range inputString {
		switch char {
		case '(':
			openParentheses++
		case ')':
			openParentheses--
		case ',':
			if openParentheses == 0 {
				// Split only if not inside parentheses
				parts = append(parts, currentPart.String())
				currentPart.Reset()
				continue
			}
		}

		currentPart.WriteRune(char)
	}

	// Add the last part
	if currentPart.Len() > 0 {
		parts = append(parts, currentPart.String())
	}

	return parts
}

//endregion
