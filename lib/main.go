package main

import "C"

import (
	"log"

	"github.com/radareorg/r2pm/internal/features"
)

var debug = false

func getReturnValue(err error) C.int {
	if err != nil {
		if debug {
			log.Print(err.Error())
		}

		return -1
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
	entries, err := features.List(C.GoString(r2pmDir))
	return entries, getReturnValue(err)
}

//export R2pmUninstall
func R2pmUninstall(r2pmDir, packageName *C.char) C.int {
	err := features.Uninstall(C.GoString(r2pmDir), C.GoString(packageName))
	return getReturnValue(err)
}

//export R2pmSetDebug
func R2pmSetDebug(value C.int) {
	debug = value != 0

	if debug {
		log.Print("debug enabled")
	}
}

func main() {}
