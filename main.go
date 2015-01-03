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
	showsVersion bool
	dirOnly      bool
	fileOnly     bool
	withDotFiles bool

	keyword string
	dirname string
	version string = "0.8.3"
)
var env FinderEnv

func DefaultParameter() FinderEnv {
	return FinderEnv{
		dirOnly,
		fileOnly,
		withDotFiles,
		keyword,
		dirname,
	}
}

func init() {
	flag.StringVar(&dirname, "p", "", "specifies directory. Default is a current directory.")
	flag.StringVar(&keyword, "k", "", "specifies keyword. Default represents all entries.")

	flag.BoolVar(&fileOnly, "f", false, "shows files only.")
	flag.BoolVar(&dirOnly, "d", false, "shows directories only.")
	flag.BoolVar(&withDotFiles, "a", false, "shows invisible files, too.")
	flag.BoolVar(&showsVersion, "v", false, "shows version if specified.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Printf("\n  Exsample: %s -k=.zip -p=/absolute/dir/name -f\n\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()
	env = DefaultParameter()
}

func main() {
	if showsVersion {
		fmt.Println(path.Base(os.Args[0] + " " + "ver." + version))
		return
	}

	finder := NewFinderWith(env)
	printForDebug()
	names, err := finder.Find()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(strings.Join(names, "\n"))
}

func printForDebug() {
	fmt.Println("jfind command debug---")
	fmt.Println("  ", "Version    :", version)
	fmt.Println("  ", "Keyword    :", keyword)
	fmt.Println("  ", "Dirname    :", dirname)
	fmt.Println("  ", "DirOnly    :", dirOnly)
	fmt.Println("  ", "FileOnly   :", fileOnly)
	fmt.Println("  ", "VisibleOnly:", withDotFiles)
	fmt.Println()
}
