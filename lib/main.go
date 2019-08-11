package main

import "C"

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/radareorg/r2pm/internal/features"
)

func init() {
	// Disable the logger by default
	R2pmSetDebug(0)
}

func getReturnValue(err error) C.int {
	if err != nil {
		log.Fatal(err.Error())
	}

	return 0
}

//export R2pmDelete
func R2pmDelete(r2pmDir *C.char) C.int {
	err := features.Delete(C.GoString(r2pmDir))
	return getReturnValue(err)
}

//export R2pmInit
func R2pmInit(r2pmDir *C.char) C.int {
	err := features.Init(C.GoString(r2pmDir))
	return getReturnValue(err)
}

//export R2pmInstall
func R2pmInstall(r2pmDir, packageName *C.char) C.int {
	err := features.Install(C.GoString(r2pmDir), C.GoString(packageName))
	return getReturnValue(err)
}

//export R2pmList
func R2pmList(r2pmDir *C.char) ([]string, C.int) {
	entries, err := features.ListInstalled(C.GoString(r2pmDir))
	return entries, getReturnValue(err)
}

//export R2pmUninstall
func R2pmUninstall(r2pmDir, packageName *C.char) C.int {
	err := features.Uninstall(C.GoString(r2pmDir), C.GoString(packageName))
	return getReturnValue(err)
}

//export R2pmSetDebug
func R2pmSetDebug(value C.int) {
	if value == 0 {
		log.SetOutput(ioutil.Discard)
		return
	}

	log.SetOutput(os.Stderr)
	log.SetPrefix("libr2pm: ")

	log.Print("debug enabled")
}

func main() {}
