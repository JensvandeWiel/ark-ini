package main

import ini "github.com/JensvandeWiel/ark-ini"

func main() {
	file := ini.NewIniFile()

	// test=test1
	// test=test2
	section := file.GetOrCreateSection("test")
	section.AddKey("test", "test1")
	section.AddKey("test", "test2")
	println(file.ToString())

	file2 := ini.NewIniFile()

	// test=test2
	section2 := file2.GetOrCreateSection("test")
	section2.AddOrReplaceKey("test", "test1")
	section2.AddOrReplaceKey("test", "test2")
	println(file2.ToString())

	file3 := ini.NewIniFile("test")
	// test=test1
	// test=test2
	file3.SafelyAddKeyToSection("test", "test", "test1")
	file3.SafelyAddKeyToSection("test", "test", "test2")
	println(file3.ToString())

	file4 := ini.NewIniFile()
	// test=test2
	file4.SafelyAddKeyToSection("test", "test", "test1")
	file4.SafelyAddKeyToSection("test", "test", "test2")
	println(file4.ToString())

}
