package main

import (
	"fmt"
	"io"
	"os/exec"
)

func main() {

	// the exec.Command helper creates a command object, the command does not run here
	dateCmd := exec.Command("date")
	// the output method runs the command, if there were no errors, dateOut will hold bytes of the result
	dateOut, err := dateCmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println("> date")
	// print the result of dateOut
	fmt.Println(string(dateOut))

	// Output and other methods of Command will return *exec.Error if
	// there was a problem executing the command and *exec.ExitError
	// if the command ran but exited with a non-zero return code
	_, err = exec.Command("date", "-x").Output() // will return rc = 1
	if err != nil {
		switch e := err.(type) {
		case *exec.Error:
			fmt.Println("failed executing:", err)
		case *exec.ExitError:
			fmt.Println("command exit rc =", e.ExitCode())
		default:
			panic(err)
		}
	}

	// grep command with pipes
	grepCmd := exec.Command("grep", "hello")
	// grab input and output pipes
	grepIn, _ := grepCmd.StdinPipe()
	grepOut, _ := grepCmd.StdoutPipe()
	// start the process
	grepCmd.Start()
	// write input to the inpipe
	grepIn.Write([]byte("hello grep 123\ngoodbye grep"))
	// close the inpipe
	grepIn.Close()
	// read from the outpipe
	grepBytes, _ := io.ReadAll(grepOut)
	// wait for the process to exit
	grepCmd.Wait()

	fmt.Println("> grep hello")
	// output the bytes read from the outpipe
	fmt.Println(string(grepBytes))

	// when spawning commands we need to provide a delineated command and argument array.
	// if you want to spawn a full command with a string, you can use bash's -c option
	lsCmd := exec.Command("bash", "-c", "ls -a -l -h") // one line command, not delineated
	lsOut, err := lsCmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println("> ls -a -l -h")
	fmt.Println(string(lsOut))
}
