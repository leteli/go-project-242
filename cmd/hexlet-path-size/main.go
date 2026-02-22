package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	(&cli.Command{
		Name:  "hexlet-path-size",
		Usage: "print size of a file or directory",
	}).Run(context.Background(), os.Args)
}
