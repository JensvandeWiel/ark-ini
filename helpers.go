package ini

import (
	"math"
	"path/filepath"
	"strconv"
	"strings"
)

func removeFileExtension(filename string) string {
	return filepath.Base(filename[:len(filename)-len(filepath.Ext(filename))])
}

func hasDecimal(floatNum float64) bool {
	return floatNum != math.Floor(floatNum)
}

func checkValueType(value string) KeyType {
	if strings.HasPrefix(value, "(") && strings.HasSuffix(value, ")") {
		return Container
	} else if f, err := strconv.ParseFloat(value, 64); err == nil {
		if hasDecimal(f) {
			return Float64
		} else {
			return Int
		}
	} else if _, err := strconv.ParseBool(value); err == nil {
		return Boolean
	} else {
		return String
	}
}

func toGuessedType(value string) interface{} {
	switch checkValueType(value) {
	case Int:
		intValue, _ := strconv.Atoi(value)
		return intValue
	case Float64:
		floatValue, _ := strconv.ParseFloat(value, 64)
		return floatValue
	case Boolean:
		boolValue, _ := strconv.ParseBool(value)
		return boolValue
	case Container:
		container, _ := NewIniContainerFromString(value)
		return container
	default:
		return value
	}
}
