package ini

import (
	"strings"
	"testing"
)

func TestINISerialization(t *testing.T) {

	const createIniSectionTestData = `[test]
test=test
test2=test2`

	ini, err := DeserializeIniFile(createIniSectionTestData)
	if err != nil {
		t.Errorf("DeserializeIniFile() = %v", err)
	}

	if strings.TrimSpace(ini.ToString()) != createIniSectionTestData {
		t.Errorf("ini.ToString() = %v, want %v", ini.ToString(), createIniSectionTestData)
	}
}
