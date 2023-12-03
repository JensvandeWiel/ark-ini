package ini

import (
	"errors"
	"fmt"
	"strings"
)

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

type ContainerKeyValue struct {
	Key   string
	Value interface{}
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

	/*if len(inputSlice) > 1 {
		return "(" + strings.Join(parts, ",") + ")"
	}*/
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

//region helpers

// ToString returns the values of the container as a string e.g. "EngramClassName="EngramEntry_CryoGun_Mod_C",EngramHidden=True,EngramPointsCost=0,EngramLevelRequirement=90,RemoveEngramPreReq=False"
func (c *IniContainer) ToString() string {
	return serializeToContainerKV(c.KeyValues)
}

//endregion
