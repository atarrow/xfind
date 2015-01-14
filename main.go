// xfind project main.go
// Copyright 2015 atarrow. All rights reserved.

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	showsVersion bool
	version      string = "0.8.6.b"
	sharedEnv    FinderEnv
)

// FinderEnv represents states & predicates that Finder has.
// In many cases, it is used for creating a new Finder.
type FinderEnv struct {
	IsRecursive bool
	Keyword     string
	Dir         string
	FileType
}

const help = `
Usage of xfind:

	xfind -k=.zip -p=/absolute/dir/name

  -v : shows version.
  -p : specifies path.
  -r : specifies wheter or not to search in subdirectories.
  -k : specifies keyword which means a file name.
  -t : specifies file types.
      [a]  adds invisible files and directories to targets.
      [d]  adds directories to targets.
      [f]  adds files to targets.
`

func init() {
	flag.StringVar(&sharedEnv.Dir, "p", "", "")
	flag.StringVar(&sharedEnv.Keyword, "k", "", "")
	flag.BoolVar(&sharedEnv.IsRecursive, "r", false, "")
	flag.Var(&sharedEnv.FileType, "t", "")
	flag.BoolVar(&showsVersion, "v", false, "")

	flag.Usage = func() {
		//path := path.Base(os.Args[0])
		fmt.Println(help)
	}
	if len(sharedEnv.Dir) == 0 {
		sharedEnv.Dir, _ = os.Getwd()
	}
	flag.Parse()
}

func main() {
	if showsVersion {
		fmt.Println(os.Args[0] + " " + "ver." + version)
		return
	}
	printSetting()
	finder := DefaultFinder()
	finder.PrintResults()

}

func printSetting() {
	fmt.Println("jfind current setting---")
	fmt.Println("  ", "Version    :", version)
	fmt.Println("  ", "Keyword    :", sharedEnv.Keyword)
	fmt.Println("  ", "Dirname    :", sharedEnv.Dir)
	fmt.Println("  ", "Recusive   :", sharedEnv.IsRecursive)
	fmt.Println(&(sharedEnv.FileType))
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
