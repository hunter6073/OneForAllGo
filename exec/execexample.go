package main

import (
	"os"
	"os/exec"
	"syscall"
)

func main() {

	// exec completely replace the current go process with another one.

	// go requires an absolute path to the binary we want to execute,
	// so we'll use exec.Lookpath to find it
	binary, lookErr := exec.LookPath("ls")
	if lookErr != nil {
		panic(lookErr)
	}
	// exec requires arguments in slice form, note that the first argument should be the program name
	args := []string{"ls", "-a", "-l", "-h"}

	//exec also needs a set of environment variables, here we just provide our current environment
	env := os.Environ()

	// here's the actual syscall.Exec. if this call is successful, the execution of
	// our process will end here and be replaced by the /bin/ls -a -l -h process
	execErr := syscall.Exec(binary, args, env)
	// if there is an error we'll get a return value
	if execErr != nil {
		panic(execErr)
	}
}
