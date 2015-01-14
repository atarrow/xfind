// xfind project default.go
// Copyright 2015 atarrow. All rights reserved.

package main

import ()

func DefaultPredicates(e FinderEnv) []Predicate {
	p := make([]Predicate, 0, 4)
	if e.IsFileOnly() {
		p = append(p, isFileOnly)
	}
	if e.IsDirOnly() {
		p = append(p, isDirOnly)
	}
	if e.OmitsDot() {
		p = append(p, isNotDotFile)
	}
	if len(e.Keyword) > 0 {
		pattern, err := PredicateWithPattern(e.Keyword)
		if err == nil {
			p = append(p, pattern)
		}
	}
	return p
}

func DefaultFinder() *Finder {
	f := NewFinder(sharedEnv)
	f.predicates = append(f.predicates, DefaultPredicates(sharedEnv)...)
	return f
}
