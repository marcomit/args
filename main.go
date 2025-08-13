package main

import (
	"args/parser"
	"os"
)

func main() {
	parser := args.New()

	dock := parser.Command("dock")

	req := parser.Command("req")

	req.Command("new")
	req.Command("run")

	dock.Command("init")
	dock.Command("use")
	dock.Command("list")
	dock.Command("status")

	parser.Parse(os.Args)
}
