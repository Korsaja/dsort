package main

import (
	"os"

	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slog"
)

func main() {
	os.Exit(run(os.Args, os.Stdout, os.Stderr))
}

func run(args []string, stdout *os.File, stderr *os.File) int {
	var logger = slog.New(slog.NewTextHandler(stderr))

	app := cli.App{
		Name:      "dir sorter",
		Usage:     "sorting files in the specified directory by date",
		UsageText: "dsort --dir ~/Download --remove=true --skip dir1,dir2",
		Suggest:   true,
		Writer:    stdout,
		ErrWriter: stderr,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dir",
				Usage:    "dir for sorting",
				Required: true,
				Aliases:  []string{"d"},
			},
			&cli.StringSliceFlag{
				Name:  "skip",
				Usage: "skip dirs",
			},
			&cli.BoolFlag{
				Name:    "remove",
				Aliases: []string{"r"},
				Usage:   "remove old file",
				Value:   false,
			},
		},
		Action: func(c *cli.Context) error {
			actions := SortAction{
				Removed:  c.Bool("remove"),
				DirPath:  c.String("dir"),
				SkipDirs: c.StringSlice("skip"),
			}
			return DoSort(actions, logger)
		},
	}

	if err := app.Run(args); err != nil {
		logger.Error("app run failed", err)
		return 127
	}

	return 0
}
