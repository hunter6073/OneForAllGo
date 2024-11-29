package main

import (
	"fmt"
	"testing"
)

// this is the code to be tested
func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// test function, t is a pointer to a testing.T object
func TestIntMinBasic(t *testing.T) {
	ans := IntMin(2, -2) // run the function
	if ans != -2 {       // check against answer, if answer is wrong, call t.Errorf
		// t.errorf will report test failure but continue executing the test, to stop the test immediately, call t.Fatal
		t.Errorf("IntMin(2, -2) = %d; want -2", ans)
	}
}

// it is idiomatic to use a table to store test cases
func TestIntMinTableDriven(t *testing.T) {
	// tests is an array of structs, each of which represents a test case
	var tests = []struct {
		a, b int
		want int
	}{
		{0, 1, 0},
		{1, 0, 0},
		{2, -2, -2},
		{0, -1, -1},
		{-1, 0, -1},
	}
	// use loop to run the test cases
	for _, tt := range tests {
		//print test case
		testname := fmt.Sprintf("%d,%d", tt.a, tt.b)
		// run test
		t.Run(testname, func(t *testing.T) {
			ans := IntMin(tt.a, tt.b)
			if ans != tt.want {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
}

// bench mark tests. the testing runner executes each benchmark function several times, increasing b.N on each run until it collects a precise measurement.
func BenchmarkIntMin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IntMin(1, 2)
	}
}

// use go test -v to run all tests in the current project in verbose mode
// use go test -bench=. to run all the benchmark tests in the current project. all tests are run prior to benchmarks
