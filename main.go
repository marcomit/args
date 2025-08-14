package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println(os.Args[1:])
	parser := New()

	req := parser.Command(Command{
		Name: "req",
	})

	req.Command(Command{
		Name: "new",
	})

	req.Positional(Positional{
		Name: "request name",
	})

	dock := parser.Command(Command{
		Name: "dock",
	})

	run := dock.Command(Command{
		Name: "run",
	})

	run.Flag(Flag{
		Name: "output",
		Abbr: "o",
	})

	parser.Parse(os.Args[1:])
	fmt.Println()
	fmt.Println("Print")
	parser.Print(0)
}
