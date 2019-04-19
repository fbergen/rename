package main

import (
	"fmt"
	"github.com/fbergen/rename/src"
)

func main() {
	args := rename.ParseArgs()
	err := rename.Run(args)
	if err != nil {
		fmt.Println(err)
	}

	// os.Rename()
}
