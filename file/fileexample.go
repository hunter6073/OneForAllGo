package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// reading files requires checking most calls for errors.This helper will streamline our error checks below.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func visit(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	fmt.Println(" ", path, d.IsDir())
	return nil
}

// run this program using "echo abcbcd| go run file/fileexample.go"
// TODO: fix this code so that the subdir folder is generated and deleted
func main() {
	/******************************* file paths **************************************/

	// use filepath.join to construct paths
	path := filepath.Join("dir1", "dir2", "filename") // construct a path in hierarchial order
	fmt.Println("join dir1, dir2 and filename: ", path)
	// you should always use Join instead of concatenating /s or \s manually
	// the following two joins both return "dir1/filename"
	fmt.Println(filepath.Join("dir1//", "filename"))
	fmt.Println(filepath.Join("dir1/../dir1", "filename"))
	// dir and base ca nbe used to split a path to the directory and the file.alternatively. split will return both in the same call
	fmt.Println("Dir(p):", filepath.Dir(path))   // dir gets the directory
	fmt.Println("Base(p):", filepath.Base(path)) // base gets the final file
	// check if a path is absolute, both returns false
	fmt.Println("dir/file is absolute: ", filepath.IsAbs("dir/file"))       // just false
	fmt.Println("/dir/file is absolute: ", filepath.IsAbs("/dir/file"))     // true in linux
	fmt.Println("D:/dir/file is absolute: ", filepath.IsAbs("D:/dir/file")) // true in windows

	// extensions
	filename := "config.json"
	ext := filepath.Ext(filename) // get the extension of a path
	fmt.Println("extension of config.json is: ", ext)
	// to find the file's name with the extension removed, use strings.TrimSuffix
	fmt.Println("file name with extension removed: ", strings.TrimSuffix(filename, ext))

	// rel finds a relative path between a base and a target, it returns an error if the target cannot be made relative to the base
	rel, err := filepath.Rel("a/b", "a/b/t/file")
	if err != nil {
		panic(err)
	}
	fmt.Println("relative path between a/b and a/b/t/file: ", rel)

	rel, err = filepath.Rel("a/b", "a/c/t/file")
	if err != nil {
		panic(err)
	}
	fmt.Println("relative path between a/b and a/c/t/file: ", rel)

	/******************************* directories **************************************/

	// create a new sub-directory in the current working directory
	errmkdir := os.Mkdir("subdir", 0755) // 0755 is the permission of the directory
	check(errmkdir)
	defer os.RemoveAll("subdir") // remove the directory at the end of the program

	// helper function to create a new empty file
	createEmptyFile := func(name string) {
		d := []byte("")
		check(os.WriteFile(name, d, 0644))
	}
	createEmptyFile("subdir/file1")

	// we can create a hierarchy of directories, including parents with mkdirall
	err = os.MkdirAll("subdir/parent/child", 0755)
	check(err)

	// create 3 file in the temp directory
	createEmptyFile("subdir/parent/file2")
	createEmptyFile("subdir/parent/file3")
	createEmptyFile("subdir/parent/child/file4")

	// readdir lists directory contents, returning a slice of os.direntry objects
	cmkdir, err := os.ReadDir("subdir/parent")
	check(err)

	fmt.Println("Listing subdir/parent")
	for _, entry := range cmkdir {
		fmt.Println(" ", entry.Name(), entry.IsDir())
	}

	// chdir lets us change the current working directory, similar to cd
	err = os.Chdir("subdir/parent/child")
	check(err)

	// // now we'll see the contents of subdir/parent/child when listing the current directory
	c, err := os.ReadDir(".")
	check(err)

	fmt.Println("Listing subdir/parent/child")
	for _, entry := range c {
		fmt.Println(" ", entry.Name(), entry.IsDir())
	}

	err = os.Chdir("../../..") // return to original location ,allowing the successful deletion of subdir
	check(err)

	//We can also visit a directory recursively, including all its sub-directories. WalkDir accepts a callback function to handle every file or directory visited.
	fmt.Println("Visiting subdir")
	err = filepath.WalkDir("subdir", visit)

	/********************************** temporary files ********************************/
	// creating temporary files
	f, err := os.CreateTemp("", "data.txt")
	check(err)
	fmt.Println("Temp file name:", f.Name())
	defer os.Remove(f.Name()) // always nice to remove the temporary file using defer

	_, err = f.Write([]byte{1, 2, 3, 4}) // write to file
	check(err)

	dname, err := os.MkdirTemp("", "sampledir") // create temporary directory
	check(err)
	fmt.Println("Temp dir name:", dname)

	defer os.RemoveAll(dname) // remove the temporary directory

	fname := filepath.Join(dname, "file1")        // join name to form directory
	err = os.WriteFile(fname, []byte{1, 2}, 0666) // write to file
	check(err)

	/************************* file write **********************************************/
	filePath := "data.txt"
	// write s string to file directly
	d1 := []byte("hello world\ngo\n")        // notice that d1 must be a byte array and not a string
	errw := os.WriteFile(filePath, d1, 0644) // 0644 if the permission of the file
	check(errw)

	// for more granular writes, open a file for writing
	filewrite, errw := os.Create("tempfile") // this creates a temp file
	check(errw)
	// it's idiomatic to defer a close immediately after opening a file
	defer filewrite.Close()

	// d2 is an array of bytes
	d2 := []byte{115, 111, 109, 101, 10} // translates to "some"
	fn2, err := filewrite.Write(d2)      // write bytes to file
	check(err)
	fmt.Printf("wrote %d bytes to tempfile\n", fn2)

	// write a string to the file
	fn3, err := filewrite.WriteString("writes\n")
	check(err)
	fmt.Printf("wrote %d bytes\n", fn3)
	// issue a sync to flush writes to stable storage
	filewrite.Sync()

	// bufio provides buffered writers too
	w := bufio.NewWriter(filewrite)
	n4, err := w.WriteString("buffered\n")
	check(err)
	fmt.Printf("wrote %d bytes\n", n4)
	// use flush to ensure all buffered operatiosn have been applied to the underlying writer
	w.Flush()

	/*********************** file read **********************************************/

	// read afile's entire contents into memory
	dataf, errf := os.ReadFile(filePath)
	check(errf)
	fmt.Println("./data.txt contains: ", string(dataf))

	// more often, you'll want to control over how and what parts of a file are read.
	// first ,start by opening a file
	filef, err := os.Open(filePath)
	check(err)
	defer filef.Close() // use defer to make sure to close the file upon program end.

	b1 := make([]byte, 14)    // create a slice holding 14 bytes
	n1, err := filef.Read(b1) // read 14 bytes into the file, n1 is the number of the bytes
	check(err)                // check for error
	fmt.Printf("%d bytes: %s\n", n1, string(b1[:n1]))

	// you can also seek to a know location in the file and read from there
	o2, err := filef.Seek(6, io.SeekStart)
	check(err)

	b2 := make([]byte, 2)
	n2, err := filef.Read(b2)
	check(err)
	fmt.Printf("%d bytes @ %d: ", n2, o2)
	fmt.Printf("%v\n", string(b2[:n2]))

	// other methods of seeking are relative to the current cursor position
	_, err = filef.Seek(4, io.SeekCurrent) // set file current position to 4 bytes
	check(err)

	// and relative tot the end of the file
	_, err = filef.Seek(-3, io.SeekEnd) // set file current position to end -3 bytes
	check(err)

	// rewind the filef pointer to the beginning
	o3, err := filef.Seek(0, io.SeekStart)
	check(err)

	// reads like the ones above can be more robustly implemented with ReadAtLeast
	b3 := make([]byte, 2)
	n3, err := io.ReadAtLeast(filef, b3, 2)
	check(err)
	fmt.Printf("%d bytes @ %d: %s\n", n3, o3, string(b3))

	// the bufio package implements a buffered reader that may be useful both for its efficiency
	// with many small reads and because of the additional reading methods it provides
	r4 := bufio.NewReader(filef)
	b4, err := r4.Peek(5) // peek returns the next n bytes without advancing the reader
	check(err)
	fmt.Printf("5 bytes: %s\n", string(b4))

	/****************************** line filters *********************************/

	// wrappping the unbuffered os.stdin with a buffered scanner gives us a convenient scan method that
	//advances the scanner to the next token; which is the next line in the default scanner
	scanner := bufio.NewScanner(os.Stdin) // scans stdin, if you didn't give any from the running of the program, it will have you input directly.
	//Text returns the current token, here the next line, from the input.
	for scanner.Scan() {
		ucl := strings.ToUpper(scanner.Text()) // scanner.Text() returns each line
		//Write out the uppercased line.
		fmt.Println(ucl) // turns each line into upper case variants
	}
	//Check for errors during Scan. End of file is expected and not reported by Scan as an error.
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

}
