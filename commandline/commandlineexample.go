package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	// use "./command-line-arguments a b c d" to intake arguments ,note that it is best to build first
	// os.Args is a slice of the command-line arguments, starting with the program name.
	argsWithProg := os.Args
	// os.Args[1:] holds the argument to the program(lose the command name)
	argsWithoutProg := os.Args[1:]
	// get individual argument with normal indexing
	arg := os.Args[1]
	// output these arguments
	fmt.Println("arguments with program: ", argsWithProg)
	fmt.Println("arguments without program: ", argsWithoutProg)
	fmt.Println("first argument: ", arg)

	// format of a flag: -word=opt
	// basic flag declarations are available for string ,integer and boolean options.
	// here we declare a string flag word with a default value "foo" and a short description.
	// this flag.String function returns a string pointer (not a string value);
	wordPtr := flag.String("word", "foo", "a string")
	// then we declare an int flag
	numbPtr := flag.Int("numb", 42, "an int")
	forkPtr := flag.Bool("fork", false, "a bool")
	// note that it is possible to declare an option that uses an existing var declared elsewhere in the program. note that we need to pass in a pointer to the flag declaration function.
	var svar string
	flag.StringVar(&svar, "svar", "bar", "a string var")
	flag.Parse() // once all the flags are declared, call flag.parse to execute the command-line parsing
	// Here we’ll just dump out the parsed options and any trailing positional arguments. Note that we need to dereference the pointers with e.g. *wordPtr to get the actual option values.
	fmt.Println("word:", *wordPtr)
	fmt.Println("numb:", *numbPtr)
	fmt.Println("fork:", *forkPtr)
	fmt.Println("svar:", svar)
	fmt.Println("tail:", flag.Args())
	// do this to get the flags: ./command-line-flags -word=opt -numb=7 -fork -svar=flag
	// if the commandline flags are nil, they automatically take their default values
	// note that the flag package requires all flags to appear before positional arguments(or they will be interpreteed as positional arguments)
	// -h will automatically stop the program and print the help text
	// if you provide a flag that wasn't specified to the falg package, the program will print an error message and show the help text again

	// we declare a subcommand using the NewFlagSet function, and proceed to declare flags specific to this subcommand.
	fooCmd := flag.NewFlagSet("foo", flag.ExitOnError)
	fooEnable := fooCmd.Bool("enable", false, "enable")
	fooName := fooCmd.String("name", "", "name")

	// For a different subcommand we can define different supported flags.
	barCmd := flag.NewFlagSet("bar", flag.ExitOnError)
	barLevel := barCmd.Int("level", 0, "level")
	//The subcommand is expected as the first argument to the program.
	if len(os.Args) < 2 {
		fmt.Println("expected 'foo' or 'bar' subcommands")
		os.Exit(1)
	}

	// check which subcommand is invoked
	switch os.Args[1] {
	// for every subcommand, we parse its own flags and have access to trailing positional arguments
	case "foo":
		fooCmd.Parse(os.Args[2:])
		fmt.Println("subcommand 'foo'")
		fmt.Println("  enable:", *fooEnable)
		fmt.Println("  name:", *fooName)
		fmt.Println("  tail:", fooCmd.Args())
	case "bar":
		barCmd.Parse(os.Args[2:])
		fmt.Println("subcommand 'bar'")
		fmt.Println("  level:", *barLevel)
		fmt.Println("  tail:", barCmd.Args())
	default:
		fmt.Println("expected 'foo' or 'bar' subcommands")
		os.Exit(1)
	}

	// environment variables
	// use os.setenv to set a key/value pair.
	os.Setenv("FOO", "1")
	// use os.getenv to get a value of a key
	fmt.Println("FOO:", os.Getenv("FOO"))
	fmt.Println("BAR:", os.Getenv("BAR"))
	// Use os.Environ to list all key/value pairs in the environment. This returns a slice of strings in the form KEY=value.
	// You can strings.SplitN them to get the key and value. Here we print all the keys.
	fmt.Println("os.Environ:")
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		fmt.Println(pair[0])
	}
	// BAR=2 go run environment-variables.go , set the bar in the environment first
}
