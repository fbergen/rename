package rename

import (
	"fmt"
	"github.com/fatih/color"
	"regexp"
	"strings"
)

// Pre-compiled instructions to run
type instruction interface {
	Run(src string) (string, error)
	Highlight(src string) (string, string, error)
}

// instruction, for now chaining expressions is not supported
type Engine struct {
	ins instruction
}

type substitution struct {
	pattern *regexp.Regexp
	repl    string
	global  bool
}

func (s *substitution) Highlight(src string) (string, string, error) {
	first := true
	srcHighlight := s.pattern.ReplaceAllStringFunc(src, func(match string) string {
		if !first && !s.global {
			return match
		}
		first = false
		return s.pattern.ReplaceAllString(match, color.RedString("${0}"))
	})

	first = true
	destHighLight := s.pattern.ReplaceAllStringFunc(src, func(match string) string {
		if !first && !s.global {
			return match
		}
		first = false
		return s.pattern.ReplaceAllString(match, color.GreenString(s.repl))
	})

	return srcHighlight, destHighLight, nil
}

func (s *substitution) Run(src string) (string, error) {
	first := true
	return s.pattern.ReplaceAllStringFunc(src, func(match string) string {
		if !first && !s.global {
			return match
		}
		first = false
		return s.pattern.ReplaceAllString(match, s.repl)
	}), nil
}

func NewEngine(expression string) (*Engine, error) {
	ins, err := parse(expression)
	return &Engine{ins: ins}, err
}

func (e *Engine) Run(src string) (string, error) {
	return e.ins.Run(src)
}
func (e *Engine) Highlight(src string) (string, string, error) {
	return e.ins.Highlight(src)
}

func parse(expression string) (instruction, error) {
	//TODO(fbergen): Figure out the separator character.
	parts := strings.Split(expression, "/")
	if len(parts) < 3 {
		return nil, fmt.Errorf("Invalid expression '%s'", expression)
	}
	switch parts[0] {
	case "s":
		l := len(parts)
		if l < 4 {
			return nil, fmt.Errorf("Unterminated substitution command, required format 's/from/to/'")
		}
		var flags []rune
		if l > 3 {
			flags = []rune(parts[3])
		}
		subs, err := newSubstitution(parts[1], parts[2], flags)
		if err != nil {
			return nil, err
		}

		return subs, nil
	default:
		return nil, fmt.Errorf("unrecognized command '%s'", parts[0])
	}
}

func newSubstitution(pattern, replacement string, flags []rune) (instruction, error) {
	subs := &substitution{repl: replacement}
	caseInsensitive := false
	var err error
	for _, char := range flags {
		switch char {
		case 'g':
			subs.global = true
		case 'i':
			caseInsensitive = true
		default:
			err = fmt.Errorf("Unrecognized substitution flag '%v'", string(char))
		}
		if err != nil {
			return nil, err
		}
	}

	if caseInsensitive {
		pattern = "(?i)" + pattern
	}
	subs.pattern, err = regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return subs, err
}
