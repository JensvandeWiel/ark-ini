package main

import (
	"github.com/JensvandeWiel/ark-ini"
	"strconv"
)

func main() {
	str := `[bob]
test=1
test2=1
test2=2`
	iniFile, _ := ini.DeserializeIniFile(str)
	sec, exists := iniFile.GetSection("bob")
	if !exists {
		panic("section doesn't exist")
	}

	println(iniFile.ToString())
	println(sec.CheckForMultipleKeys("test2"))
	keys := sec.GetMultipleKeys("test2")
	for _, key := range keys {
		intkey, err := key.AsInt()
		if err != nil {
			panic(err)
		}
		key.Value = strconv.Itoa(intkey + 1)
	}
	println(iniFile.ToString())

	sec.AddOrReplaceKey("test3", "(Bob=1,Bob2=2)")

	println(iniFile.ToString())

	ke, exists := sec.GetKey("test3")
	if !exists {
		panic("key doesn't exist")
	}

	val, typ := ke.AsGuessedValue()
	if typ != ini.Container {
		panic(typ)
	}
	wow := val.(ini.IniContainer)

	ketro, found := wow.FindKey("Bob2")
	if !found {
		panic("key doesn't exist")
	}

	valfj, _ := ketro.AsInt()

	println(valfj)

}
