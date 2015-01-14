// xfind project  pred.go
// Copyright 2015 atarrow. All rights reserved.
package main

import (
	"os"
	"regexp"
	"strings"
)

// Predicate is used to define logical conditions used to constrain a search
type Predicate func(os.FileInfo) bool

// useful predicates:
var (
	isDotFile    Predicate = func(v os.FileInfo) bool { return strings.HasPrefix(v.Name(), ".") }
	isNotDotFile Predicate = func(v os.FileInfo) bool { return !isDotFile(v) }
	isDirOnly    Predicate = func(v os.FileInfo) bool { return v.IsDir() }
	isFileOnly   Predicate = func(v os.FileInfo) bool { return !v.IsDir() }
)

// PredicateWithPattern creates a new predicate that represents
// whether an argument matches a regular expression pattern.
func PredicateWithPattern(pattern string) (Predicate, error) {
	re, err := regexp.Compile(NFDString(pattern))
	if err != nil {
		return nil, err
	}
	pred := func(v os.FileInfo) bool { return re.MatchString(NFDString(v.Name())) }
	return pred, nil
}
