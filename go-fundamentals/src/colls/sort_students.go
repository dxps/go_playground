package main

import (
	"fmt"
	"sort"
)

// Student is a struct containing minimal details about a student.
type Student struct {
	Name  string
	Grade float32
}

func (s Student) String() string {
	return fmt.Sprintf("%s: %f", s.Name, s.Grade)
}

// ByGrade is a type alias that implements the Sort sort.Interface.
type ByGrade []Student

// Len returns the length of ByGrade. It is one or the methods implemented as per sort.Interface.
func (coll ByGrade) Len() int {
	return len(coll)
}

// Less reports whether the element with index i should sort before the element with index j.
func (coll ByGrade) Less(i, j int) bool {
	return coll[i].Grade < coll[j].Grade
}

// Swap swaps the elements with indexes i and j.
func (coll ByGrade) Swap(i, j int) {
	coll[i], coll[j] = coll[j], coll[i]
}

func main() {
	students := []Student{
		{"Bob", 9.2},
		{"Alice", 7.5},
		{"John", 8.3},
		{"Tom", 6.3},
	}
	fmt.Printf(" (initial set) students: %s\n", students)
	sort.Sort(ByGrade(students))
	fmt.Printf(" (sorted by grade) students: %s\n", students)

}
