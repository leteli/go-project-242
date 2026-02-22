package main

import (
	"code/files"
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	(&cli.Command{
		Name:  "hexlet-path-size",
		Usage: "print size of a file or directory",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			path := cmd.Args().Get(0)
			if path == "" {
				return fmt.Errorf("No arguments provided: please specify a path")
			}
			sizeBytes, err := files.GetSize(path)

			if err != nil {
				return fmt.Errorf("Failed to get size: %w", err)
			}
			fmt.Printf("%dB\t%s\n", sizeBytes, path)
			return nil
		},
	}).Run(context.Background(), os.Args)
}
