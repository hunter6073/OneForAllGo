package main

import (
	"embed"
	"fmt"
	"net/http"
)

// embed is a compiler directive that allows programs to include arbitrary files and folders in the go binary at build time
// embed directives accept paths relative to the directory containing the go source file.
// this directive embeds the contents of the file into the string variable immediately following it

// we would use embed when we want to embed static files into the go binary,
// since the build environment might not be the same as the runtime environment

// Note: embedding large files can signifcantly increase the size of the go binary,
// also changes to the embedded files require recompilation

//go:embed folder/single_file.txt
var fileString string

// or embed the contents of the file into a []byte

//go:embed folder/single_file.txt
var fileByte []byte

// we can also embed multiple files or even folders with wildcards
// this uses a variable of the embed.FS type, which implements a simple virtual file system

//go:embed folder/single_file.txt
//go:embed folder/*.hash
var folder embed.FS

// you can also directly embed entire folders

//go:embed folder
var folder2 embed.FS

func main() {

	// print out the contents of single_file.txt
	println("fileString embeds single_file.txt(hello go): ", fileString)
	// byte array format requires string conversion to print
	println("fileByte is []byte format: ", string(fileByte))

	// folder embeds multiple files
	files, _ := folder.ReadDir("folder") // use ReadDir to get all files in the folder
	for _, file := range files {         // range over all the files, file name only returns the name
		data, _ := folder.ReadFile("folder/" + file.Name())     // use ReadFile to read a single file
		fmt.Printf("File: %s Content: %s\n", file.Name(), data) // output the data in the file
	}

	// folder 2 embeds the folder directly, the reading of the embedded file is the same
	data, _ := folder2.ReadFile(("folder/single_file.txt"))
	fmt.Println("got data from single_file.txt: ", string(data))

	// we can use http.FileServer to serve the embedded files
	// TODO: move this to http example
	http.Handle("/", http.FileServer(http.FS(folder2))) // access the file using http://localhost:8080/folder/single_file.txt
	http.ListenAndServe(":8080", nil)
}
