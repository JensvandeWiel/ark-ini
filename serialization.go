package ini

import (
	"strings"
)

// SerializeIniFile converts IniFile to INI format string
func SerializeIniFile(file *IniFile) string {
	var result string
	for _, section := range file.Sections {
		result += section.ToString()
	}
	return result
}

// DeserializeIniFile converts INI format string to IniFile
func DeserializeIniFile(data string, allowedDuplicateKeys ...string) (IniFile, error) {
	// Initialize an empty IniFile
	file := NewIniFile("", allowedDuplicateKeys...)
	var currentSection *IniSection
	var currentLine string

	// Iterate over the data character by character
	for i := 0; i < len(data); i++ {
		// If the current character is a newline character
		if data[i] == '\n' {
			// Trim leading and trailing spaces from the current line
			currentLine = strings.TrimSpace(currentLine)
			// If the line is empty or a comment, skip it
			if currentLine == "" || strings.HasPrefix(currentLine, ";") || strings.HasPrefix(currentLine, "#") {
				currentLine = ""
				continue
			}

			// If the line is a section
			if strings.HasPrefix(currentLine, "[") && strings.HasSuffix(currentLine, "]") {
				// Extract the section name and create a new section
				sectionName := strings.TrimPrefix(strings.TrimSuffix(currentLine, "]"), "[")
				currentSection = NewIniSection(sectionName, &file.AllowedDuplicateKeys)
				// Add the new section to the IniFile
				file.Sections = append(file.Sections, currentSection)
			} else if currentSection != nil {
				// If the line is a key-value pair
				if strings.Contains(currentLine, "=") {
					// Split the line into key and value
					keyValuePair := strings.SplitN(currentLine, "=", 2)
					key := keyValuePair[0]
					value := keyValuePair[1]

					// If the value is a container
					if strings.HasPrefix(value, "(") {
						// Keep appending lines to the value until the closing parenthesis is found
						for !strings.HasSuffix(value, ")") && i < len(data) {
							i++
							value += string(data[i])
						}
					}

					// Add the key-value pair to the current section
					currentSection.AddKey(key, toGuessedType(value))
				}
			}
			// Reset the current line
			currentLine = ""
		} else {
			// If the current character is not a newline character, add it to the current line
			currentLine += string(data[i])
		}
	}

	// If there's a remaining line after the loop, handle it
	if currentLine != "" {
		currentLine = strings.TrimSpace(currentLine)
		// If the line is not a comment and there's a current section
		if !(strings.HasPrefix(currentLine, ";") || strings.HasPrefix(currentLine, "#")) && currentSection != nil {
			// If the line is a key-value pair
			if strings.Contains(currentLine, "=") {
				// Add the key-value pair to the current section
				currentSection.AddParsedKey(currentLine)
			}
		}
	}

	// Return the IniFile
	return file, nil
}
