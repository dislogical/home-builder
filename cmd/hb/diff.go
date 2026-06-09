package main

import (
	"context"
	"log/slog"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/urfave/cli/v3"
)

var diff *cli.Command = &cli.Command{
	Name: "diff",
	Action: func(ctx context.Context, c *cli.Command) error {
		dmp := diffmatchpatch.New()

		diff := dmp.DiffMain("asdf", "hjkl;", true)
		slog.InfoContext(ctx, dmp.DiffPrettyText(diff))

		return nil
	},
}
