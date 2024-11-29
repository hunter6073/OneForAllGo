package main

import "iter"

// As an example of a generic function, SlicesIndex takes a slice of any comparable
// type and an element of that type and returns the index of the first occurrence of
// v in s, or -1 if not present. The comparable constraint means that we can compare
//  values of this type with the == and != operators
func SlicesIndex[S ~[]E, E comparable](s S, v E) int {
	for i := range s {
		if v == s[i] {
			return i
		}
	}
	return -1
}

// As an example of a generic type, List is a singly-linked list with values of any type.
type List[T any] struct {
	head, tail *element[T]
}

type element[T any] struct {
	next *element[T]
	val  T
}

//We can define methods on generic types just like we do on regular types,
// but we have to keep the type parameters in place. The type is List[T], not List.
func (lst *List[T]) Push(v T) {
	if lst.tail == nil {
		lst.head = &element[T]{val: v}
		lst.tail = lst.head
	} else {
		lst.tail.next = &element[T]{val: v}
		lst.tail = lst.tail.next
	}
}

// AllElements returns all the List elements as a slice.
func (lst *List[T]) AllElements() []T {
	var elems []T
	for e := lst.head; e != nil; e = e.next {
		elems = append(elems, e.val)
	}
	return elems
}

// all returns an iterator
func (lst *List[T]) All() iter.Seq[T] {
	// the iterator function takes another function as a parameter,
	// it will call yield for every element we want to iterate over
	return func(yield func(T) bool) {
		for e := lst.head; e != nil; e = e.next {
			if !yield(e.val) {
				return
			}
		}
	}
}

// iteration doesn't require underlying data structure
func genFib() iter.Seq[int] {
	return func(yield func(int) bool) {
		a, b := 1, 1

		for { // doesn't need to be finite, can keep on yielding
			if !yield(a) {
				return
			}
			a, b = b, a+b
		}
	}
}
