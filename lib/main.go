package main

import (
	"C"

	"github.com/radareorg/r2pm/internal/features"
)

//export R2pmDelete
func R2pmDelete(r2pmDir string) error {
	return features.Delete(r2pmDir)
}

//export R2pmInit
func R2pmInit(r2pmDir string) error {
	return features.Init(r2pmDir)
}

//export R2pmInstall
func R2pmInstall(r2pmDir, packageName string) error {
	return features.Install(r2pmDir, packageName)
}

//export R2pmList
func R2pmList(r2pmDir string) ([]string, error) {
	return features.List(r2pmDir)
}

//export R2pmUninstall
func R2pmUninstall(r2pmDir, packageName string) error {
	return features.Uninstall(r2pmDir, packageName)
}

func main() {}
