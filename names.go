// xfind project names.go
// Copyright 2015 atarrow. All rights reserved.

package main

import (
	"golang.org/x/text/unicode/norm"
	"io/ioutil"
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
	isDir        Predicate = func(v os.FileInfo) bool { return v.IsDir() }
	isNotDir     Predicate = func(v os.FileInfo) bool { return !v.IsDir() }
)

// PredicateWithPattern creates a new predicate that represents
// whether an argument matches a regular expression pattern.
func PredicateWithPattern(pattern string) (Predicate, error) {
	re, err := regexp.Compile(NFDString(pattern))
	if err != nil {
		return nil, err
	}
	pred := func(v os.FileInfo) bool { return re.MatchString(v.Name()) }
	return pred, nil
}

// Finder searches names of file or directory.
// Finder.dirname represents a working directory.
// Finder.predicates define logical conditions used to constrain a search.
// All the predicates in Finder.predicates works as conjunction.
type Finder struct {
	dirname    string      // directory to search entries in
	predicates []Predicate // predicates to use for search
}

// NewFinder() creates Finder with working directory specified by dirname.
func NewFinder(dirname string) Finder {
	finder := Finder{dirname, make([]Predicate, 0, 10)}
	return finder
}

// Append() adds specified predicates to the end of Finder.predicates.
func (finder *Finder) Append(p ...Predicate) {
	finder.predicates = append(finder.predicates, p...)
}

// Find() searches with specified predicates & working directory.
// It returns the result as []string & error.
func (finder Finder) Find() ([]string, error) {
	fileInfos, err := ioutil.ReadDir(finder.dirname)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(fileInfos))
	for _, v := range fileInfos {
		if finder.MeetsPreds(v) {
			names = append(names, distinctName(v))
		}
	}
	return names, nil
}

// MeetsPreds returns wheter specified v meets all Finder.predicates.
func (finder Finder) MeetsPreds(v os.FileInfo) bool {
	if len(finder.predicates) == 0 {
		return true
	}
	var result bool = true
	for _, pred := range finder.predicates {
		result = (result && pred(v))
	}
	return result
}

// distinctName returns FileInfo.Name() appended slash to if "fi" represents directory.
func distinctName(v os.FileInfo) string {
	name := v.Name()
	if v.IsDir() {
		name += "/"
	}
	return name
}

// NFDString returns a string converted into Unicode NFD normalized string.
func NFDString(s string) string {
	buf := []byte(s)
	buf = norm.NFD.Bytes(buf)
	return string(buf)
}
