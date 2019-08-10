package features

import (
	"log"

	"github.com/radareorg/r2pm/pkg/database"
)

func Delete(r2pmDir string) error {
	return database.Delete(r2pmDir)
}

func Init(r2pmDir string) error {
	return database.Init(r2pmDir)
}

func Install(r2pmDir, packageName string) error {
	pi, err := database.FindPackage(r2pmDir, packageName)
	if err != nil {
		log.Fatalf("could not find package %s: %v", packageName, err)
	}

	if err := pi.Install(r2pmDir); err != nil {
		log.Fatalf("could not install %s: %v", packageName, err)
	}

	return nil
}

func List(r2pmDir string) ([]string, error) {
	return database.List(r2pmDir)
}

func Uninstall(r2pmDir, packageName string) error {
	pi, err := database.FindPackage(r2pmDir, packageName)
	if err != nil {
		log.Fatalf("could not find package %s: %v", packageName, err)
	}

	if err := pi.Uninstall(r2pmDir); err != nil {
		log.Fatalf("could not uninstall %s: %v", packageName, err)
	}

	return nil
}
