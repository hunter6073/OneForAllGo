package main

import (
	"embed"
)

// embed is a compiler directive that allows programs to include arbitrary files and folders in the go binary at build time

// embed directives accept paths relative to the directory containing the go source file.
// this directive embeds the contents of the file into the string variable immediately following it

//go:embed folder/single_file.txt
var fileString string

// or embed the contents of the file into a []byte
//
//go:embed folder/single_file.txt
var fileByte []byte

// we can also embed multiple files or even folders with wildcards
// this uses a variable of the embed.FS type, which implements a simple virtual file system

//go:embed folder/single_file.txt
//go:embed folder/*.hash
var folder embed.FS

func main() {
	// print out the contents of single_file.txt
	println("fileString embeds single_file.txt(hello go): ", fileString)
	println("fileByte is []byte format: ", string(fileByte))

	// retrieve some files from the embedded folder
	content1, _ := folder.ReadFile("folder/file1.hash")
	println("content1 embeds file1.hash(123): ", string(content1))

	content2, _ := folder.ReadFile("folder/file2.hash")
	println("content2 embeds file2.hash(456): ", string(content2))
}
