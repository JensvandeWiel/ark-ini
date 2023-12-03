package ini

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//region IniContainer

type IniContainer struct {
	KeyValues []ContainerKeyValue
}

// NewIniContainerFromString returns a new IniContainer. The input string must be '
func NewIniContainerFromString(inputString string) (IniContainer, error) {
	keyValues, err := deserializeToContainerKv(inputString)
	if err != nil {
		return IniContainer{}, err
	}

	return IniContainer{KeyValues: keyValues}, nil
}

//endregion

//region ContainerKeyValue

type ContainerKeyValue struct {
	Key   string
	Value interface{}
}

func (c *ContainerKeyValue) ToString() string {
	return fmt.Sprintf("%s=%v", c.Key, c.Value)
}

//region Key Conversions

// AsString returns the key value as a string
func (c *ContainerKeyValue) AsString() (string, error) {
	value, ok := c.Value.(string)
	if !ok {
		return "", errors.New("key value is not a string")
	}
	return value, nil
}

// AsInt returns the key value as an int
func (c *ContainerKeyValue) AsInt() (int, error) {
	strValue, ok := c.Value.(string)
	if !ok {
		return 0, errors.New("key value is not a string")
	}

	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		return 0, errors.New("failed to convert key value to int")
	}

	return intValue, nil
}

// AsFloat64 returns the key value as a float64
func (c *ContainerKeyValue) AsFloat64() (float64, error) {
	strValue, ok := c.Value.(string)
	if !ok {
		return 0, errors.New("key value is not a string")
	}

	floatValue, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		return 0, errors.New("failed to convert key value to float64")
	}

	return floatValue, nil
}

// AsBool returns the key value as a bool
func (c *ContainerKeyValue) AsBool() (bool, error) {
	strValue, ok := c.Value.(string)
	if !ok {
		return false, errors.New("key value is not a string")
	}

	boolValue, err := strconv.ParseBool(strValue)
	if err != nil {
		return false, errors.New("failed to convert key value to bool")
	}

	return boolValue, nil
}

// AsContainer returns the key value as a container
func (c *ContainerKeyValue) AsContainer() (IniContainer, error) {
	strValue, ok := c.Value.(string)
	if !ok {
		return IniContainer{}, errors.New("key value is not a string")
	}

	if strings.HasPrefix(strValue, "(") && strings.HasSuffix(strValue, ")") {
		return NewIniContainerFromString(strValue)
	}

	return IniContainer{}, errors.New("key value is not a container")
}

// AsGuessedValue returns the key value as a guessed value and the value type
func (c *ContainerKeyValue) AsGuessedValue() (interface{}, KeyType, error) {
	switch v := c.Value.(type) {
	case string:
		if strings.HasPrefix(v, "(") && strings.HasSuffix(v, ")") {
			cont, err := NewIniContainerFromString(v)
			if err == nil {
				return cont, Container, nil
			}
			return nil, Fail, err
		} else if val, err := strconv.ParseFloat(v, 64); err == nil && hasDecimal(val) {
			return val, Float64, nil
		} else if val, err := strconv.Atoi(v); err == nil {
			return val, Int, nil
		} else if val, err := strconv.ParseBool(v); err == nil {
			return val, Boolean, nil
		} else {
			return v, String, nil
		}
	case []ContainerKeyValue:
		return v, Container, nil
	}

	return nil, Fail, errors.New("unsupported type")
}

//endregion

//endregion

//region helpers

// ToString returns the values of the container as a string e.g. "EngramClassName="EngramEntry_CryoGun_Mod_C",EngramHidden=True,EngramPointsCost=0,EngramLevelRequirement=90,RemoveEngramPreReq=False"
func (c *IniContainer) ToString() string {
	return serializeToContainerKV(c.KeyValues)
}

// serializeToContainerKV serializes a slice of key-value pairs to a string
func serializeToContainerKV(inputSlice []ContainerKeyValue) string {
	var parts []string

	for _, kv := range inputSlice {
		if nestedSlice, ok := kv.Value.([]ContainerKeyValue); ok {
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
func deserializeToContainerKv(inputString string) ([]ContainerKeyValue, error) {
	var result []ContainerKeyValue

	if inputString == "" || strings.TrimSpace(inputString) == "" {
		return nil, errors.New("inputString is empty")
	}

	var parts []string

	if strings.HasPrefix(inputString, "(") && strings.HasSuffix(inputString, ")") {
		//it is only the value

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
			result = append(result, ContainerKeyValue{Key: key, Value: nestedSlice})
		} else {
			// If the value is a simple value, store it directly
			result = append(result, ContainerKeyValue{Key: key, Value: strings.TrimSpace(kv[1])})
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
