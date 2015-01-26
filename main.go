// xfind project main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const version = "0.8.8.a"

var mainV bool

func main() {

	commands := map[string]*Command{
		"name": nameCmd(),
	}

	flag.BoolVar(&mainV, "v", false, "shows version")
	flag.Parse()

	if mainV {
		fmt.Println("flagx Version " + version)
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Need subcommand.")
		os.Exit(1)
	}
	if cmd, ok := commands[args[0]]; !ok {
		log.Fatalf("Unknown command: %s", args[0])

	} else if err := cmd.Run(args[1:]); err != nil {
		log.Fatal(err)
	}
}

type Command struct {
	Run       func(args []string) error
	UsageLine string
	Short     string
	Long      string
	Flag      *flag.FlagSet
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", c.Long)
	os.Exit(2)
}
