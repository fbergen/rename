package rename

import (
	"bufio"
	"fmt"
	"github.com/manifoldco/promptui"
	flag "github.com/ogier/pflag"
	"io"
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
	Copy        bool
}

type FromTo struct {
	From string
	To   string
}

func ParseArgs() *Args {
	verbosePtr := flag.BoolP("verbose", "v", false, "Show which files where renamed, if any.")
	noActPtr := flag.BoolP("no-action", "n", false, "Don't perform any changes. Show what files would have been renamed.")
	forcePtr := flag.BoolP("force", "f", false, "Overwrite existing files.")
	copyPtr := flag.BoolP("copy", "c", false, "Copy instead of move.")
	helpPtr := flag.BoolP("help", "h", false, "Show help dialog.")
	interactivePtr := flag.BoolP("interactive", "i", false, "Ask for confirmation, before renaming")

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
		Copy:        *copyPtr,
		Force:       *forcePtr,
		Interactive: *interactivePtr,
	}
}

func GetReplacements(engine *Engine, args *Args) ([]FromTo, error) {
	destinations := make(map[string]bool)
	var replacements []FromTo
	for _, file := range args.Files {
		dir := path.Dir(file)
		filename := path.Base(file)
		dest, err := engine.Run(filename)
		if err != nil {
			return nil, err
		}
		dest = path.Join(dir, dest)
		if destinations[dest] {
			return nil, fmt.Errorf("Conflicting rename pattern, multiple files will be renamed to the same destination  '%s'", dest)
		}
		destinations[dest] = true

		replacements = append(replacements, FromTo{path.Join(dir, filename), dest})
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

	if args.Interactive || args.Verbose || args.NoAct {
		for _, fromto := range replacements {
			PrintRename(engine, fromto)
		}
	}
	prompt := promptui.Prompt{
		Label:     "Continue?",
		IsConfirm: true,
	}

	_, err = prompt.Run()
	if err != nil {
		return nil
	}

	for _, fromto := range replacements {
		if !args.NoAct {
			act := true
			if _, err := os.Stat(fromto.To); err == nil {
				// File exists
				if !args.Force {
					fmt.Printf("Not overwriting file: '%s'\n", fromto.To)
					act = false
				}
			}
			if act {
				if args.Copy {
					err := copy(fromto.From, fromto.To)
					if err != nil {
						fmt.Printf("Failed to copy file '%s'\n", err)
					}
				} else {
					os.Rename(fromto.From, fromto.To)
					if err != nil {
						fmt.Printf("Failed to rename file '%s'\n", err)
					}
				}
			}
		}
	}

	return nil
}

func copy(source, destination string) error {
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()

	dest, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer dest.Close()
	_, err = io.Copy(dest, src)
	return err
}
