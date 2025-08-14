package main

import (
	"fmt"
)

type Parser struct {
	command     *Command
	positionals []*Positional
	flags       []*Flag
	options     []*Option
	children    []*Parser
	selected    int
	parsed      bool
}

type Command struct {
	Name  string
	Descr string
	Abbr  string
	Help  string
}

type Positional struct {
	Name     string
	Descr    string
	Abbr     string
	Help     string
	value    string
	Validate func(string) bool
}

type Flag struct {
	Name   string
	Descr  string
	Abbr   string
	Help   string
	active bool
}

type Option struct {
	Name     string
	Descr    string
	Abbr     string
	Help     string
	Options  []string
	selected string
}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) Command(command Command) *Parser {
	child := New()

	child.command = &command

	p.children = append(p.children, child)

	return child
}

func (p *Parser) Positional(positional Positional) *Parser {
	p.positionals = append(p.positionals, &positional)
	return p
}

func (p *Parser) Flag(flag Flag) *Parser {
	p.flags = append(p.flags, &flag)
	return p
}
func (p *Parser) Option(option Option) *Parser {
	p.options = append(p.options, &option)
	return p
}

func (p *Parser) Print(depth int) {
	if p.command != nil {
		fmt.Println(p.command.Name, depth)
	}

	for _, positional := range p.positionals {

		fmt.Println(positional.Name)
	}

	for _, flag := range p.flags {

		fmt.Println(flag.Name)
	}

	for _, option := range p.options {

		fmt.Println(option.Name)
	}
	for i := range p.children {
		p.children[i].Print(depth + 1)
	}
}

func (p *Parser) checkCommand(name string) (*Parser, bool) {
	for _, child := range p.children {
		if child.command == nil {
			continue
		}

		if child.command.Name == name {
			p.children = []*Parser{child}
			return child, true
		}
	}
	return p, false
}

func (p *Parser) checkPositional(i *int, args []string) bool {
	//
	// for j, positional := range p.positionals {
	// 	positional.value = args[*i+j]
	// }
	return false
}

func (p *Parser) checkFlag(name string) bool {

	// for _, flag := range p.flags {
	// 	if flag.Name == "--"+name {
	// 		flag.active = true
	// 		return true
	// 	}
	// 	if flag.Abbr == "-"+name {
	// 		flag.active = true
	// 		return true
	// 	}
	// }
	//
	// p.flags = []*Flag{}
	return false
}

func (p *Parser) checkOption(args []string) bool {
	// if len(args) == 0 {
	// 	return false
	// }
	//
	// for _, option := range p.options {
	// 	if !strings.HasPrefix(args[0], "--"+option.Name) && !strings.HasPrefix(args[0], "-"+option.Abbr) {
	// 		continue
	// 	}
	// 	val := strings.Split(args[0], "=")
	// 	if len(val) == 1 {
	// 		if len(args) < 2 {
	// 			return false
	// 		}
	//
	// 		option.selected = args[1]
	// 		return true
	// 	} else if len(val) == 2 {
	// 		option.selected = args[1]
	// 		return true
	// 	}
	// }
	//
	// return false
	return false
}

func (p *Parser) Parse(args []string) {
	for i := 0; i < len(args); i++ {
		fmt.Println("CURR:", args[i])
		p, isCommand := p.checkCommand(args[i])

		if !isCommand {
			return
		}
		fmt.Println("Is command", p.command.Name)

		i++

		isPositional := p.checkPositional(&i, args[i:])

		isFlag := p.checkFlag(args[i])
		if !isFlag {
			p.checkOption(args[i:])
		}

		if !isFlag && !isPositional {
			i--
		}
	}
}
