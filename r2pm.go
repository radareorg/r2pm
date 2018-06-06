package main

import "flag"
import "fmt"
import "os"
import "os/exec"
import "runtime"
import "path"
import "io/ioutil"
import "strconv"

const VERSION string = "1.0"
const R2PM_LOCAL string = ".local/share/radare2/r2pm"

var R2PM_DIR string
var R2PM_GITDIR string
var R2PM_DB string
var WD string

/* Utils functions */
func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}

/* R2PM functions */
func r2pm_init() {
	// Make sure git directory exists
	if _, err := os.Stat(R2PM_GITDIR); os.IsNotExist(err) {
		os.Mkdir(R2PM_GITDIR, 0)
	}

	// Check if radare2-pm was already cloned
	repo := path.Join(R2PM_GITDIR, "radare2-pm")
	if _, err := os.Stat(repo); os.IsNotExist(err) {
		os.Chdir(R2PM_GITDIR)
		fmt.Println("Downloading repository...")
		cmd := exec.Command("git", "clone", "https://github.com/radare/radare2-pm", "--depth=3", "--recursive")
		_, err := cmd.Output()
		if err != nil {
			panic(err)
		}
		fmt.Println("Download complete.")
	} else {
		fmt.Println("Repository already initialized.")
	}
}

func r2pm_info() {
	fmt.Println("# Repository Database:")
	files, _ := ioutil.ReadDir(R2PM_DB)
	nbPackages := strconv.Itoa(len(files))
	fmt.Println("# " +  nbPackages + " Packages")
	fmt.Println("# Installed:")
	fmt.Println("TODO")
}

func r2pm_install(pck string) bool {
	fmt.Println("Installing " +  pck)
	return true;
}

func main() {
	// Initialize environment variables
	var r2pmdir string
	if isWindows() {
		r2pmdir = getenv("APPDATA", "")
	} else {
		r2pmdir = path.Join(getenv("HOME", ""), R2PM_LOCAL)
	}
	R2PM_DIR = getenv("R2PM_DIR", r2pmdir)
	R2PM_GITDIR = getenv("R2PM_GITDIR", path.Join(r2pmdir, "git"))
	R2PM_DB = getenv("R2PM_DB", path.Join(R2PM_GITDIR, "radare2-pm", "db"))
	WD, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Parse arguments
	versionPtr := flag.Bool("v", false, "Show r2pm version")
	initPtr := flag.Bool("init", false, "Init the repository")
	infoPtr := flag.Bool("i", false, "Show information")
	flag.Parse()

	if *versionPtr == true {
		fmt.Println("r2pm " + VERSION)
		return
	}

	if *initPtr == true {
		r2pm_init()
		return
	}

	if *infoPtr == true {
		if len(flag.Args()) == 1 {
			r2pm_install(flag.Args()[0])
		} else {
			r2pm_info()
		}
		return
	}

	fmt.Println("Using this dir: ", R2PM_DIR)
	fmt.Println("Using this dir: ", WD)
}
