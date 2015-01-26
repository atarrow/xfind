// xfind project  options.go
// Copyright 2015 atarrow. All rights reserved.package search
package search

import (
	"errors"
	"fmt"
)

// nameOpts represents states & predicates that Finder has.
// In many cases, it is used for creating a new Finder.
type Options struct {
	Key         string // keyword
	Path        string // path
	IsRecursive bool   // IsRecursive
	FileType           // file type
}

func (e *Options) ToPredicates() Predicates {
	p := make(Predicates, 0, 4)
	if e.IsFileOnly() {
		p = append(p, isFileOnly)
	}
	if e.IsDirOnly() {
		p = append(p, isDirOnly)
	}
	if !e.IncludesDot() {
		p = append(p, isNotDotFile)
	}
	if len(e.Key) > 0 {
		pattern, err := PredicateWithPattern(e.Key)
		if err == nil {
			p = append(p, pattern)
		}
	}
	return p
}

func (e *Options) ToPredicate() Predicate {
	p := e.ToPredicates()
	return p.ANDPredicate()
}

func (e Options) String() string {
	const tmplate = `
    Keyword      %v
    Path         %v
    Recursive    %v
    Directories  %v
    Files        %v
    Invisibles   %v
	`
	return fmt.Sprintf(
		tmplate,
		e.Key,
		e.Path,
		e.IsRecursive,
		e.IncludesDir(),
		e.IncludesFile(),
		e.IncludesDot(),
	)
}

type FileType struct {
	omitsFile   bool
	omitsDir    bool
	includesDot bool
}

func (t *FileType) String() string {
	const template = `
	Directories  %v
	Files        %v
	Invisibles   %v
	`
	return fmt.Sprintf(
		template,
		t.IncludesDir(),
		t.IncludesFile(),
		t.IncludesDot(),
	)
}

func (t *FileType) Set(value string) error {
	if len(value) == 0 {
		return errors.New("No Value")
	}
	set := map[rune]bool{
		'a': false,
		'd': false, 'f': false}
	for _, r := range value {
		if _, ok := set[r]; ok {
			set[r] = true
		}
	}
	if set['a'] {
		t.includesDot = true
	}
	switch {
	case set['d'] && !set['f']:
		t.omitsFile = true
	case !set['d'] && set['f']:
		t.omitsDir = true
	}
	return nil
}

func (t FileType) IsFileOnly() bool   { return t.omitsDir && !t.omitsFile }
func (t FileType) IsDirOnly() bool    { return !t.omitsDir && t.omitsFile }
func (t FileType) IncludesFile() bool { return !t.omitsFile }
func (t FileType) IncludesDir() bool  { return !t.omitsDir }
func (t FileType) IncludesDot() bool  { return t.includesDot }
