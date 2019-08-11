package main

import (
	"log"
	"os"
	"path"
	"runtime"

	"github.com/urfave/cli"

	"github.com/radareorg/r2pm/internal/features"
	"github.com/radareorg/r2pm/pkg"
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

	app := cli.NewApp()
	app.Name = "r2pm"
	app.Usage = "r2 package manager"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:  "delete",
			Usage: "delete the local package database",
			Action: func(*cli.Context) error {
				return features.Delete(r2pmDir)
			},
		},
		{
			Name:  "init",
			Usage: "initialize the local package database",
			Action: func(*cli.Context) error {
				return features.Init(r2pmDir)
			},
		},
		{
			Name:      "install",
			Usage:     "install a package",
			ArgsUsage: "PACKAGE",
			Action: func(c *cli.Context) error {
				packageName := getArgumentOrExit(c)

				return features.Install(r2pmDir, packageName)
			},
		},
		{
			Name:  "list",
			Usage: "list all the available packages",
			Action: func(c *cli.Context) error {
				packages, err := features.List(r2pmDir)
				if err != nil {
					return err
				}

				log.Printf("%d packages available", len(packages))

				for _, p := range packages {
					log.Print(p)
				}

				return nil
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

				log.Printf("Your search returned %d matches", len(matches))

				for _, m := range matches {
					log.Println(m)
				}

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
