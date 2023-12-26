package ini

import (
	"strings"
	"testing"
)

func TestNewIniFile(t *testing.T) {
	ini := NewIniFile("")
	if ini == nil {
		t.Error("NewIniFile() failed.")
	}

	ini.AddKeyToSection("default", "key", "value")
	key, err := ini.GetKeyFromSection("default", "key")
	if err != nil {
		return
	}

	if key.ToValueString() != "value" {
		t.Fail()
	}
}

func TestSerializeAndDeserializeIniFile(t *testing.T) {
	data := `[default]
key=value
key2=(key=(bob=bab),wow=22.1)`

	ini, _ := DeserializeIniFile(data)
	if ini == nil {
		t.Error("NewIniFile() failed.")
	}

	if strings.TrimSpace(ini.ToString()) != data {
		t.Error(`ini.ToString() is not the same as data
data:
` + data + `
ini.ToString():
` + ini.ToString() + `
`)
	}
}

func TestIniFile_SafelyAddKeyToSection(t *testing.T) {
	ini1 := NewIniFile()
	ini1.SafelyAddKeyToSection("default", "key", "value")
	ini1.SafelyAddKeyToSection("default", "key", "value2")
	keys1, err := ini1.GetKeyFromSectionWithMultipleValues("default", "key")
	if err != nil {
		t.Error(err)
		return
	}

	if len(keys1) != 1 {
		t.Fail()
	}

	if keys1[0].ToValueString() != "value2" {
		t.Fail()
	}

	ini2 := NewIniFile("key")
	ini2.SafelyAddKeyToSection("default", "key", "value")
	ini2.SafelyAddKeyToSection("default", "key", "value2")
	keys2, err := ini2.GetKeyFromSectionWithMultipleValues("default", "key")
	if err != nil {
		t.Error(err)
		return
	}

	if len(keys2) != 2 {
		t.Fail()
	}

	if keys2[0].ToValueString() != "value" {
		t.Fail()
	}

	if keys2[1].ToValueString() != "value2" {
		t.Fail()
	}
}

func TestToMap(t *testing.T) {
	data := `[default]
bob=value
bob=value2
key2=(key=(bob=bab),wow=22.1)`

	ini, _ := DeserializeIniFile(data)
	mapp := ToMap(ini)
	if strings.TrimSpace(DeserializeFromMap(mapp).ToString()) != strings.TrimSpace(data) {
		t.Fail()
	}
}
