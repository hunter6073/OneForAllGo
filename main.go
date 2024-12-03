// there are several files containing go code, they are all under the main package,
// thus there is no need to explicitly declare function usage
package main

// import other packages to use
import (
	// in order to use a go file in another folder, must set the other file to another package, but under the same module
	"RobotTask/helper"
	"cmp"
	"errors"
	"fmt"
	"maps"
	"math"
	"slices"
	"unicode/utf8"
)

const constantS string = "constant" // declaring a constant value

func plus(a int, b int) int { // regular function declaring seperate variables
	return a + b
}

func plus2(a, b, c int) (int, int) { // if all the variables have the same type, you can use a single type declaration
	return a + b, c // notice that this function can return multiple values
}

func sum(nums ...int) { // this function can take in (1,2,3) or []int{1,2,3,4} as parameters, both are then transformed into a slice called nums
	fmt.Println(nums) // print the array of ints
	total := 0
	for _, num := range nums { // use range to iterate over the slice
		total += num
	}
	// both for loops are valid ways to iterate over a slice
	for i := 0; i < len(nums); i++ {
		total += nums[i]
	}
	fmt.Println(total)
}

func intSeq() func() int { // the intset function returns a function that returns an int
	i := 0
	return func() int {
		i++
		return i
	}
}

func fact(n int) int { // a recursive function
	if n == 0 {
		return 1
	}
	return n * fact(n-1)
}

type rect struct { // defining a struct
	width, height int    // multiple variables with the same type
	name          string // different types
}

// function that generates a new rect struct every time
// go is a garbage collected language, so you don't have to worry about memory management
// when there are no active referenced to it, the pointer will be automatically cleaned by by the garbage collector
func newRect(name string) *rect {
	r := rect{name: name}
	r.width = 10
	r.height = 20
	return &r
}

// you cannot designate a method in the struct,but is available to do this outside, *rect is a receiver type
func (r *rect) area() int { // a method of the rect struct, notice that the method designated a struct before the function name
	return r.width * r.height
}

// go automatically handles conversions between values and pointers for method calls
// but you may want to use a pointer receiver type to avoid copying on method calls or to allow the method to mutate the receiving struct
func (r rect) perim() int { // Methods can be defined for either pointer or value receiver types
	return 2 * (r.width + r.height)
}

type geometry interface { //If a variable has an interface type, then we can call methods that are in the named interface
	area() float64
	perim() float64
}

type circle struct { // defining a struct
	radius float64
}

// implementing the methods that are in the named interface
func (c circle) area() float64 { // a method of the circle struct
	return math.Pi * c.radius * c.radius
}
func (c circle) perim() float64 {
	return 2 * math.Pi * c.radius
}
func measure(g geometry) { // a function that takes in a geometry interface
	fmt.Println("printing geometry: ", g)
	fmt.Println("area of geometry is: ", g.area())
	fmt.Println("perimeter of geometry is: ", g.perim())
}

type ServerState int // defining an int type

const (
	StateIdle ServerState = iota // this special keyword generates successive constant values automatically,
	StateConnected
	StateError
	StateRetrying
)

var stateName = map[ServerState]string{
	StateIdle:      "idle",
	StateConnected: "connected",
	StateError:     "error",
	StateRetrying:  "retrying",
}

// overwritign the string method of the ServerState type, this works for Println
func (ss ServerState) String() string {
	return stateName[ss]
}

func transition(s ServerState) ServerState {
	switch s {
	case StateIdle:
		return StateConnected
	case StateConnected, StateRetrying:
		return StateIdle
	case StateError:
		return StateError
	default:
		panic(fmt.Errorf("unknown state: %s", s))
	}
}

// regular struct
type base struct {
	num int
}

// function for a struct
func (b base) describe() string {
	// sprintf prints and returns a string, %+v will output the fields in the struct
	return fmt.Sprintf("base with num=%+v", b.num)
}

// struct that embeds a struct
type container struct {
	base // an embedded struct, looks like a field without a name
	str  string
}

// by convention, errors are the last return value and have type error, a built-in interface
func fun(arg int) (int, error) {
	if arg == 42 {
		// errors.New constructs a basic error value with the given error message
		return -1, errors.New("can't work with 42")
	}
	// a nil value in the error position indicates that there was no error
	return arg + 3, nil
}

// a sentinel error is a predeclared variable that is used to signify a specific error condition
var ErrOutOfTea = fmt.Errorf("no more tea available")
var ErrPower = fmt.Errorf("can't boil water")

func makeTea(arg int) error {
	if arg == 2 {
		return ErrOutOfTea // using a predefined error
	} else if arg == 4 {
		// we can wrap errors with higher-level errors to add context
		// wrapped errors create a logical chain that can be queried with functions like errors.Is and errors.As
		return fmt.Errorf("making tea: %w", ErrPower)
	}
	return nil
}

// creating a custom error type
type argError struct {
	arg     int
	message string
}

// adding this Error method makes argError implement the error interface
func (e *argError) Error() string {
	return fmt.Sprintf("%d - %s", e.arg, e.message)
}

// return custom error
func fun2(arg int) (int, error) {
	if arg == 42 {
		return -1, &argError{arg, "can't work with it"}
	}
	return arg + 3, nil
}

// like c or c++, go require a main entrance to run the program
func main() {

	/********************************** variables ************************************************/

	// using var to define a variable, assigning the string type is optional, go will automatically infer the type
	var variable1 = "test" // in this case, the type of variable1 is inferred as string

	// using %T to print the type of a variable, you can also use the reflect.TypeOf() to get the type of a variable
	fmt.Printf("variabl11 is:%s, type of variable1 is: %T\n", variable1, variable1)

	// using a[0] returns the first byte of the string, if you want to print the char,
	// you need to convert it to string, also, use + to concatenate strings
	fmt.Println("getting the first char of variable1: " + string(variable1[0]))

	// defining multiple variables and setting their type
	var variable2, variable3 int = 1, 2

	// defining a boolean variable
	var variable4 = true

	// variables declared without initialization are zero-valued. for int, it's 0
	var variable5 int

	// short hand for var f string = "apple", this syntax is only available inside functions
	variable6 := "apple"

	// all declared variables must be used or it won't compile， notice you can use comma to separately print multiple variables
	fmt.Println("variable1:", variable1, ",variable2:", variable2, ",variable3:", variable3, ",variable4:", variable4, ",variable5:", variable5, ",variable6:", variable6)

	// a const statement can appear anywhere a var statement can
	const constnum = 500000000
	// Constant expressions perform arithmetic with arbitrary precision.
	const constd = 3e20 / constnum
	fmt.Println("constd is: ", constd)
	// A numeric constant has no type until it’s given one, such as by an explicit conversion.
	fmt.Println("int64(d) is: ", int64(constd)) // constd is now type int64
	// A number can be given a type by using it in a context that requires one,
	// such as a variable assignment or function call. For example, here math.Sin expects a float64.
	fmt.Println("math.Sin(constnum) is: ", math.Sin(constnum))

	/************************************* loops,ifs and switches ************************************************/
	// in go, there is no while keyword, for can be used as both, for example:
	for variable2 < 3 { // this is the same as while(b<3)
		variable2++
	}

	// example of a classic for loop
	for i := 0; i < 3; i++ {
		fmt.Println("for loop i: ", i)
		// demonstrating basic if else, there is no brackets for if in go
		//also, there is no a=b>c?b:c in go, must use full if statement
		if i%2 == 0 {
			continue // continue to next iteration
		} else {
			break // use break to break out of a loop
		}
	}

	for i := range 3 { // using the range keyword to iterate 3 times
		fmt.Println("using range: ", i) // i is from 0 to 2
	}

	// note that a for loop without a condition will loop reapeatedly until you break or return
	for {
		break
	}
	// go allows preceding statements for if conditions, like the followoing:
	if num := 9; num < 0 {
		fmt.Println("this will never happen")
	}

	// example of a classic switch
	cas := 0
	switch cas { // classic switch case
	case 2: // classic usage
		fmt.Println("cas is 2")
	case -1, 0: // case allows using comma to separate multiple expressions
		fmt.Println("cas+1 is: ", cas+1)
	default: // default case
		fmt.Println("cas+1 is: ", cas+3)
	}

	//switch without an expression is an alternate way to express if/else logic
	cas2 := 3
	switch { // catches in order, if the first catch succeeds, will not proceed to the next catches
	case cas2 < 4: // catch by cas2<4 first
		fmt.Println("cas2<4")
	case cas2 > 2:
		fmt.Println("cas2>2")
	default:
		fmt.Println("default cas2")
	}

	// whatAmI is a function that takes an interface as a parameter, an interface type means i can be anything
	whatAmI := func(i interface{}) {
		switch t := i.(type) { // switch the type of i, note that t is not explicitly required.
		case bool: // if the type is bool
			fmt.Println("I'm a bool")
		case int:
			fmt.Println("I'm an int")
		default:
			fmt.Printf("Don't know type %T\n", t)
		}
	}
	whatAmI(3) // running the function

	/************************************ arrays(not that commonly used) *******************************/

	// in typical go code, slices are much more common, the traditional arrays are userful in certain special scenarios
	// so apparently arrays have predefined lengths, if the length is not defined, it's a slice
	var arr [5]int        // declaring an array of 5 integers
	arr[4] = 100          // setting the 5th element of the array to 100
	fmt.Println(len(arr)) // use len() to get the length of an array

	// go's arrays are values, an array variable denotes the entire array. note that it is not a pointer to the first array element like in c.
	// this means that when you assign or pass around array value, you will make a copy of its contents.(to avoid this, use pointers to the array)
	var arr2 = [5]int{1, 2, 3, 4, 5} // declaring an array with values
	fmt.Println(arr2)

	var brr = [...]int{1, 2, 3, 4, 5, 10: 100} // using ... to automaticaly count the length of the array, while setting the initial values
	fmt.Println(brr)                           // since we used 10:100, the remaining elements will be set to 0, and the 10th element will be set to 100

	var twoD = [2][3]int{ // declaring a 2d array with values
		{1, 2, 3},
		{4, 5, 6},
	}
	// using for loop to traverse the array
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			twoD[i][j]++
		}
	}
	fmt.Println("2d: ", twoD)

	/*********************************** slices *************************************************/

	// initialize an empty slice, slices do not require element count
	var slice1 []string
	// len returns the length of the slice， cap returns the capacity of the slice
	fmt.Println("uninit slice1:", slice1, "slice1==nil is: ", slice1 == nil, "len(slice1) is: ", len(slice1), "cap(slice1)  is: ", cap(slice1))

	// using make to create an empty slice with non-zero length, other items are set to zero values, for the case of string ,they are set to ""
	slice1 = make([]string, 3) // since capacity is not explicitly declared here, the len and cap of slice1 are both 3
	slice1[0] = "a"            // setting an item's value in the slice
	fmt.Println("made slice1:", slice1, "slice1==nil is: ", slice1 == nil, "len(slice1) is: ", len(slice1), "cap(slice1) is: ", cap(slice1))

	// append to the slice, the append operation is for slices but not arrays, this example appends two items to the slice
	slice1 = append(slice1, "e") // you can also append multiple items, such as append(slice1,"e","f")
	// now slice1 has length 4 and capacity 6, this is because append slice1 has already reached its capacity, so append will create a new copy and the new slice size and capacity is now doubled
	fmt.Println("slice1 after appending: ", slice1, "len is: ", len(slice1), "cap is: ", cap(slice1))
	slice1 = append(slice1, "f")
	// slice1's capacity is still 6, because we haven't reached the capacity, so nothing happens.
	fmt.Println("slice1 after appending: ", slice1, "len is: ", len(slice1), "cap is: ", cap(slice1))

	// here we actually grew the size of the slice, note that the length of slice1 was originally 5, but now it's 6
	// the point here is you can shrink or grow the size of the slice anywhere from 0 to the capacity, but not out of this bound
	// if you do slice1 = sliece1[:7], this will render a panic, as you've grown the size ver the capacity
	slice1 = slice1[:6]
	fmt.Println("slice1 after slicing: ", slice1, "len is: ", len(slice1), "cap is: ", cap(slice1))

	// copy a slice
	copiedSlice := make([]string, len(slice1)) // use len to get the length of the slice, use cap to get the capacity. both returns 0 for a nil slice
	copy(copiedSlice, slice1)                  // copy into cs from s, operations to cs will not affect s, copy will return the number of elements copied
	fmt.Println("copiedSlice is an exact copy of slice1: ", copiedSlice)

	// slicing does not copy the slice's data, it creates a new slice value that points to the original array. modifyin the elements of a reslice will modify the original
	sub := copiedSlice[2:5] // getting a subslice, use [2:] to begin at index 2, and [:5] to end at index 5(excluding 5)
	fmt.Println("subslice 2:5 of copiedSlice: ", sub)
	sub[0] = "z"
	fmt.Println("when you change sub, copiedSlice is also changed: ", sub, copiedSlice)

	slice2 := []string{"g", "h", "i"}  // initializing a slice with values
	fmt.Println("slice2 is: ", slice2) // use slices.Equal() to compare slices

	if slices.Equal(copiedSlice, slice2) { // use slices.Equal to check if two slices are equal
		fmt.Println("copiedSlice and slice2 are equal")
	} else {
		fmt.Println("copiedSlice and slice2 are not equal")
	}

	// this is a dynamic 2d slice
	var twoDSlice [][]int
	for i := 0; i < 3; i++ {
		var oneDSlice []int
		for j := 0; j < i+1; j++ {
			oneDSlice = append(oneDSlice, j)
		}
		twoDSlice = append(twoDSlice, oneDSlice)
	}
	fmt.Println("dynamic twoDSlice: ", twoDSlice)

	// slice pop
	// var stack []int; // this will initialize an empty slice as well
	stack := []int{} // initialize a stack
	for i := range 3 {
		stack = append(stack, i) // push to the stack
	}
	fmt.Println("original stack is: ", stack)

	// pop item from stack(last in first out)
	popItem := stack[len(stack)-1]
	stack = stack[:len(stack)-1]
	fmt.Println("popped item is:", popItem, "stack is: ", stack)
	// pop item from queue(first in first out)
	queueItem := stack[0]
	stack = stack[1:]
	fmt.Println("queue item is:", queueItem, "queue is: ", stack)

	/*********************************** maps *************************************************/

	//mp := make(map[string]int) // make a map, the key must be a string, and the value must be an int
	mp := map[string]int{} // this is the same as the above
	mp["k1"] = 7           // set a value to a key
	delete(mp, "k1")       // delete one key in the map
	clear(mp)              // delete all keys in the map

	_, exists := mp["k2"] // the optional second value is the boolean value of whether the key exists
	fmt.Println("does k2 exists in mp? ", exists)

	mpn := map[string]int{"foo": 1, "bar": 2} // initializing a map with values, we can use maps.Equal() to compare maps
	fmt.Println("mpn is: ", mpn)
	mpn2 := map[string]int{"foo": 1, "bar": 2}
	if maps.Equal(mpn, mpn2) { // the maps package contains a number of useful utility functions for maps
		fmt.Println("mpn == mpn2")
	}

	/*********************************** functions *************************************************/

	// regular functions
	fmt.Println("plus 1 and 2 and you'll get: ", plus(1, 2)) // use the function
	plus2_1, plus2_2 := plus2(1, 2, 3)                       // function returning multiple variables
	fmt.Println("trying out a function that returns two variables: ", plus2_1, plus2_2)

	// variadic functions can take any number of parameters
	sum(1, 2, 3)
	nums := []int{1, 2, 3, 4}
	sum(nums...) // if you have a slice, you can use the ... to pass it as a parameter

	// closures
	nextInt := intSeq()                                                          // nextInt is a function that returns an int
	fmt.Println("calling intSeq three times: ", nextInt(), nextInt(), nextInt()) // use the function, the i gets incremented by 1
	newInts := intSeq()                                                          // defining a new function, not that the state of nextInt is unique
	fmt.Println("starting the intSeq from zero: ", newInts())                    // i starts from 1 again

	// recursive functions
	fmt.Println("preforming a recursive function fact(n*n-1): ", fact(7))
	// Anonymous functions can also be recursive, but this requires explicitly
	// declaring a variable with var to store the function before it’s defined.
	var fib func(n int) int
	fib = func(n int) int {
		if n < 2 {
			return n
		}

		// Since fib was previously declared in main, Go knows which function to call with fib here
		return fib(n-1) + fib(n-2)
	}
	fmt.Println("fib(7): ", fib(7))

	// Using the Add function from another file under another folder/package
	result := helper.Add(3, 4)
	fmt.Println("Result:", result)

	/*********************************** range *************************************************/

	// range can be performed on strings, arrays, maps, channels, and slices
	for i, num := range nums { // i is the index and num is the value, note that you can use _ to ignore the index
		if num == 3 {
			fmt.Println("found 3 at index:", i)
		}
	}
	kvs := map[string]string{"a": "apple", "b": "banana"}
	for k, v := range kvs { // k is the key and v is the value, you can omit the v and only get the keys
		fmt.Printf("%s -> %s\n", k, v)
	}
	for i, c := range "go" { // i is the index, c is the byte value of the char
		fmt.Println("ranging over the string: ", i, c, string(c))
	}

	/*********************************** pointers *************************************************/

	var zeroval = func(ival int) {
		ival = 0
	}
	var zeroptr = func(iptr *int) {
		*iptr = 0
	}
	ivalue := 1 // create a variable ivalue and assign it with type int, value 1
	fmt.Println("initially, ivalue is: ", ivalue)
	zeroval(ivalue) // use zeroval on the value, nothing changes because the function received a copy of ivalue and only changed the copy
	fmt.Println("after zeroval, ivalue is:", ivalue)
	zeroptr(&ivalue) // using the & to pass the address of ivalue, this will actually change ivalue
	fmt.Println("after zeroptr, ivalue is:", ivalue)
	fmt.Println("address of ivalue is: ", &ivalue) // access the address of an object using &

	/*********************************** strings and runes *************************************************/
	// a go string is a read-only slice of bytes, use len() to get length of string
	svar := "Hello你是谁"

	for i := 0; i < len(svar); i++ { // use len to get the length of the string, note that the length is not the same as the byte length
		fmt.Printf("%x ", svar[i]) // cast to string to print actual char, does not correctly print chinese chars , use printf for that
	}
	fmt.Println("")
	fmt.Println("Rune count for svar:", utf8.RuneCountInString(svar)) // use runecount to get the number of runes in the string
	for idx, runeValue := range svar {                                // use range to automatcally get the rune value and the index
		fmt.Printf("%#U starts at %d\n", runeValue, idx)
	}
	runeValue, width := utf8.DecodeRuneInString(svar[5:]) // use decoderuneinstring and give the index of the rune to get the decoded rune
	fmt.Println("rune value:", runeValue, "width:", width)
	// Values enclosed in single quotes are rune literals. We can compare a rune value to a rune literal directly.
	if runeValue == '你' {
		fmt.Println("found character 你")
	}

	// example of functions supported by the strings package
	helper.ShowStringPackageExample()

	/*********************************** structs *************************************************/

	// naming the fields is optional, but it's best to do so
	fmt.Println("here's a rectangle: ", rect{width: 10, height: 20, name: "rectangle1"})
	// ommitted fields will be zero valued
	rectangle := rect{width: 10, height: 5}
	// will print the source code that produces the rectangle object
	fmt.Printf("rectangle test: %#v\n", rectangle)
	// structs are mutable
	rectangle.height = 20
	// use . to access fields
	// an & prefix yields a pointer to a struct
	fmt.Println("printing the pointer to rectangle", &rectangle, "the width of the rectangle is: ", rectangle.width)
	fmt.Println("creating a new rectangle:", newRect("rectangle2"))
	// pointer to a struct
	rectanglePointer := &rectangle
	// pointers are automatically dereferenced
	fmt.Println("using a pointer to access the rectangle's height", rectanglePointer.height)

	// creating an anonymous struct
	dog := struct {
		name   string
		isGood bool
	}{
		"Rex",
		true,
	}
	fmt.Println("dog struct is: ", dog)

	fmt.Println("rectange area is: ", rectangle.area())
	fmt.Println("rectangle perim is: ", rectangle.perim())

	// go automatically handles conversions between values and pointers for method calls
	// but you may want to use a pointer receiver type to avoid copying on method calls or to allow the method to mutate the receiving struct
	rectanglepointer := &rectangle
	fmt.Println("rectangle area is: ", rectanglepointer.area())
	fmt.Println("rectangle perimeter is: ", rectanglepointer.perim())

	// using interface
	circle1 := circle{radius: 5}
	measure(circle1)

	// enums
	// ns is a StateIdle ServerState, the String() function now returns the string representation of the state
	nowstate := StateIdle
	fmt.Println("ns is now: ", nowstate)
	transstate := transition(nowstate)
	fmt.Println("after transition, ns is now: ", transstate)

	// if a struct embeds another struct, this is how you create it
	containerobject := container{
		base: base{
			num: 1,
		},
		str: "some name",
	}

	// you can access the embedded struct's fields through the embedded struct or directly without
	fmt.Printf("co={num: %v,alsonum: %v, str: %v}\n", containerobject.base.num, containerobject.num, containerobject.str)
	fmt.Println("describe:", containerobject.describe()) // since co embeds base, it can access the base's describe function

	type describer interface {
		describe() string
	}

	//Embedding structs with methods may be used to bestow interface implementations onto other structs.
	//Here we see that a container now implements the describer interface because it embeds base.
	var des describer = containerobject
	fmt.Println("describer:", des.describe())

	/*********************************** generic functions  *************************************************/
	// the generic functions part of the code requires "go run tutorial.go main.go" to successfully link up the two files
	// we can use "go run ." to simplify the process

	var sList = []string{"foo", "bar", "zoo"}
	// note that we don't have to specify the type of the slice, it's inferred automatically
	fmt.Println("index of zoo:", SlicesIndex(sList, "zoo"))
	// though you can specify the type of the slice, it's not necessary
	_ = SlicesIndex[[]string, string](sList, "zoo")

	// creating a int list
	lst := List[int]{}
	lst.Push(10) // push into list
	lst.Push(13)
	lst.Push(23)
	fmt.Println("list.AllElements :", lst.AllElements()) // use all elements to display list

	// use range on iterator
	for e := range lst.All() {
		fmt.Println("using range on iterator: ", e)
	}

	all := slices.Collect(lst.All()) // Collect takes any iterator and collects all its values into a slice
	fmt.Println("all:", all)

	// iterators
	for n := range genFib() {
		if n >= 10 {
			break // once the loop hits break or an early return, the yield function passed to the iterator will return false
		}
		fmt.Println("using iterator to generate fibbonacci sequence: ", n)
	}
	/*********************************** errors *************************************************/
	for _, i := range []int{7, 42} {

		if r, e := fun(i); e != nil { // inline error check, consist of init and is nil
			fmt.Println("f failed:", e) // print error
		} else {
			fmt.Println("f worked:", r)
		}
	}

	for i := range 5 {
		if err := makeTea(i); err != nil {
			if errors.Is(err, ErrOutOfTea) { // check to see if err is ErrOutOfTea
				fmt.Println("We should buy new tea!")
			} else if errors.Is(err, ErrPower) { // notice that when i=4,we returned a higher error that included ErrPower, but here it is captured
				fmt.Println("Now it is dark.")
			} else {
				fmt.Printf("unknown error: %s\n", err)
			}
			continue
		}
		fmt.Println("Tea is ready!")
	}

	_, err := fun2(42)
	var ae *argError
	// errors.As checks that a given error(or any error in its chain) matches
	// a specific error type and converts to a value of that type,returning ture.
	// if there's no match, return false
	if errors.As(err, &ae) {
		fmt.Println("argument: ", ae.arg)
		fmt.Println("message: ", ae.message)
	} else {
		fmt.Println("err doesn't match argError")
	}

	/*********************************** sorting  *************************************************/

	// sorting functions are generic, and work for any ordered built in type
	// sorting strings
	strs := []string{"c", "a", "b"}
	slices.Sort(strs)
	fmt.Println("Strings:", strs)
	// sorting ints
	ints := []int{7, 2, 4}
	slices.Sort(ints)
	fmt.Println("Ints:   ", ints)
	// We can also use the slices package to check if a slice is already in sorted order.
	sv := slices.IsSorted(ints)
	fmt.Println("is Sorted: ", sv)

	// using custom sorting criteria
	fruits := []string{"peach", "banana", "kiwi"}
	//  create custom compare criteria function
	lenCmp := func(a, b string) int {
		return cmp.Compare(len(a), len(b))
	}
	// sort the slice using custom function
	slices.SortFunc(fruits, lenCmp)
	fmt.Println("sorted fruits via string length: ", fruits)

	// create a person type
	type Person struct {
		name string
		age  int
	}

	// create a person slice
	people := []Person{
		Person{name: "Jax", age: 37},
		Person{name: "TJ", age: 25},
		Person{name: "Alex", age: 72},
	}

	// sort the people slice using age
	// note that f the person struct is large, you may want the slice to contain
	// *person instead and adjusting the sorting function accordingly
	slices.SortFunc(people,
		func(a, b Person) int {
			// compare returns -1 if x<y, 0 if x==y, 1 if x>y
			return cmp.Compare(a.age, b.age)
		})
	fmt.Println("sorted people slice: ", people)

	// the following are examples of packages that we might use on a regular basis
	// these functions are in the helper.go file to better format the code.
	helper.ShowRegularExpressionExample()
	helper.ShowLoggerExample()
	helper.ShowTextTemplateExample()
	helper.ShowSha256Example()
	helper.ShowRandExample()
	helper.ShowNumberParsingExample()
	helper.ShowBase64EncodingExample()
	helper.ShowJsonExample()
	helper.ShowUrlParsingExample()
	helper.ShowTimeExample()
	helper.ShowStringFormattingExample()
	helper.ShowXMLExample()

}
