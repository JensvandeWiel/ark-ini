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
func DeserializeIniFile(data string) (*IniFile, error) {
	lines := strings.Split(data, "\n")
	file := NewIniFile("") // Initialize an empty IniFile

	var currentSection *IniSection

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Ignore empty lines and comments
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}

		// Check if the line represents a section
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			sectionName := strings.TrimPrefix(strings.TrimSuffix(line, "]"), "[")
			currentSection = NewIniSection(sectionName)
			file.Sections = append(file.Sections, currentSection)
		} else if currentSection != nil {
			// If not a section, assume it's a key-value pair
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				currentSection.AddKey(key, value)
			}
		}
	}

	return file, nil
}