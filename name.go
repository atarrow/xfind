// xfind project  name.go
// Copyright 2015 atarrow. All rights reserved.
package main

import (
	"flag"
	"fmt"
	"xfind/search"
)

func nameCmd() *Command {
	cmd := &Command{
		Flag:  flag.NewFlagSet("xfind name", flag.ExitOnError),
		Short: "xfind name [-r][-v][-p path][-k kewword]",
		Long:  nameHelp}
	cmd.Flag.Usage = cmd.Usage
	cmd.Run = func(args []string) error {
		opt := nameopt(cmd.Flag)
		cmd.Flag.Parse(args)
		sch := search.New(*opt)
		fmt.Println(sch.Options.String())
		sch.Run()
		return nil
	}
	return cmd
}

func nameopt(f *flag.FlagSet) *search.Options {
	o := &search.Options{}
	f.StringVar(&o.Key, "k", "", "")
	f.StringVar(&o.Path, "p", search.Getwd(), "")
	f.BoolVar(&o.IsRecursive, "r", false, "")
	f.Var(&o.FileType, "t", "")
	return o
}

const nameHelp = `

  -v : shows version.
  -p : specifies path.
  -r : specifies wheter or not to search in subdirectories.
  -k : specifies keyword which means a file name.
  -t : specifies file types.
      [a]  adds invisible files and directories to targets.
      [d]  adds directories to targets.
      [f]  adds files to targets.
`
