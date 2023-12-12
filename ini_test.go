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
	key, err := ini.FindKeyFromSection("default", "key")
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
