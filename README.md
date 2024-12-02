hi, welcome to the one for all project for go
all code in this project is from https://gobyexample.com/ i merely organized it into a fast food-like vrsion.
since some part of the code require termination of the program, this project is being separated into several entrances. the list of entrances and the command to run them are:
1. main.go 
    go run .  // this is because main.go uses other files' code, you can also use "go run a.go b.go"
2. embed/embedexample.go
    go run embed/embedexample.go
3. exit/exitexample.go
    go run exit/exitexample.go
4. file/fileexample.go
    echo abcbcd| go run file/fileexample.go // there is piping in the example
5. http/httpexample.go
    go run http/httpexample.go
6. recover/recoverexample.go
    go run recover/recoverexample.go
7. signal/signalexample.go
    go run signal/signalexample.go
8. goroutine/goroutineexample.go
    go run goroutine/goroutineexample.go
9. commandline/commandlineexample.go
    go run commandline/commandlineexample.go -word=opt -numb=7 -fork -svar=flag -bar=1 // TODO: not the final version
// exec and spawn requires linux system to run, it uses ls for demonstration
10. exec/execexample.go
    go run exec/execexample.go 
11. spawn/spawnexample.go
    go run spawn/spawnexample.go
// testing is special, you'll need to go into testing folder then run the following
    go test -v  // run all tests in the current project in verbose mode
    go test -bench=.  // run all the benchmark tests in the current project. all tests are run prior to benchmarks