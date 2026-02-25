package main

import (
	"code/files"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "hexlet-path-size",
		Usage: "print size of a file or directory",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			path := cmd.Args().Get(0)
			if path == "" {
				return fmt.Errorf("no arguments provided: please specify a path")
			}
			sizeBytes, err := files.GetSize(path)

			if err != nil {
				return fmt.Errorf("failed to get size: %w", err)
			}
			fmt.Printf("%dB\t%s\n", sizeBytes, path)
			return nil
		},
	}
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
