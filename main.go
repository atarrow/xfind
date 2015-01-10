// xfind project main.go
// Copyright 2015 atarrow. All rights reserved.

package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

var (
	showsVersion bool
	dirOnly      bool
	fileOnly     bool
	withDotFiles bool
	recursive    bool

	keyword string
	dirname string
	version string = "0.8.5"
)
var sharedEnv FinderEnv

func DefaultEnv() FinderEnv {
	return FinderEnv{
		IsDirOnly:       dirOnly,
		IsFileOnly:      fileOnly,
		IncludesDotFile: withDotFiles,
		IsRecursive:     recursive,
		Keyword:         keyword,
		Dir:             dirname,
	}
}

func DefaultPredicates(e FinderEnv) []Predicate {
	p := make([]Predicate, 0, 4)
	if e.IsFileOnly {
		p = append(p, isFile)
	}
	if e.IsDirOnly {
		p = append(p, isDir)
	}
	if !e.IncludesDotFile {
		p = append(p, isNotDotFile)
	}
	if len(e.Keyword) > 0 {
		pattern, err := PredicateWithPattern(keyword)
		if err != nil {

		} else {
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

func init() {
	flag.StringVar(&dirname, "p", "", "specifies directory. Default is a current directory.")
	flag.StringVar(&keyword, "k", "", "specifies keyword. Default represents all entries.")

	flag.BoolVar(&fileOnly, "f", false, "shows files only.")
	flag.BoolVar(&dirOnly, "d", false, "shows directories only.")
	flag.BoolVar(&withDotFiles, "a", false, "shows invisible files, too.")
	flag.BoolVar(&showsVersion, "v", false, "shows version if specified.")
	flag.BoolVar(&recursive, "r", false, "tries to search in subdirectories.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Printf("\n  Exsample: %s -k=.zip -p=/absolute/dir/name -f\n\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()
	if len(dirname) == 0 {
		dirname, _ = os.Getwd()
	}

	sharedEnv = DefaultEnv()
}

func main() {
	if showsVersion {
		fmt.Println(path.Base(os.Args[0] + " " + "ver." + version))
		return
	}
	printForDebug()
	finder := DefaultFinder()
	finder.PrintResults()
}

func printForDebug() {
	fmt.Println("jfind command debug---")
	fmt.Println("  ", "Version    :", version)
	fmt.Println("  ", "Keyword    :", keyword)
	fmt.Println("  ", "Dirname    :", dirname)
	fmt.Println("  ", "DirOnly    :", dirOnly)
	fmt.Println("  ", "FileOnly   :", fileOnly)
	fmt.Println("  ", "Recusive   :", recursive)
	fmt.Println("  ", "VisibleOnly:", withDotFiles)
	fmt.Println()
}
