package ini

import (
	"strings"
	"testing"
)

func TestCreateINISection(t *testing.T) {

	const createIniSectionTestData = `[test]
test=test
test2=test2`

	section := NewIniSection("test")
	section.AddKey("test", "test")
	section.AddParsedKey("test2=test2")

	//check is the section is the same as the test data using section.ToString()
	if strings.TrimSpace(section.ToString()) != createIniSectionTestData {
		t.Errorf("section.ToString() = %v, want %v", section.ToString(), createIniSectionTestData)
	}

}
