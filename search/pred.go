// xfind project  pred.go
// Copyright 2015 atarrow. All rights reserved.
package search

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
	isNotDotFile           = func(v os.FileInfo) bool { return !isDotFile(v) }
	isDirOnly              = func(v os.FileInfo) bool { return v.IsDir() }
	isFileOnly             = func(v os.FileInfo) bool { return !v.IsDir() }
	AlwaysTrue             = func(v os.FileInfo) bool { return true }
)

// PredicateWithPattern creates a new predicate that represents
// whether an argument matches a regular expression pattern.
func PredicateWithPattern(pattern string) (Predicate, error) {
	re, err := regexp.Compile(NFDString(pattern))
	if err != nil {
		return nil, err
	}
	pred := func(v os.FileInfo) bool {
		return re.MatchString(
			NFDString(v.Name()),
		)
	}
	return pred, nil
}

type Predicates []Predicate

// Pred return the predictate that is a union of all predicates in it.
func (preds Predicates) ANDPredicate() Predicate {
	return func(v os.FileInfo) bool {
		if len(preds) == 0 {
			return true
		}
		b := true
		for _, p := range preds {
			b = (b && p(v))
		}
		return b
	}
}
