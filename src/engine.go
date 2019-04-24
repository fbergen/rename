package rename

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Pre-compiled instructions to run
type instruction func(src string) (string, error)

// instruction, for now chaining expressions is not supported
type Engine struct {
	ins instruction
}

type substitution struct {
	pattern *regexp.Regexp
	repl    string
	global  bool
}

func (s *substitution) run(src string) (string, error) {
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
	return e.ins(src)
}

func parse(expression string) (instruction, error) {
	//TODO(fbergen): Figure out the separator character.
	parts := strings.Split(expression, "/")
	if len(parts) < 3 {
		return nil, errors.New("err")
	}
	switch parts[0] {
	case "s":
		var flags []rune
		if len(parts) > 3 {
			flags = []rune(parts[3])
		}
		subs, _ := newSubstitution(parts[1], parts[2], flags)
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
			err = fmt.Errorf("Bad regexp flag <%v>", char)
		}
		if err != nil {
			break
		}
	}

	if caseInsensitive {
		pattern = "(?i)" + pattern
	}
	subs.pattern, err = regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return subs.run, err
}
