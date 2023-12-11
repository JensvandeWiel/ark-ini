package ini

import (
	"strings"
)

//TODO add parsing of values and save them as the type and not a string

// SerializeIniFile converts IniFile to INI format string
func SerializeIniFile(file *IniFile) string {
	var result string
	for _, section := range file.Sections {
		result += section.ToString()
	}
	return result
}

//TODO add parsing of values and save them as the type and not a string

// DeserializeIniFile converts INI format string to IniFile
func DeserializeIniFile(data string) (IniFile, error) {
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
			currentSection.AddParsedKey(line)
		}
	}

	return file, nil
}
