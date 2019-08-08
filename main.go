package main

import (
	"log"
	"path"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/radareorg/r2pm/pkg"
	"github.com/radareorg/r2pm/pkg/database"
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
	// init
	//

	initCmd := &cobra.Command{
		Use:  "init",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if err := database.Init(r2pmDir); err != nil {
				log.Fatal(err)
			}
		},
	}

	rootCmd.AddCommand(initCmd)

	//
	// init
	//

	installCmd := &cobra.Command{
		Use:   "install",
		Short: "install a package",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			packageName := args[0]

			pi, err := database.FindPackage(r2pmDir, packageName)
			if err != nil {
				log.Fatalf("could not find package %s: %v", packageName, err)
			}

			if err := pi.Install(r2pmDir); err != nil {
				log.Fatalf("could not install %s: %v", packageName, err)
			}
		},
	}

	rootCmd.AddCommand(installCmd)

	//
	// uninstall
	//

	uninstallCmd := &cobra.Command{
		Use:  "uninstall",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			packageName := args[0]

			pi, err := database.FindPackage(r2pmDir, packageName)
			if err != nil {
				log.Fatalf("could not find package %s: %v", packageName, err)
			}

			if err := pi.Uninstall(r2pmDir); err != nil {
				log.Fatalf("could not uninstall %s: %v", packageName, err)
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
