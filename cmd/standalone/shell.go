package main

import (
	"fmt"
	"strings"

	minutes "github.com/42minutes/go-42minutes"
	ishell "github.com/abiosoft/ishell"
)

func (d *daemon) startShell() {
	// configure our interactive shell
	shell := ishell.New()
	shell.Println("* 42minutes standalone client")

	// register add function to add new series
	shell.Register("add", func(args ...string) (string, error) {
		tt := ""
		if len(args) > 0 {
			tt = strings.Join(args, " ")
		}
		shs, err := d.glibrary.QueryShowsByTitle(tt)
		if err != nil || len(shs) == 0 {
			return "Could not find series by name.", err
		}

		sh := shs[0]

		if err := d.ulibrary.UpsertShow(sh); err != nil {
			shell.Println()
			return "Could not find series by name.", err
		}

		return fmt.Sprintf("Added '%s' to your library.", sh.Title), nil
	})

	// register function to add new series
	shell.Register("add", func(args ...string) (string, error) {
		tt := ""
		if len(args) > 0 {
			tt = strings.Join(args, " ")
		}
		shs, err := d.glibrary.QueryShowsByTitle(tt)
		if err != nil || len(shs) == 0 {
			return "Could not find any series matching the name you provided.", err
		}

		sh := shs[0]

		if ush, err := d.ulibrary.GetShow(sh.ID); err != nil {
			if err != minutes.ErrNotFound {
				return "There was an issue getting the user's library.", err
			}
		} else if ush != nil {
			return "You already have this series in your library.", nil
		}

		if err := d.ulibrary.UpsertShow(sh); err != nil {
			shell.Println()
			return "Could not find series by name.", err
		}

		return fmt.Sprintf("Added '%s' to your library.", sh.Title), nil
	})

	// register function to list series
	shell.Register("list", func(args ...string) (string, error) {
		shs, err := d.ulibrary.GetShows()
		if err != nil {
			return "There was an issue getting the user's library.", err
		} else if len(shs) == 0 {
			return "You do not have any shows in your library.", err
		}

		shell.Println(fmt.Sprintf("You have %d series.", len(shs)))

		for _, sh := range shs {
			shell.Println(fmt.Sprintf("> %s", sh.Title))
		}

		return "", nil
	})

	// register function to run diff
	shell.Register("diff", func(args ...string) (string, error) {
		shell.Println("Running diff...")
		d.Diff()
		return fmt.Sprintf("Diff completed."), nil
	})

	// register function to run watch
	shell.Register("watch", func(args ...string) (string, error) {
		pt := d.config.SeriesPath
		if len(args) > 0 {
			pt = strings.Join(args, " ")
		}
		shell.Println(fmt.Sprintf("Running a new watcher on %s.", pt))
		d.watcher.Watch(pt)
		return fmt.Sprintf("Watcher started."), nil
	})

	// start shell
	shell.Start()
}
