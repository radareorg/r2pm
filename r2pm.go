package main

import "flag"
import "fmt"
import "os"
import "os/exec"
import "runtime"
import "path"
import "io/ioutil"
import "strconv"
import "path/filepath"
import "encoding/json"

const VERSION string = "1.0"
const R2PM_LOCAL string = ".local/share/radare2/r2pm"

var R2PM_DIR string
var R2PM_GITDIR string
var R2PM_DB string
var DBFILE string

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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir(), err
}

/* R2PM functions */
type PackageInfo struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Repo      string   `json:"repo"`
	Desc      string   `json:"desc"`
	Install   []string `json:"install"`
	Uninstall []string `json:"uninstall"`
}

func r2pmInit() {
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
		check(err)
		fmt.Println("Download complete.")
	} else {
		fmt.Println("Repository already downloaded, updating...")
		os.Chdir(repo)
		cmd := exec.Command("git", "pull")
		_, err := cmd.Output()
		check(err)
	}

	// Initialize database
	packageRepo := path.Join(repo, "db")
	var validPackages []string
	err := filepath.Walk(packageRepo, func(path string, info os.FileInfo, err error) error {
		dir, _ := isDirectory(path)
		if dir {
			return nil
		}

		// Read file content and parse it
		dat, err := ioutil.ReadFile(path)
		check(err)
		pinfo := PackageInfo{}
		err = json.Unmarshal([]byte(dat), &pinfo)
		if err != nil {
			return nil
		}

		// Validate package
		if pinfo.Name != filepath.Base(path) {
			fmt.Println("Invalid package name in '" + path + "': '" + pinfo.Name + "'")
			return nil
		}
		validPackages = append(validPackages, pinfo.Name)

		return nil
	})

	// Save valid packages list
	validPackagesJson, _ := json.Marshal(validPackages)
	err = ioutil.WriteFile(DBFILE, validPackagesJson, 0644)
	check(err)
}

func r2pmInfo() {
	fmt.Println("# Repository Database:")

	// Read database file
	dat, err := ioutil.ReadFile(DBFILE)
	check(err)
	var packagesList []string
	json.Unmarshal(dat, &packagesList)
	packagesNumber := strconv.Itoa(len(packagesList))

	fmt.Println("# " + packagesNumber + " Packages")
	fmt.Println("# Installed:")
	fmt.Println("TODO")
}

func r2pmInstall(pkg string) bool {
	fmt.Println("Installing " + pkg)
	fmt.Println("TODO")
	return true
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
	DBFILE = path.Join(R2PM_DIR, "db.json")

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
		r2pmInit()
		return
	}

	if *infoPtr == true {
		if len(flag.Args()) == 1 {
			r2pmInstall(flag.Args()[0])
		} else {
			r2pmInfo()
		}
		return
	}
}
