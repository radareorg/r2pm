package main

import (
	"log"
	"path"
	"runtime"

	"github.com/spf13/cobra"

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

func main() {
	r2pmDir := r2pmDir()

	rootCmd := &cobra.Command{
		Use:     "r2pm",
		Short:   "r2 package manager",
		Version: "1.0.0",
	}

	//
	// delete
	//

	deleteCmd := &cobra.Command{
		Use:  "delete",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if err := features.Delete(r2pmDir); err != nil {
				log.Fatal(err)
			}
		},
	}

	rootCmd.AddCommand(deleteCmd)

	//
	// init
	//

	initCmd := &cobra.Command{
		Use:  "init",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if err := features.Init(r2pmDir); err != nil {
				log.Fatal(err)
			}
		},
	}

	rootCmd.AddCommand(initCmd)

	//
	// install
	//

	installCmd := &cobra.Command{
		Use:   "install",
		Short: "install a package",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := features.Install(r2pmDir, args[0]); err != nil {
				log.Fatal(err)
			}
		},
	}

	rootCmd.AddCommand(installCmd)

	//
	// list
	//

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "list all available packages",
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			packages, err := features.List(r2pmDir)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("%d packages available:", len(packages))

			for _, p := range packages {
				log.Print(p)
			}
		},
	}

	rootCmd.AddCommand(listCmd)

	//
	// uninstall
	//

	uninstallCmd := &cobra.Command{
		Use:  "uninstall",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := features.Uninstall(r2pmDir, args[0]); err != nil {
				log.Fatal(err)
			}
		},
	}

	rootCmd.AddCommand(uninstallCmd)

	//
	// update
	//

	updateCmd := &cobra.Command{
		Use: "update",
	}

	rootCmd.AddCommand(updateCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
