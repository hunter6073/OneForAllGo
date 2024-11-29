hi, welcome to the one for all project for go
all code in this project is from https://gobyexample.com/ i merely organized it into a fast food-like vrsion.
since some part of the code require termination of the program, this project is being separated into several entrances. the list of entrances are below:
1. main.go // use go run . to run main.go, this is because main.go uses other files' code
// the rest below can run using go run <pathname>
2. commandline/commandlineexample.go
3. embed/embedexample.go
4. exec/execexample.go // exec requires linux system to run, it uses ls for demonstration
5. exit/exitexample.go
6. file/fileexample.go
7. goroutine/goroutineexample.go
8. http/httpexample.go
9. recover/recoverexample.go
10. signal/signalexample.go
11. spawn/spawnexample.go // spawn requires linux system to run, it uses ls for demonstration
testing is special, you'll need to go into testing folder and run go test -v to run tests