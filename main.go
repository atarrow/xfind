// xfind project main.go
// Copyright 2015 atarrow. All rights reserved.

package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

var (
	showsFileOnly  bool
	showsInvisible bool
	showsDirOnly   bool
	keyword        string
	dirname        string
	showsVersion   bool
	version        string = "0.8.1"
)

func init() {
	flag.StringVar(&dirname, "p", "", "specifies directory. Default is a current directory.")
	flag.StringVar(&keyword, "k", "", "specifies keyword. Default represents all entries.")

	flag.BoolVar(&showsFileOnly, "f", false, "shows files only.")
	flag.BoolVar(&showsDirOnly, "d", false, "shows directories only.")
	flag.BoolVar(&showsInvisible, "a", false, "shows invisible files, too.")
	flag.BoolVar(&showsVersion, "v", false, "shows version if specified.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Printf("\n  Exsample: %s -k=.zip -p=/absolute/dir/name -f\n\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()
}

func main() {
	if showsVersion {
		fmt.Println(path.Base(os.Args[0] + " " + "ver." + version))
		return
	}
	finder := NewFinder(dirname)
	finder.loadPredicates()
	names, err := finder.Find()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(strings.Join(names, "\n"))
}

func (finder *Finder) loadPredicates() {
	switch {
	case showsFileOnly && !showsDirOnly:
		finder.Append(isNotDir)
		fallthrough
	case !showsFileOnly && showsDirOnly:
		finder.Append(isDir)
		fallthrough
	case !showsInvisible:
		finder.Append(isNotDotFile)
		fallthrough
	case len(dirname) > 0:
		pwd, _ := os.Getwd()
		finder.dirname = pwd
		fallthrough
	case len(keyword) > 0:
		pattern, _ := PredicateWithPattern(keyword)
		finder.Append(pattern)
	}
}
