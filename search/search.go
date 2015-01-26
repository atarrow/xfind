// xfind project finder.go
// Copyright 2015 atarrow. All rights reserved.

package search

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
)

type Search struct {
	*Options
	Predicates // predicates to use for search
}

// New() creates Search with working directory specified by dirname.
func New(opts Options) *Search {
	return &Search{
		Options:    &opts,
		Predicates: make(Predicates, 0, 10),
	}
}

func (s *Search) Run() {
	r := s.result()
	if r.Error != nil {
		log.Fatal(r.Error.Error())
	}
	a := r.Subpaths
	r.Print()
	if s.IsRecursive && len(a) > 0 {
		search := New(*s.Options)
		for _, path := range a {
			search.Path = path
			search.Run()
		}
	}
}

func (f *Search) result() Result {
	fis, err := ioutil.ReadDir(f.Path)
	if err != nil {
		return Result{
			Path:  f.Path,
			Error: err,
		}
	}
	list := make([]string, 0, len(fis))
	dirs := make([]string, 0, len(fis))

	for _, v := range fis {
		if v.IsDir() {
			dirs = append(dirs, path.Join(f.Path, v.Name()))
		}
		p := f.ToPredicate()
		if p(v) {
			name := v.Name()
			if v.IsDir() {
				name += "/"
			}
			list = append(list, name)
		}
	}
	return Result{
		Path:     f.Path,
		List:     list,
		Subpaths: dirs,
	}
}

func (f *Search) toPred() Predicate {
	x := f.Options.ToPredicates()
	y := f.Predicates
	n := len(x) + len(y)
	p := make(Predicates, 0, n)
	if n == 0 {
		return AlwaysTrue
	}
	p = append(p, x...)
	p = append(p, y...)
	return p.ANDPredicate()
}

type Result struct {
	Path     string
	List     []string
	Subpaths []string
	Error    error
}

func (r Result) Print() {
	if r.Error != nil {
		log.Fatal(r.Error.Error())
	}
	if i := len(r.List); i == 0 {
		return
	} else {
		s := format(r.List, "  ・", "\n")
		fmt.Printf(" %d entries found: %s", i, r.Path+"/")
		if len(s) > 0 {
			fmt.Println(s)
		}
	}
}

func format(a []string, pre, pos string) string {
	if len(a) == 0 {
		return ""
	}
	if len(a) == 1 {
		return "\n" + pre + a[0] + pos
	}
	n := len(pre) * len(a)
	n += len(pos) * len(a)
	for i := 0; i < len(a); i++ {
		n += len(a[i])
	}
	b := make([]byte, n+1)
	bp := copy(b, "\n")
	for _, s := range a {
		bp += copy(b[bp:], pre)
		bp += copy(b[bp:], s)
		bp += copy(b[bp:], pos)
	}
	return string(b)
}
