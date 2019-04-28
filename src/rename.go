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
	verbosePtr := flag.BoolP("verbose", "v", false, "Show which files where renamed, if any.")
	noActPtr := flag.BoolP("no-action", "n", false, "Don't perform any changes. Show what files would have been renamed.")
	helpPtr := flag.BoolP("help", "h", false, "Show help dialog")
	interactivePtr := flag.BoolP("interactive", "i", false, "Interactive mode")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Usage: rename [options] files... expression\n\n")
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	l := flag.NArg()
	if l < 1 || *helpPtr {
		flag.Usage()
		os.Exit(2)
	}

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

func GetReplacements(engine *Engine, args *Args) ([]FromTo, error) {
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

func PrintRename(engine *Engine, fromto FromTo) {
	// color match.
	from, to, _ := engine.Highlight(fromto.From)
	fmt.Printf("%s\t-> %s\n", from, path.Base(to))
}

func Run(args *Args) error {
	engine, err := NewEngine(args.Expression)
	if err != nil {
		return err
	}
	replacements, err := GetReplacements(engine, args)
	if err != nil {
		return err
	}
	for _, fromto := range replacements {
		if args.Verbose || args.NoAct {
			PrintRename(engine, fromto)
		}
		if !args.NoAct {
			os.Rename(fromto.From, fromto.To)
		}
	}

	return nil
}
