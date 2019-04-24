package rename

import (
	"bufio"
	"fmt"
	flag "github.com/ogier/pflag"
	"os"
	"path"
)

type Args struct {
	Files       []string
	Expression  string
	NoAct       bool
	Verbose     bool
	Interactive bool
	Force       bool
}

type FromTo struct {
	From string
	To   string
}

func ParseArgs() *Args {
	verbosePtr := flag.BoolP("verbose", "v", false, "Verbose")
	noActPtr := flag.BoolP("no-action", "n", false, "No-action")
	interactivePtr := flag.BoolP("interactive", "i", false, "Interactive")

	flag.Parse()

	l := flag.NArg()
	var files []string
	if l < 2 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			files = append(files, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return nil
		}
	} else {
		files = flag.Args()[0 : l-1]
	}

	expression := flag.Args()[l-1]
	return &Args{
		Files:       files,
		Expression:  expression,
		NoAct:       *noActPtr,
		Verbose:     *verbosePtr,
		Interactive: *interactivePtr,
	}
}

func GetReplacements(args *Args) ([]FromTo, error) {
	engine, err := NewEngine(args.Expression)
	if err != nil {
		return nil, err
	}
	var replacements []FromTo
	for _, file := range args.Files {
		dir := path.Dir(file)
		filename := path.Base(file)
		output, err := engine.Run(filename)
		if err != nil {
			return nil, err
		}
		replacements = append(replacements, FromTo{path.Join(dir, filename), path.Join(dir, output)})
	}
	return replacements, nil
}

func Run(args *Args) error {
	replacements, err := GetReplacements(args)
	if err != nil {
		return err
	}
	for _, fromto := range replacements {
		if args.Verbose || args.NoAct {
			fmt.Printf("%s\t-> %s\n", fromto.From, path.Base(fromto.To))
		}
		if !args.NoAct {
			os.Rename(fromto.From, fromto.To)
		}
	}

	return nil
}
