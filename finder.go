// xfind project finder.go
// Copyright 2015 atarrow. All rights reserved.

package main

import (
	"fmt"
	"golang.org/x/text/unicode/norm"
	"io/ioutil"
	"os"
)

// Finder searches names of file or directory.
// Finder.Dir represents a working directory.
// Finder.predicates define logical conditions used to constrain a search.
// All the predicates in Finder.predicates works as conjunction.
type Finder struct {
	FinderEnv
	predicates []Predicate // predicates to use for search
}

// NewFinder() creates Finder with working directory specified by dirname.
func NewFinder(env FinderEnv) *Finder {
	if len(env.Dir) == 0 {
		env.Dir, _ = os.Getwd()
	}
	return &Finder{env, make([]Predicate, 0, 10)}
}

func (f *Finder) PrintResults() {
	list, dirs, err := f.Results()
	f.Print(list, err)
	if f.IsRecursive {
		recurse(dirs)
	}
}

func (f *Finder) Print(a []string, err error) {
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if i := len(a); i == 0 {
		return
	} else {
		s := format(a, "  ・", "\n")
		fmt.Printf(" %d entries found: %s", i, f.Dir+"/")
		if len(s) > 0 {
			fmt.Println(s)
		}
	}
}

func (f *Finder) Results() (list, dirs []string, err error) {
	fis, err := ioutil.ReadDir(f.Dir)
	if err != nil {
		return nil, nil, err
	}
	list = make([]string, 0, len(fis))
	dirs = make([]string, 0, len(fis))
	for _, v := range fis {
		if v.IsDir() {
			dirs = append(dirs, f.Dir+"/"+v.Name())
		}
		if f.MeetsPreds(v) {
			name := v.Name()
			if v.IsDir() {
				name += "/"
			}
			list = append(list, name)
		}
	}
	return list, dirs, nil
}

// MeetsPreds returns wheter specified v meets all Finder.predicates.
func (f *Finder) MeetsPreds(v os.FileInfo) bool {
	if len(f.predicates) == 0 {
		return true
	}
	var b bool = true
	for _, p := range f.predicates {
		b = (b && p(v))
	}
	return b
}

func recurse(dirs []string) {
	f := DefaultFinder()
	for _, path := range dirs {
		f.Dir = path
		list, subdirs, err := f.Results()
		f.Print(list, err)
		if len(subdirs) > 0 {
			recurse(subdirs)
		}
	}
}

// FinderEnv represents states & predicates that Finder has.
// In many cases, it is used for creating a new Finder.
type FinderEnv struct {
	IsDirOnly       bool
	IsFileOnly      bool
	IncludesDotFile bool
	IsRecursive     bool
	Keyword         string
	Dir             string
}

// NFDString returns a string converted into Unicode NFD normalized string.
func NFDString(s string) string {
	buf := []byte(s)
	buf = norm.NFD.Bytes(buf)
	return string(buf)
}

func format(a []string, pre, pos string) string {
	if len(a) == 0 {
		return ""
	}
	if len(a) == 1 {
		return "\n" + pre + a[0] + pos
	}
	n := len(pre) * len(a)
	n += len(pos) * len(a)
	for i := 0; i < len(a); i++ {
		n += len(a[i])
	}
	b := make([]byte, n+1)
	bp := copy(b, "\n")
	for _, s := range a {
		bp += copy(b[bp:], pre)
		bp += copy(b[bp:], s)
		bp += copy(b[bp:], pos)
	}
	return string(b)
}
