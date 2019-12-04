package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/radareorg/r2pm/internal/features"
	"github.com/radareorg/r2pm/internal/util/dir"
	"github.com/radareorg/r2pm/pkg/r2package"
)

func getArgumentOrExit(c *cli.Context) string {
	packageName := c.Args().First()

	if packageName == "" {
		if err := cli.ShowSubcommandHelp(c); err != nil {
			log.Fatal(err)
		}

		os.Exit(1)
	}

	return packageName
}

func main() {
	r2Dir := dir.R2Dir()
	r2pmDir := dir.SiteDir()

	listAvailablePackages := func(c *cli.Context) error {
		packages, err := features.ListAvailable(r2pmDir)
		if err != nil {
			return err
		}

		fmt.Printf("%d available packages\n", len(packages))
		printPackageSlice(packages)

		return nil
	}

	const flagNameDebug = "debug"

	app := cli.NewApp()
	app.Name = "r2pm"
	app.Usage = "r2 package manager"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    flagNameDebug,
			Usage:   "enable debug logs",
			EnvVars: []string{features.DebugEnvVar},
		},
	}

	app.Before = func(c *cli.Context) error {
		features.SetDebug(c.Bool(flagNameDebug))
		return nil
	}

	app.Commands = []*cli.Command{
		{
			Name:  "delete",
			Usage: "delete the local package database",
			Action: func(*cli.Context) error {
				return features.Delete(r2pmDir)
			},
		},
		{
			Name:    "init",
			Aliases: []string{"update"},
			Usage:   "initialize or update the local package database",
			Action: func(*cli.Context) error {
				return features.Init(r2pmDir)
			},
		},
		{
			Name:      "install",
			Usage:     "install a package",
			ArgsUsage: "[PACKAGE]",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "f",
					Usage: "install a package described by a local file",
				},
			},
			Action: func(c *cli.Context) error {
				if path := c.String("f"); path != "" {
					log.Print("Installing " + path)
					return features.InstallFromFile(r2pmDir, path)
				}

				packageName := getArgumentOrExit(c)

				return features.Install(r2pmDir, packageName)
			},
			Subcommands: []*cli.Command{
				{
					Name:      "radare2",
					Usage:     "install radare2",
					ArgsUsage: "VERSION",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:  "p",
							Usage: "radare2's configure --prefix",
							Value: r2Dir,
						},
					},
					Action: func(c *cli.Context) error {
						if c.NArg() != 1 {
							return errors.New("a version number is required")
						}

						version := c.Args().First()

						prefix := c.String("p")
						if prefix == "" {
							return errors.New("A prefix is required")
						}

						return features.InstallRadare2(r2pmDir, r2Dir, version)
					},
				},
			},
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "list packages",
			Action:  listAvailablePackages,
			Subcommands: []*cli.Command{
				{
					Name:   "available",
					Usage:  "list all the available packages",
					Action: listAvailablePackages,
				},
				{
					Name:  "installed",
					Usage: "list all the installed packages",
					Action: func(c *cli.Context) error {
						packages, err := features.ListInstalled(r2pmDir)
						if err != nil {
							return err
						}

						fmt.Printf("%d installed packages\n", len(packages))
						printPackageSlice(packages)

						return nil
					},
				},
			},
		},
		{
			Name:      "search",
			Usage:     "search for a package in the database",
			ArgsUsage: "PATTERN",
			Action: func(c *cli.Context) error {
				pattern := getArgumentOrExit(c)

				matches, err := features.Search(r2pmDir, pattern)
				if err != nil {
					return err
				}

				fmt.Printf("Your search returned %d matches\n", len(matches))
				printPackageSlice(matches)

				return nil
			},
		},
		{
			Name:      "uninstall",
			Usage:     "uninstall a package",
			ArgsUsage: "PACKAGE",
			Action: func(c *cli.Context) error {
				packageName := getArgumentOrExit(c)

				return features.Uninstall(r2pmDir, packageName)
			},
			Subcommands: []*cli.Command{
				{
					Name:  "radare2",
					Usage: "uninstall radare2",
					Action: func(c *cli.Context) error {
						return features.UninstallRadare2(r2pmDir, r2Dir)
					},
				},
			},
		},
		{
			Name:      "upgrade",
			Usage:     "upgrade (uninstall and reinstall) a package",
			ArgsUsage: "[PACKAGE]",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "a, all",
					Usage: "upgrade all packages",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("a") {
					return features.UpgradeAll(r2pmDir)
				}

				packageName := getArgumentOrExit(c)

				return features.Upgrade(r2pmDir, packageName)
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

func printPackageSlice(packages []r2package.Info) {
	for _, p := range packages {
		fmt.Printf("%s: %s\n", p.Name, p.Desc)
	}
}
