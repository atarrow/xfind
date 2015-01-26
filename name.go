// xfind project  name.go
// Copyright 2015 atarrow. All rights reserved.
package main

import (
	"fmt"
	"xfind/search"
)

func runName(opts *search.Options) error {
	fmt.Println(opts)
	return nil
}

func nameCmd() *Command {
	cmd := &Command{}
	cmd.Run = func(args []string) error {
		search.Parse(args)
		search.Run()
		return nil
	}
	return cmd
}
