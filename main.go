// xfind project main.go
// Copyright 2015 atarrow. All rights reserved.

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

var (
	showsVersion bool
	recursive    bool
	fileType     FileType
	keyword      string
	dirname      string
	version      string = "0.8.5.h"
)

var sharedEnv FinderEnv

func DefaultEnv() FinderEnv {
	return FinderEnv{
		IsRecursive: recursive,
		Keyword:     keyword,
		Dir:         dirname,
		FileType:    fileType,
	}
}

func DefaultPredicates(e FinderEnv) []Predicate {
	p := make([]Predicate, 0, 4)
	if e.IsFileOnly() {
		p = append(p, isFile)
	}
	if e.IsDirOnly() {
		p = append(p, isDirOnly)
	}
	if e.OmitsDot() {
		p = append(p, isNotDotFile)
	}
	if len(e.Keyword) > 0 {
		pattern, err := PredicateWithPattern(keyword)
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

const msgFileTYpe = ` specifies file type:
		[a]  adds invisible files and directories to targets.
		[d]  adds directories to targets.
		[f]  adds files to targets.`

func init() {
	flag.StringVar(&dirname, "p", "", "specifies directory. Default is a current directory.")
	flag.StringVar(&keyword, "k", "", "specifies keyword. Default represents all entries.")
	flag.BoolVar(&showsVersion, "v", false, "shows version if specified.")
	flag.BoolVar(&recursive, "r", false, "tries to search in subdirectories.")
	flag.Var(&fileType, "t", msgFileTYpe)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Printf("\n  Exsample: %s -k=.zip -p=/absolute/dir/name -f\n\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	if len(dirname) == 0 {
		dirname, _ = os.Getwd()
	}
	flag.Parse()
	sharedEnv = DefaultEnv()
}

func main() {
	if showsVersion {
		fmt.Println(path.Base(os.Args[0] + " " + "ver." + version))
		return
	}
	printSetting()
	finder := DefaultFinder()
	finder.PrintResults()

}

func printSetting() {
	fmt.Println("jfind current setting---")
	fmt.Println("  ", "Version    :", version)
	fmt.Println("  ", "Keyword    :", keyword)
	fmt.Println("  ", "Dirname    :", dirname)
	fmt.Println("  ", "Recusive   :", recursive)
	fmt.Println(&fileType)
	fmt.Println()
}

type FileType struct {
	omitsFile   bool
	omitsDir    bool
	includesDot bool
}

func (t *FileType) String() string {
	var (
		m   = map[bool]string{true: "YES", false: "NO"}
		tab = "   "
	)
	a := []string{
		"",
		tab + "directories: " + m[!t.OmitsDir()],
		tab + "files      : " + m[!t.OmitsFile()],
		tab + "invisibles : " + m[!t.OmitsDot()],
	}
	return strings.Join(a, "\n")
}

func (t *FileType) Set(value string) error {
	if len(value) == 0 {
		return errors.New("No Value")
	}
	set := map[rune]bool{
		'a': false,
		'd': false,
		'f': false,
	}
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

func (t FileType) IsFileOnly() bool {
	return t.OmitsDir() && !t.OmitsFile()
}
func (t FileType) IsDirOnly() bool {
	return !t.OmitsDir() && t.OmitsFile()
}
func (t FileType) OmitsFile() bool { return t.omitsFile }
func (t FileType) OmitsDir() bool  { return t.omitsDir }
func (t FileType) OmitsDot() bool  { return !t.includesDot }
