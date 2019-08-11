package main

import (
	"flag"
	"fmt"
	"strings"
	"os"
	"os/exec"
	"runtime"
	"path"
	"io/ioutil"
	"strconv"
	"path/filepath"
	"errors"
	"encoding/json"
	"gopkg.in/yaml.v2"
)

const VERSION string = "1.0"

var R2PM_DIR string
var R2PM_GITDIR string
var R2PM_DB string
var DBFILE string

/* General util functions */
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
	if err != nil {
		return false, err
	}
	isDir := fileInfo.IsDir() == true
	return isDir, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func caseInsensitiveContains(a, b string) bool {
	return strings.Contains(strings.ToUpper(a), strings.ToUpper(b))
}

func gitClone(repoPath string, repoUrl string, args ...string) {
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		fmt.Println("Downloading repository...")
		cmdArgs := append([]string{"clone", repoUrl}, args...)
		cmdArgs = append(cmdArgs, repoPath)
		cmd := exec.Command("git", cmdArgs...)
		_, err := cmd.CombinedOutput()
		check(err)
		fmt.Println("Download complete.")
	} else if err == nil {
		fmt.Println("Repository already downloaded, updating...")
		os.Chdir(repoPath)
		cmd := exec.Command("git", "pull")
		_, err := cmd.Output()
		check(err)
	} else {
		check(err)
	}
}

/* R2PM utils */
type PackageInfo struct {
	Name      string   `yaml:"name"`
	Type      string   `yaml:"type"`
	Repo      string   `yaml:"repo"`
	Desc      string   `yaml:"desc"`
	Install   []string `yaml:"install"`
	Uninstall []string `yaml:"uninstall"`
}

func getPackagesList() []string {
	dat, err := ioutil.ReadFile(DBFILE)
	var packagesList []string
	if err != nil {
		r2pmInit()
		dat, err = ioutil.ReadFile(DBFILE)
		if err != nil {
			fmt.Println("Could not read database file " + DBFILE + ". Did you initialize r2pm? (via r2pm -init)")
			return packagesList
		}
	}
	json.Unmarshal(dat, &packagesList)
	return packagesList
}

func getPackageInfo(pkg string) (PackageInfo, error) {
	var file string
	// Empty name, return empty
	if pkg == "" {
		return PackageInfo{}, errors.New("Package name invalid.")
	}

	// If the name contains a file separator, it's surely the full path
	idx := strings.Index(pkg, string(os.PathSeparator))
	if idx >= 0 {
		file = pkg
	} else {
		file = path.Join(R2PM_DB, pkg)
	}
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Package '" + pkg + "' not found.")
		return PackageInfo{}, err
	}
	pinfo := PackageInfo{}
	err = yaml.Unmarshal([]byte(dat), &pinfo)
	if err != nil {
		return PackageInfo{}, err
	}

	return pinfo, err
}

/* R2PM functions */
func r2pmInit() {
	// Make sure git directory exists
	if _, err := os.Stat(R2PM_GITDIR); os.IsNotExist(err) {
		os.MkdirAll(R2PM_GITDIR, 0755)
	}

	// Check if radare2-pm was already cloned
	repoPath := path.Join(R2PM_GITDIR, "r2pm-db")
	repoUrl := "https://github.com/radareorg/r2pm-db"
	gitClone(repoPath, repoUrl, "--depth=3", "--recursive")

	// Initialize database
	var validPackages []string
	err := filepath.Walk(R2PM_DB, func(file string, info os.FileInfo, err error) error {
		dir, err := isDirectory(file)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		if dir {
			return nil
		}

		// Read file content and parse it
		pinfo, err := getPackageInfo(file)
		if err != nil {
			return nil
		}

		// Validate package
		if (pinfo.Name != filepath.Base(file)) {
			fmt.Println("Invalid package name in '" + file + "': '" + pinfo.Name + "'")
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

func r2pmRemove() {
	if _, err := os.Stat(R2PM_DIR); os.IsNotExist(err) {
		return
	}

	err := os.RemoveAll(R2PM_DIR)
	check(err)

	fmt.Println("Deleted " + R2PM_DIR)
}

func r2pmInfo() {
	// Read database file
	packagesList := getPackagesList()
	if len(packagesList) == 0 {
		return
	}
	packagesNumber := strconv.Itoa(len(packagesList))

	fmt.Println("# Repository Database:")
	fmt.Println("There are " + packagesNumber + " packages available.")
	fmt.Println("# Installed:")
	fmt.Println("TODO")
}

func r2pmInstall(pkg string) bool {
	// Get package information
	pinfo, err := getPackageInfo(pkg)
	if (err != nil) {
		return false
	}

	// Download content
	var newdir string
	if pinfo.Type == "git" {
		repoPath := path.Join(R2PM_GITDIR, filepath.Base(pinfo.Repo))
		newdir = repoPath
		gitClone(repoPath, pinfo.Repo)
	}
	fmt.Println("Entering " + newdir)
	os.Chdir(newdir)

	// Install
	fmt.Println("Installing " + pkg + "...")
	for _, command := range pinfo.Install {
		fmt.Println(command)
		commandArgs := strings.Fields(command)
		prog := commandArgs[0]
		if prog == "cd" {
			os.Chdir(commandArgs[1])
			continue
		}
		commandArgs = commandArgs[1:]
		cmd := exec.Command(prog, commandArgs...)
		_, err := cmd.Output()
		check(err)
	}
	return true
}

func r2pmUninstall(pkg string) bool {
	// Get package information
	pinfo, err := getPackageInfo(pkg)
	check(err)

	// Go into package folder
	dir := path.Join(R2PM_GITDIR, filepath.Base(pinfo.Repo))
	fmt.Println("Entering " + dir)
	os.Chdir(dir)

	// Uninstall
	fmt.Println("Uninstalling " + pkg + "...")
	for _, command := range pinfo.Uninstall {
		fmt.Println(command)
		commandArgs := strings.Fields(command)
		prog := commandArgs[0]
		if prog == "cd" {
			os.Chdir(commandArgs[1])
			continue
		}
		commandArgs = commandArgs[1:]
		cmd := exec.Command(prog, commandArgs...)
		_, err := cmd.Output()
		check(err)
	}
	return true
}

func r2pmSearch(pkg string) bool {
	packagesList := getPackagesList()
	anyFound := false

	if len(packagesList) == 0 {
		return false
	}

	headMsg := "List of available packages: "
	if pkg != "" {
		headMsg += "(filter: " + pkg + ")"
	}
	fmt.Println(headMsg)
	for _, pkgname := range packagesList {
		pinfo, err := getPackageInfo(pkgname)
		if err != nil {
			continue
		}
		if !caseInsensitiveContains(pinfo.Name, pkg) && !caseInsensitiveContains(pinfo.Desc, pkg) {
			continue
		}
		fmt.Println(pinfo.Name + "\t\t" + pinfo.Desc)
		anyFound = true
	}
	if !anyFound {
		fmt.Println("No packages found.")
	}

	return true
}

func _main() {
	// Initialize environment variables
	var r2pmdir string
	if isWindows() {
		r2pmdir = getenv("APPDATA", "")
	} else {
		// TODO Use XDG env variable and fallback to this
		r2pmdir = path.Join(getenv("HOME", ""), ".local/share/radare2/r2pm")
	}
	R2PM_DIR = getenv("R2PM_DIR", r2pmdir)
	R2PM_GITDIR = getenv("R2PM_GITDIR", path.Join(r2pmdir, "git"))
	R2PM_DB = getenv("R2PM_DB", path.Join(R2PM_GITDIR, "r2pm-db", "db"))
	DBFILE = path.Join(R2PM_DIR, "db.json")

	// Parse arguments
	versionPtr := flag.Bool("v", false, "Show r2pm version")
	initPtr := flag.Bool("init", false, "Init the repository")
	deletePtr := flag.Bool("delete", false, "Delete the whole local r2pm repository")
	infoPtr := flag.Bool("i", false, "Show information or install a package")
	uninstallPtr := flag.Bool("u", false, "Uninstall a package")
	searchPtr := flag.Bool("s", false, "Search into database")
	flag.Parse()

	if *versionPtr == true {
		fmt.Println("r2pm " + VERSION)
		return
	}

	if *deletePtr == true {
		r2pmRemove()
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

	if *uninstallPtr == true {
		if len(flag.Args()) != 1 {
			fmt.Println("No package name specified!")
		} else {
			r2pmUninstall(flag.Args()[0])
		}
		return
	}

	if *searchPtr == true {
		if len(flag.Args()) == 1 {
			r2pmSearch(flag.Args()[0])
		} else {
			r2pmSearch("")
		}
		return
	}

	fmt.Println("No action given.")
	flag.PrintDefaults();
}
