package main

import "github.com/JensvandeWiel/ark-ini"

func main() {
	iniFile := ini.NewIniFile("C://bob/bob.ini", "allowedDuplicateValue1", "allowedDuplicateValue2")
}
