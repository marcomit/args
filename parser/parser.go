package args

type ArgParser struct {
	name     string
	flags    []*Flag
	options  []*Option
	children []*ArgParser
	on       func(*ArgResult)
}

type ArgResult struct {
}

type Flag struct {
	Name string
	Abbr string
	Help string
}
type Option struct {
	Name string
	Abbr string
	Help string
}

type Command struct {
	Name string
	Abbr string
	Help string
}

func New(on func(*ArgResult)) *ArgParser {
	return &ArgParser{children: []*ArgParser{}, on: on}
}

func (parser *ArgParser) Command(name string) *ArgParser {
	child := &ArgParser{
		name:     name,
		children: []*ArgParser{},
	}
	parser.children = append(parser.children, child)
	return child
}

func (parser *ArgParser) Flag(flag Flag) *ArgParser {
	parser.flags = append(parser.flags, &flag)
	return parser
}
func (parser *ArgParser) Option(option Option) *ArgParser {
	parser.options = append(parser.options, &option)
	return parser
}

func (parser *ArgParser) checkFlag() bool {
	return false

}

func (parser *ArgParser) checkOption() bool {
	return false
}

func (parser *ArgParser) subcommand(name string) *ArgParser {
	for _, child := range parser.children {
		if child.name == name {
			return child
		}
	}
	return nil
}

func (parser *ArgParser) Parse(args []string) (*ArgResult, error) {
	// dummy := parser
	// for i := 0; i < len(args); i++ {
	// 	if child := parser.checkCommand(args[i]); child != nil {
	// 		dummy = child
	// 		dummy.on()
	// 	}
	// }
	return nil, nil
}
