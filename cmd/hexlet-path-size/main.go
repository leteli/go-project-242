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
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "human",
				Aliases: []string{"H"},
				Value:   false,
				Usage:   "human-readable sizes (auto-select unit)",
			},
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Value:   false,
				Usage:   "include hidden files and directories",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			path := cmd.Args().Get(0)
			if path == "" {
				return fmt.Errorf("no arguments provided: please specify a path")
			}
			sizeBytes, err := files.GetSize(path, cmd.Bool("all"))

			if err != nil {
				return fmt.Errorf("failed to get size: %w", err)
			}
			res, err := files.FormatSize(sizeBytes, cmd.Bool("human"))
			if err != nil {
				return fmt.Errorf("failed to format size %w", err)
			}
			fmt.Printf("%s\t%s\n", res, path)
			return nil
		},
	}
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
