package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/urfave/cli"

	"github.com/radareorg/r2pm/internal/features"
	"github.com/radareorg/r2pm/pkg"
	"github.com/radareorg/r2pm/pkg/r2package"
)

func r2pmDir() string {
	var defaultDir string

	if runtime.GOOS == "windows" {
		defaultDir = pkg.GetenvDefault("APPDATA", "")
	} else {
		// TODO Use XDG env variable and fallback to this
		defaultDir = path.Join(
			pkg.GetenvDefault("HOME", ""),
			".local/share/radare2/r2pm")
	}

	return pkg.GetenvDefault("R2PM_DIR", defaultDir)
}

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
	r2pmDir := r2pmDir()

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
		cli.BoolFlag{
			Name:   flagNameDebug,
			Usage:  "enable debug logs",
			EnvVar: "R2PM_DEBUG",
		},
	}

	app.Before = func(c *cli.Context) error {
		features.SetDebug(c.Bool(flagNameDebug))
		return nil
	}

	app.Commands = []cli.Command{
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
				cli.StringFlag{
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
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "list packages",
			Action:  listAvailablePackages,
			Subcommands: []cli.Command{
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
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func printPackageSlice(packages []r2package.Info) {
	for _, p := range packages {
		fmt.Printf("%s: %s\n", p.Name, p.Desc)
	}
}
