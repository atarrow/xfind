// xfind project default.go
// Copyright 2015 atarrow. All rights reserved.

package search

import (
	"flag"
	"fmt"
	"os"
)

var (
	CommandLine *Search
	Flag        *flag.FlagSet
)

const help = `

  -v : shows version.
  -p : specifies path.
  -r : specifies wheter or not to search in subdirectories.
  -k : specifies keyword which means a file name.
  -t : specifies file types.
      [a]  adds invisible files and directories to targets.
      [d]  adds directories to targets.
      [f]  adds files to targets.
`

func Parse(args []string) {
	var opts *Options
	opts, Flag = ParsedOptions(args)
	CommandLine = &Search{opts, nil}
}

func ParsedOptions(args []string) (*Options, *flag.FlagSet) {
	var (
		f = flag.NewFlagSet("search", flag.ExitOnError)
		o = &Options{}
	)

	f.StringVar(&o.Key, "k", "", "")
	f.StringVar(&o.Path, "p", "", "")
	f.BoolVar(&o.IsRecursive, "r", false, "")
	f.Var(&o.FileType, "t", "")
	f.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s", "name [-r][-v][-p path][-k kewword]")
		fmt.Fprintf(os.Stderr, "%s\n", help)
		os.Exit(2)

	}

	f.Parse(args)

	if len(o.Path) == 0 {
		o.Path, _ = os.Getwd()
	}

	return o, f
}

func Run() {
	fmt.Println(CommandLine.Options.String())
	CommandLine.Run()
}
