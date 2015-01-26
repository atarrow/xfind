// xfind project path.go
// Copyright 2015 atarrow. All rights reserved.
package search

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"golang.org/x/text/unicode/norm"
)

// NFDString returns a string converted into Unicode NFD normalized string.
func NFDString(s string) string {
	buf := []byte(s)
	buf = norm.NFD.Bytes(buf)
	return string(buf)
}

type Path string

func (p Path) Exists() bool {
	path := p.AbsPath()
	return exists(path)
}

func (p Path) AbsPath() string {
	if len(p) == 0 {
		return ""
	} else {
		return expand(string(p))
	}
}

func exists(name string) bool {
	_, err := os.Stat(name)
	if err != nil {
		return false
	} else {
		return true
	}
}

// expand returns tilde-expanded path.
func expand(name string) string {
	if beginsWithTilde(name) {
		fmt.Println("begins with tilde")
		usr, _ := user.Current()
		name = path.Join(usr.HomeDir, name[1:])
	}
	return name
}

// HasPrefix tests whether the string s begins with prefix.
func beginsWithTilde(s string) bool {
	tilde := "~/"
	return len(s) >= len(tilde) && s[0:len(tilde)] == tilde
}
