package args

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Parser struct {
	name        string
	description string
	flags       map[string]*flagDef
	options     map[string]*optionDef
	commands    map[string]*Parser
	parent      *Parser
	handler     func(*Result) error
	positionals []string
}

type flagDef struct {
	name, short, help string
}

type optionDef struct {
	name, short, help string
	choices           []string
}

type Result struct {
	Command     []string
	Flags       map[string]bool
	Options     map[string]string
	Positionals []string
	Args        []string
}

func New(name string) *Parser {
	return &Parser{
		name:     name,
		flags:    make(map[string]*flagDef),
		options:  make(map[string]*optionDef),
		commands: make(map[string]*Parser),
	}
}

func (p *Parser) Flag(name, short, help string) *Parser {
	p.flags[name] = &flagDef{name: name, short: short, help: help}
	return p
}

func (p *Parser) Option(name, short, help string, choices ...string) *Parser {
	p.options[name] = &optionDef{
		name: name, short: short, help: help, choices: choices,
	}
	return p
}

func (p *Parser) Command(name, help string) *Parser {
	child := New(name)
	child.description = help
	child.parent = p
	p.commands[name] = child
	return child
}

func (p *Parser) Action(fn func(*Result) error) *Parser {
	p.handler = fn
	return p
}

func (p *Parser) Positional(name string) *Parser {
	p.positionals = append(p.positionals, name)
	return p
}

func (p *Parser) Parse(args []string) (*Result, error) {
	if len(args) == 0 {
		return nil, errors.New("no arguments provided")
	}

	result := &Result{
		Flags:   make(map[string]bool),
		Options: make(map[string]string),
	}

	current := p
	i := 0

	for i < len(args) {
		arg := args[i]

		if cmd, exists := current.commands[arg]; exists {
			result.Command = append(result.Command, arg)
			current = cmd
			i++
			continue
		}

		if strings.HasPrefix(arg, "-") {
			if handled, consumed := current.parseFlag(arg, result); handled {
				i += consumed
				continue
			}
			if handled, consumed := current.parseOption(arg, args[i:], result); handled {
				i += consumed
				continue
			}
			return nil, fmt.Errorf("unknown flag: %s", arg)
		}

		result.Positionals = append(result.Positionals, arg)
		i++
	}

	result.Args = result.Positionals

	return result, nil
}

func (p *Parser) parseFlag(arg string, result *Result) (bool, int) {
	for _, flag := range p.flags {
		if arg == "--"+flag.name || (flag.short != "" && arg == "-"+flag.short) {
			result.Flags[flag.name] = true
			return true, 1
		}
	}
	return false, 0
}

func (p *Parser) parseOption(arg string, args []string, result *Result) (bool, int) {
	for _, opt := range p.options {
		if strings.HasPrefix(arg, "--"+opt.name+"=") {
			value := arg[len("--"+opt.name+"="):]
			result.Options[opt.name] = value
			return true, 1
		}

		if opt.short != "" && strings.HasPrefix(arg, "-"+opt.short+"=") {
			value := arg[len("-"+opt.short+"="):]
			result.Options[opt.name] = value
			return true, 1
		}

		if arg == "--"+opt.name || (opt.short != "" && arg == "-"+opt.short) {
			if len(args) < 2 {
				return false, 0
			}
			result.Options[opt.name] = args[1]
			return true, 2
		}
	}
	return false, 0
}

func (p *Parser) Run(args []string) error {
	result, err := p.Parse(args)
	if err != nil {
		return err
	}

	current := p
	for _, cmd := range result.Command {
		current = current.commands[cmd]
	}

	if current.handler != nil {
		return current.handler(result)
	}

	current.Usage()
	return nil
}

func (p *Parser) Usage() {
	fmt.Printf("Usage: %s", p.getPath())

	if len(p.commands) > 0 {
		fmt.Print(" <command>")
	}
	if len(p.flags) > 0 || len(p.options) > 0 {
		fmt.Print(" [options]")
	}
	if len(p.positionals) > 0 {
		for _, pos := range p.positionals {
			fmt.Printf(" <%s>", pos)
		}
	}
	fmt.Println()

	if p.description != "" {
		fmt.Printf("\n%s\n", p.description)
	}

	if len(p.commands) > 0 {
		fmt.Println("\nCommands:")
		for name, cmd := range p.commands {
			desc := cmd.description
			if desc == "" {
				desc = "-"
			}
			fmt.Printf("  %-15s %s\n", name, desc)
		}
	}

	flags := p.collectFlags()
	if len(flags) > 0 {
		fmt.Println("\nFlags:")
		for _, flag := range flags {
			short := ""
			if flag.short != "" {
				short = fmt.Sprintf("-%s, ", flag.short)
			}
			help := flag.help
			if help == "" {
				help = "-"
			}
			fmt.Printf("  %s--%s%s %s\n", short, flag.name, strings.Repeat(" ", 15-len(short+flag.name)), help)
		}
	}

	options := p.collectOptions()
	if len(options) > 0 {
		fmt.Println("\nOptions:")
		for _, opt := range options {
			short := ""
			if opt.short != "" {
				short = fmt.Sprintf("-%s, ", opt.short)
			}
			choices := ""
			if len(opt.choices) > 0 {
				choices = fmt.Sprintf(" {%s}", strings.Join(opt.choices, "|"))
			}
			help := opt.help
			if help == "" {
				help = "-"
			}
			fmt.Printf("  %s--%s <value>%s%s %s\n", short, opt.name, choices,
				strings.Repeat(" ", max(0, 10-len(short+opt.name+choices))), help)
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (p *Parser) getPath() string {
	var parts []string
	current := p
	for current != nil {
		if current.name != "" {
			parts = append([]string{current.name}, parts...)
		}
		current = current.parent
	}
	if len(parts) == 0 {
		return os.Args[0]
	}
	return strings.Join(parts, " ")
}

func (p *Parser) collectFlags() []*flagDef {
	var flags []*flagDef
	seen := make(map[string]bool)

	current := p
	for current != nil {
		for _, flag := range current.flags {
			if !seen[flag.name] {
				flags = append(flags, flag)
				seen[flag.name] = true
			}
		}
		current = current.parent
	}
	return flags
}

func (p *Parser) collectOptions() []*optionDef {
	var options []*optionDef
	seen := make(map[string]bool)

	current := p
	for current != nil {
		for _, opt := range current.options {
			if !seen[opt.name] {
				options = append(options, opt)
				seen[opt.name] = true
			}
		}
		current = current.parent
	}
	return options
}

func (r *Result) Flag(name string) bool {
	val, ok := r.Flags[name]
	if !ok {
		return false
	}
	return val
}
