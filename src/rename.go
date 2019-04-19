package rename

import (
	"flag"
	"fmt"
	"github.com/rwtodd/Go.Sed/sed"
	"os"
	"path"
	"strings"
)

type Args struct {
	Files   []string
	Replace string
	DryRun  bool
	Verbose bool
}

type FromTo struct {
	From string
	To   string
}

func ParseArgs() *Args {
	verbosePtr := flag.Bool("v", false, "Verbose")
	dryPtr := flag.Bool("d", false, "Dry run")

	flag.Parse()

	l := flag.NArg()
	files := flag.Args()[1 : l-1]
	replace := flag.Args()[l-1]
	return &Args{
		Files:   files,
		Replace: replace,
		DryRun:  *dryPtr,
		Verbose: *verbosePtr,
	}
}

func GetReplacements(args *Args) ([]FromTo, error) {
	engine, err := sed.New(strings.NewReader(args.Replace))
	if err != nil {
		return nil, err
	}
	var replacements []FromTo
	for _, file := range args.Files {
		dir := path.Dir(file)
		filename := path.Base(file)
		output, err := engine.RunString(filename)
		if err != nil {
			return nil, err
		}
		// go-sed always returns trailing newlines for some reason.
		output = strings.TrimSuffix(output, "\n")
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
		if args.Verbose || args.DryRun {
			fmt.Printf("%s\t-> %s\n", fromto.From, path.Base(fromto.To))
		}
		if !args.DryRun {
			os.Rename(fromto.From, fromto.To)
		}
	}

	return nil
}
