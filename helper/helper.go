package helper

import (
	"bytes"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"math/rand/v2"
	"net"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	sObj "strings"
	"time"
)

// Add is a function that adds two integers and returns the result.
// !!! note that if you're defining a function in another file outside main, the function must
// start with a capital case letter, otherwise the compiler cannot link to the external function
func Add(a, b int) int {
	return a + b
}

// we'll use these two structs to demonstrate encoding and decoding of custom types below
type response1 struct {
	Page   int
	Fruits []string
}

// only exported fields will be encoded/decoded in json, fields must start with capital letters to be exported
type response2 struct {
	Page   int `json:"page"`
	Fruits []string
}

type point struct {
	x, y int
}

// Plant will be mapped to XML. Similarly to the JSON examples,
// field tags contain directives for the encoder and decoder.
// Here we use some special features of the XML package: the XMLName field name
// dictates the name of the XML element representing this struct;
// id,attr means that the Id field is an XML attribute rather than a nested element.
type Plant struct {
	XMLName xml.Name `xml:"plant"`
	Id      int      `xml:"id,attr"`
	Name    string   `xml:"name"`
	Origin  []string `xml:"origin"`
}

func (p Plant) String() string {
	return fmt.Sprintf("Plant id=%v, name=%v, origin=%v",
		p.Id, p.Name, p.Origin)
}

func ShowStringPackageExample() {
	var p = fmt.Println // alias fmt.Println to a shorter name

	// functions that strings package contains

	// we can give the strings package an alias
	p("test Contains es:  ", sObj.Contains("test", "es"))
	// or directly use strings
	p("Count numbers of e in test:     ", strings.Count("test", "e"))
	p("test HasPrefix te: ", strings.HasPrefix("test", "te"))
	p("test HasSuffix st: ", strings.HasSuffix("test", "st"))
	p("Index of e in test:     ", strings.Index("test", "e"))
	p("Join a and b with -:      ", strings.Join([]string{"a", "b"}, "-"))
	p("Repeat a five times:    ", strings.Repeat("a", 5))
	p("Replace o with 0:   ", strings.Replace("foo", "o", "0", -1))
	p("Replace o with 0 starting from the second o:   ", strings.Replace("foo", "o", "0", 1))
	p("Split a-b-c-d-e using -:     ", strings.Split("a-b-c-d-e", "-"))
	p("ToLower TEST:   ", strings.ToLower("TEST"))
	p("ToUpper test:   ", strings.ToUpper("test"))
}

func ShowRegularExpressionExample() {
	// this test whether a pattern matches a string
	match, _ := regexp.MatchString("p([a-z]+)ch", "peach")
	fmt.Println("p([a-z]+)ch matches peach", match)

	// Above we used a string pattern directly, but for other regexp tasks
	// you’ll need to Compile an optimized Regexp struct.
	r, _ := regexp.Compile("p([a-z]+)ch")
	fmt.Println("p([a-z]+)ch matches peach", r.MatchString("peach")) // perform match test

	// This finds the match of a string for the regexp from the given string.
	fmt.Println(r.FindString("peach punch"))

	// This also finds the first match but returns the start and
	// end indexes for the match instead of the matching text.
	fmt.Println("idx for p([a-z]+)ch in peach punch:", r.FindStringIndex("peach punch"))

	// The Submatch variants include information about both
	// the whole-pattern matches and the submatches within those matches.
	// For example this will return information for both p([a-z]+)ch and ([a-z]+).
	fmt.Println(r.FindStringSubmatch("peach punch"))

	// Similarly this will return information about the indexes of matches and submatches.
	fmt.Println(r.FindStringSubmatchIndex("peach punch"))

	// The All variants of these functions apply to all matches in the input,
	// not just the first. For example to find all matches for a regexp.
	fmt.Println(r.FindAllString("peach punch pinch", -1))

	// These All variants are available for the other functions we saw above as well.
	fmt.Println("all:", r.FindAllStringSubmatchIndex("peach punch pinch", -1))

	// Providing a non-negative integer as the second argument
	// to these functions will limit the number of matches.
	fmt.Println(r.FindAllString("peach punch pinch", 2))

	// Our examples above had string arguments and used names like MatchString.
	// We can also provide []byte arguments and drop String from the function name.
	fmt.Println(r.Match([]byte("peach")))

	// When creating global variables with regular expressions you can use the
	// MustCompile variation of Compile. MustCompile panics instead of returning an error,
	// which makes it safer to use for global variables.
	r = regexp.MustCompile("p([a-z]+)ch")
	fmt.Println("regexp:", r)

	// The regexp package can also be used to replace subsets of strings with other values.
	fmt.Println(r.ReplaceAllString("a peach", "<fruit>"))

	// The Func variant allows you to transform matched text with a given function.
	in := []byte("a peach")
	out := r.ReplaceAllFunc(in, bytes.ToUpper)
	fmt.Println(string(out))
}

func ShowLoggerExample() {
	// use log to print messages
	log.Println("standard logger")

	// configure logger with flags to set their output format
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Println("with micro")

	// also support emitting the file name and line from which the log was called.
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("with file/line")

	// it may be useful to create a custom logger and pass it around
	// we are able to set a prefix to distinguish its output from other loggers
	mylog := log.New(os.Stdout, "my:", log.LstdFlags)
	mylog.Println("from mylog")

	// we can set the prefix on existing loggers with the set prefix method
	mylog.SetPrefix("ohmy:")
	mylog.Println("from mylog")

	// loggers can have custom output targets; any io.writer works
	var buf bytes.Buffer
	buflog := log.New(&buf, "buf:", log.LstdFlags)

	// this call writes the log output into buf
	buflog.Println("hello")

	// this will show it to standard output
	fmt.Print("from buflog:", buf.String())

	// slog provies structured log output. for example, logging in json format is straightforward
	jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
	myslog := slog.New(jsonHandler)
	myslog.Info("hi there")

	// in addition to the message, slog output can contain an arbitrary number of keyvalue pairs
	myslog.Info("hello again", "key", "val", "age", 25)

}

func ShowTextTemplateExample() {
	// We can create a new template and parse its body from a string.
	// Templates are a mix of static text and “actions” enclosed in {{...}} that are used to dynamically insert content.
	t1 := template.New("t1")
	// we can use template.Must to panic in case parse returns an error. this is especially useful for templates initialized in the global scope.
	t1 = template.Must(t1.Parse("Value: {{.}}\n"))
	// by executing the template we generate its text with specific values for its actions.
	// the {{.}} action is replaced by the value passed a s a parameter to execute
	t1.Execute(os.Stdout, "some text")
	t1.Execute(os.Stdout, 5)
	t1.Execute(os.Stdout, []string{
		"Go",
		"Rust",
		"C++",
		"C#",
	})
	// create a helper function
	Create := func(name, t string) *template.Template {
		return template.Must(template.New(name).Parse(t))
	}

	// If the data is a struct we can use the {{.FieldName}} action to access its fields.
	// The fields should be exported to be accessible when a template is executing.
	t2 := Create("t2", "Name: {{.Name}}\n")
	t2.Execute(os.Stdout, struct {
		Name string
	}{"Jane Doe"})

	// The same applies to maps; with maps there is no restriction on the case of key names.
	t2.Execute(os.Stdout, map[string]string{
		"Name": "Mickey Mouse",
	})

	// if/else provide conditional execution for templates. A value is considered false if it’s
	//the default value of a type, such as 0, an empty string, nil pointer, etc.
	//This sample demonstrates another feature of templates: using - in actions to trim whitespace.
	t3 := Create("t3",
		"{{if . -}} yes {{else -}} no {{end}}\n")
	t3.Execute(os.Stdout, "not empty")
	t3.Execute(os.Stdout, "")

	// range blocks let us loop through slices, arrays, maps or channels.
	// Inside the range block {{.}} is set to the current item of the iteration.
	t4 := Create("t4",
		"Range: {{range .}}{{.}} {{end}}\n")
	t4.Execute(os.Stdout,
		[]string{
			"Go",
			"Rust",
			"C++",
			"C#",
		})
}

func ShowSha256Example() {
	s4 := "sha256 this string"
	h := sha256.New()                       // start a new hash
	h.Write([]byte(s4))                     // write expects bytes, if you have string, use []byte(s) to convert a string
	bs := h.Sum(nil)                        // this gets  the finalized hash result as a byte slice. the argument to sum can be used to append to an existing byte slice, it usually isn't needed
	fmt.Println("original string is: ", s4) // print the original string
	fmt.Printf("the hash is: %x\n", bs)     // print the hash result
}

// TODO: add print message
func ShowRandExample() {
	// For example, rand.IntN returns a random int n, 0 <= n < 100.
	fmt.Print(rand.IntN(100), ",")
	fmt.Print(rand.IntN(100))
	fmt.Println()
	// rand.Float64 returns a float64 f, 0.0 <= f < 1.0.
	fmt.Println(rand.Float64())
	// This can be used to generate random floats in other ranges, for example 5.0 <= f' < 10.0.
	fmt.Print((rand.Float64()*5)+5, ",")
	fmt.Print((rand.Float64() * 5) + 5)
	fmt.Println()

	// If you want a known seed, create a new rand.Source and pass it into the New constructor.
	// NewPCG creates a new PCG source that requires a seed of two uint64 numbers.
	s2 := rand.NewPCG(42, 1024)
	r2 := rand.New(s2)
	fmt.Print(r2.IntN(100), ",")
	fmt.Print(r2.IntN(100))
	fmt.Println()
	s3 := rand.NewPCG(42, 1024)
	r3 := rand.New(s3)
	fmt.Print(r3.IntN(100), ",")
	fmt.Print(r3.IntN(100))
	fmt.Println()
}

func ShowNumberParsingExample() {
	// With ParseFloat, this 64 tells how many bits of precision to parse.
	f, _ := strconv.ParseFloat("1.234", 64)
	fmt.Println(f)

	// For ParseInt, the 0 means infer the base from the string. 64 requires that the result fit in 64 bits.
	i, _ := strconv.ParseInt("123", 0, 64)
	fmt.Println(i)

	//ParseInt will recognize hex-formatted numbers.
	d, _ := strconv.ParseInt("0x1c8", 0, 64)
	fmt.Println(d)

	// A ParseUint is also available.
	u, _ := strconv.ParseUint("789", 0, 64)
	fmt.Println(u)

	// Atoi is a convenience function for basic base-10 int parsing.
	k, _ := strconv.Atoi("135")
	fmt.Println(k)

	// Parse functions return an error on bad input.
	_, e := strconv.Atoi("wat")
	fmt.Println(e)
}

func ShowBase64EncodingExample() {
	// Here’s the string we’ll encode/decode.
	data := "abc123!?$*&()'-=@~"
	// Go supports both standard and URL-compatible base64.
	// Here’s how to encode using the standard encoder.
	// The encoder requires a []byte so we convert our string to that type.
	sEnc := b64.StdEncoding.EncodeToString([]byte(data))
	fmt.Println(sEnc)

	//Decoding may return an error, which you can check if you don’t already know the input to be well-formed.
	sDec, _ := b64.StdEncoding.DecodeString(sEnc)
	fmt.Println(string(sDec))
	fmt.Println()
	// This encodes/decodes using a URL-compatible base64 format.
	uEnc := b64.URLEncoding.EncodeToString([]byte(data))
	fmt.Println(uEnc)
	uDec, _ := b64.URLEncoding.DecodeString(uEnc)
	fmt.Println(string(uDec))
}

func ShowUrlParsingExample() {
	// We’ll parse this example URL, which includes a scheme, authentication info, host, port, path, query params, and query fragment.
	s := "postgres://user:pass@host.com:5432/path?k=v#f"
	//Parse the URL and ensure there are no errors.
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	// Accessing the scheme is straightforward.
	fmt.Println(u.Scheme)
	// User contains all authentication info; call Username and Password on this for individual values.
	fmt.Println(u.User)
	fmt.Println(u.User.Username())
	p, _ := u.User.Password()
	fmt.Println(p)
	// The Host contains both the hostname and the port, if present. Use SplitHostPort to extract them.
	fmt.Println(u.Host)
	host, port, _ := net.SplitHostPort(u.Host)
	fmt.Println(host)
	fmt.Println(port)
	// Here we extract the path and the fragment after the #.
	fmt.Println(u.Path)
	fmt.Println(u.Fragment)
	// To get query params in a string of k=v format, use RawQuery.
	// You can also parse query params into a map. The parsed query param maps
	// are from strings to slices of strings, so index into [0] if you only want the first value.
	fmt.Println(u.RawQuery)
	m, _ := url.ParseQuery(u.RawQuery)
	fmt.Println(m)
	fmt.Println(m["k"][0])
}

func ShowJsonExample() {
	// First we’ll look at encoding basic data types to JSON strings. Here are some examples for atomic values.
	bolB, _ := json.Marshal(true)
	fmt.Println(string(bolB))
	intB, _ := json.Marshal(1)
	fmt.Println(string(intB))
	fltB, _ := json.Marshal(2.34)
	fmt.Println(string(fltB))
	strB, _ := json.Marshal("gopher")
	fmt.Println(string(strB))
	// And here are some for slices and maps, which encode to JSON arrays and objects as you’d expect.
	slcD := []string{"apple", "peach", "pear"}
	slcB, _ := json.Marshal(slcD)
	fmt.Println(string(slcB))
	mapD := map[string]int{"apple": 5, "lettuce": 7}
	mapB, _ := json.Marshal(mapD)
	fmt.Println(string(mapB))
	// The JSON package can automatically encode your custom data types.
	// It will only include exported fields in the encoded output and will by default use those names as the JSON keys.
	res1D := &response1{
		Page:   1,
		Fruits: []string{"apple", "peach", "pear"}}
	res1B, _ := json.Marshal(res1D)
	fmt.Println(string(res1B))
	// You can use tags on struct field declarations to customize the encoded JSON key names.
	// Check the definition of response2 above to see an example of such tags.
	res2D := &response2{
		Page:   1,
		Fruits: []string{"apple", "peach", "pear"}}
	res2B, _ := json.Marshal(res2D)
	fmt.Println(string(res2B))
	// Now let’s look at decoding JSON data into Go values. Here’s an example for a generic data structure.
	byt := []byte(`{"num":6.13,"strs":["a","b"]}`)
	// We need to provide a variable where the JSON package can put the decoded data.
	// This map[string]interface{} will hold a map of strings to arbitrary data types.
	var dat map[string]interface{}
	// Here’s the actual decoding, and a check for associated errors.
	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat)
	// In order to use the values in the decoded map, we’ll need to convert them to their appropriate type.
	// For example here we convert the value in num to the expected float64 type.
	num := dat["num"].(float64)
	fmt.Println(num)
	// Accessing nested data requires a series of conversions.
	strs := dat["strs"].([]interface{})
	str1 := strs[0].(string)
	fmt.Println(str1)
	// We can also decode JSON into custom data types. This has the advantages of adding additional
	// type-safety to our programs and eliminating the need for type assertions when accessing the decoded data.
	str := `{"page": 1, "fruits": ["apple", "peach"]}`
	res := response2{}
	json.Unmarshal([]byte(str), &res)
	fmt.Println(res)
	fmt.Println(res.Fruits[0])
	// In the examples above we always used bytes and strings as intermediates
	// between the data and JSON representation on standard out.
	// We can also stream JSON encodings directly to os.Writers like os.Stdout or even HTTP response bodies.
	enc := json.NewEncoder(os.Stdout)
	d := map[string]int{"apple": 5, "lettuce": 7}
	enc.Encode(d)
}

func ShowTimeExample() {
	p := fmt.Println
	// We’ll start by getting the current time.
	now := time.Now()
	p(now)
	// You can build a time struct by providing the year, month, day, etc. Times are always associated with a Location, i.e. time zone.
	then := time.Date(
		2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	p(then)
	// You can extract the various components of the time value as expected.
	p(then.Year())
	p(then.Month())
	p(then.Day())
	p(then.Hour())
	p(then.Minute())
	p(then.Second())
	p(then.Nanosecond())
	p(then.Location())
	// The Monday-Sunday Weekday is also available.
	p(then.Weekday())
	// These methods compare two times, testing if the first occurs before, after, or at the same time as the second, respectively.
	p(then.Before(now))
	p(then.After(now))
	p(then.Equal(now))
	// The Sub methods returns a Duration representing the interval between two times.
	diff := now.Sub(then)
	p(diff)
	// We can compute the length of the duration in various units.
	p(diff.Hours())
	p(diff.Minutes())
	p(diff.Seconds())
	p(diff.Nanoseconds())
	// You can use Add to advance a time by a given duration, or with a - to move backwards by a duration.
	p(then.Add(diff))
	p(then.Add(-diff))

	// Here’s a basic example of formatting a time according to RFC3339, using the corresponding layout constant.
	t := time.Now()
	p(t.Format(time.RFC3339))
	// Time parsing uses the same layout values as Format.
	t1, e := time.Parse(
		time.RFC3339,
		"2012-11-01T22:08:41+00:00")
	p(t1)
	// Format and Parse use example-based layouts. Usually you’ll use a constant from time for these layouts,
	// but you can also supply custom layouts. Layouts must use the reference time Mon Jan 2 15:04:05 MST 2006
	// to show the pattern with which to format/parse a given time/string.
	// The example time must be exactly as shown: the year 2006, 15 for the hour, Monday for the day of the week, etc.
	p(t.Format("3:04PM"))
	p(t.Format("Mon Jan _2 15:04:05 2006"))
	p(t.Format("2006-01-02T15:04:05.999999-07:00"))
	form := "3 04 PM"
	t2, e := time.Parse(form, "8 41 PM")
	p(t2)
	// For purely numeric representations you can also use standard string formatting with the extracted components of the time value.
	fmt.Printf("%d-%02d-%02dT%02d:%02d:%02d-00:00\n",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	// Parse will return an error on malformed input explaining the parsing problem.
	ansic := "Mon Jan _2 15:04:05 2006"
	_, e = time.Parse(ansic, "8:41PM")
	p(e)

	// Use time.Now with Unix, UnixMilli or UnixNano to get elapsed
	// time since the Unix epoch in seconds, milliseconds or nanoseconds, respectively.
	fmt.Println(now)
	fmt.Println(now.Unix())
	fmt.Println(now.UnixMilli())
	fmt.Println(now.UnixNano())
	// You can also convert integer seconds or nanoseconds since the epoch into the corresponding time.
	fmt.Println(time.Unix(now.Unix(), 0))
	fmt.Println(time.Unix(0, now.UnixNano()))
}

func ShowStringFormattingExample() {
	// Go offers several printing “verbs” designed to format general Go values.
	// For example, this prints an instance of our point struct.
	p := point{1, 2}
	fmt.Printf("struct1: %v\n", p)

	// If the value is a struct, the %+v variant will include the struct’s field names.
	fmt.Printf("struct2: %+v\n", p)

	// The %#v variant prints a Go syntax representation of the value,
	// i.e. the source code snippet that would produce that value.
	fmt.Printf("struct3: %#v\n", p)

	// To print the type of a value, use %T.
	fmt.Printf("type: %T\n", p)

	// Formatting booleans is straight-forward.
	fmt.Printf("bool: %t\n", true)

	// There are many options for formatting integers. Use %d for standard, base-10 formatting.
	fmt.Printf("int: %d\n", 123)

	// This prints a binary representation.
	fmt.Printf("bin: %b\n", 14)

	// This prints the character corresponding to the given integer.
	fmt.Printf("char: %c\n", 33)

	// %x provides hex encoding.
	fmt.Printf("hex: %x\n", 456)

	// There are also several formatting options for floats. For basic decimal formatting use %f.
	fmt.Printf("float1: %f\n", 78.9)

	// %e and %E format the float in (slightly different versions of) scientific notation.
	fmt.Printf("float2: %e\n", 123400000.0)
	fmt.Printf("float3: %E\n", 123400000.0)

	// For basic string printing use %s.
	fmt.Printf("str1: %s\n", "\"string\"")

	// To double-quote strings as in Go source, use %q.
	fmt.Printf("str2: %q\n", "\"string\"")

	// As with integers seen earlier, %x renders the string
	// in base-16, with two output characters per byte of input.
	fmt.Printf("str3: %x\n", "hex this")

	// To print a representation of a pointer, use %p.
	fmt.Printf("pointer: %p\n", &p)

	// When formatting numbers you will often want to control the width and precision of the resulting figure.
	// To specify the width of an integer, use a number after the % in the verb.
	// By default the result will be right-justified and padded with spaces.
	fmt.Printf("width1: |%6d|%6d|\n", 12, 345)

	// You can also specify the width of printed floats, though usually you’ll
	// also want to restrict the decimal precision at the same time with the width.precision syntax.
	fmt.Printf("width2: |%6.2f|%6.2f|\n", 1.2, 3.45)

	// To left-justify, use the - flag.
	fmt.Printf("width3: |%-6.2f|%-6.2f|\n", 1.2, 3.45)

	// You may also want to control width when formatting strings,
	// especially to ensure that they align in table-like output. For basic right-justified width.
	fmt.Printf("width4: |%6s|%6s|\n", "foo", "b")

	// To left-justify use the - flag as with numbers.
	fmt.Printf("width5: |%-6s|%-6s|\n", "foo", "b")

	// So far we’ve seen Printf, which prints the formatted string to os.Stdout.
	// Sprintf formats and returns a string without printing it anywhere.
	s := fmt.Sprintf("sprintf: a %s", "string")
	fmt.Println(s)

	// You can format+print to io.Writers other than os.Stdout using Fprintf.
	fmt.Fprintf(os.Stderr, "io: an %s\n", "error")
}

func ShowXMLExample() {

	coffee := &Plant{Id: 27, Name: "Coffee"}
	coffee.Origin = []string{"Ethiopia", "Brazil"}

	// Emit XML representing our plant; using MarshalIndent to produce a more human-readable output.
	out, _ := xml.MarshalIndent(coffee, " ", "  ")
	fmt.Println(string(out))

	// To add a generic XML header to the output, append it explicitly.
	fmt.Println(xml.Header + string(out))

	// Use Unmarshal to parse a stream of bytes with XML into a data structure.
	// If the XML is malformed or cannot be mapped onto Plant, a descriptive error will be returned.
	var p Plant
	if err := xml.Unmarshal(out, &p); err != nil {
		panic(err)
	}
	fmt.Println(p)
	tomato := &Plant{Id: 81, Name: "Tomato"}
	tomato.Origin = []string{"Mexico", "California"}

	// The parent>child>plant field tag tells the encoder to nest all plants under <parent><child>...
	type Nesting struct {
		XMLName xml.Name `xml:"nesting"`
		Plants  []*Plant `xml:"parent>child>plant"`
	}
	nesting := &Nesting{}
	nesting.Plants = []*Plant{coffee, tomato}
	out, _ = xml.MarshalIndent(nesting, " ", "  ")
	fmt.Println(string(out))
}
