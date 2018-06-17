package main

import "flag"
import "fmt"
import "strings"
import "os"
import "os/exec"
import "runtime"
import "path"
import "io/ioutil"
import "strconv"
import "path/filepath"
import "encoding/json"
import "errors"

const VERSION string = "1.0"
const R2PM_LOCAL string = ".local/share/radare2/r2pm"

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

func caseInsenstiveContains(a, b string) bool {
	return strings.Contains(strings.ToUpper(a), strings.ToUpper(b))
}

func gitClone(repoPath string, repoUrl string, args ...string) {
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		os.Chdir(R2PM_GITDIR)
		fmt.Println("Downloading repository...")
		cmdArgs := []string{"clone", repoUrl} // + args
		cmd := exec.Command("git", cmdArgs...)
		_, err := cmd.Output()
		check(err)
		fmt.Println("Download complete.")
	} else {
		fmt.Println("Repository already downloaded, updating...")
		os.Chdir(repoPath)
		cmd := exec.Command("git", "pull")
		_, err := cmd.Output()
		check(err)
	}
}

/* R2PM utils */
type PackageInfo struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Repo      string   `json:"repo"`
	Desc      string   `json:"desc"`
	Install   []string `json:"install"`
	Uninstall []string `json:"uninstall"`
}

func getPackagesList() []string {
	dat, err := ioutil.ReadFile(DBFILE)
	check(err)
	var packagesList []string
	json.Unmarshal(dat, &packagesList)
	return packagesList
}

func getPackageInfo(pkg string) (PackageInfo, error) {
	var file string
	// Empty name, return empty
	if pkg == "" {
		return PackageInfo{}, errors.New("Package name invalid")
	}

	// If the name contains a file separator, it's surely the full path
	idx := strings.Index(pkg, string(os.PathSeparator))
	if idx >= 0 {
		file = pkg
	} else {
		file = path.Join(R2PM_DB, pkg)
	}
	dat, err := ioutil.ReadFile(file)
	check(err)
	pinfo := PackageInfo{}
	err = json.Unmarshal([]byte(dat), &pinfo)
	if err != nil {
		return PackageInfo{}, err
	}

	return pinfo, err
}

/* R2PM functions */
func r2pmInit() {
	// Make sure git directory exists
	if _, err := os.Stat(R2PM_GITDIR); os.IsNotExist(err) {
		os.Mkdir(R2PM_GITDIR, 0)
	}

	// Check if radare2-pm was already cloned
	repoPath := path.Join(R2PM_GITDIR, "radare2-pm")
	repoUrl := "https://github.com/radare/radare2-pm"
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
		if pinfo.Name != filepath.Base(file) {
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

func r2pmInfo() {
	fmt.Println("# Repository Database:")

	// Read database file
	packagesList := getPackagesList()
	packagesNumber := strconv.Itoa(len(packagesList))

	fmt.Println("There are " + packagesNumber + " packages available.")
	fmt.Println("# Installed:")
	fmt.Println("TODO")
}

func r2pmInstall(pkg string) bool {
	packagesList := getPackagesList()

	// Check if the package is valid
	if stringInSlice(pkg, packagesList) == false {
		fmt.Println("Package " + pkg + " not found!")
		return false
	}

	// Get package information
	pinfo, err := getPackageInfo(pkg)
	check(err)

	// Download content
	var newdir string
	if pinfo.Type == "git" {
		repoPath := path.Join(R2PM_GITDIR, filepath.Base(pinfo.Repo))
		gitClone(repoPath, pinfo.Repo)
	}
	fmt.Println("Entering " + newdir)
	os.Chdir(newdir)

	// Install
	fmt.Println("Installing " + pkg + "...")
	for _, command := range pinfo.Install {
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

	headMsg := "List of packages: "
	if pkg != "" {
		headMsg += "(filter: " + pkg + ")"
	}
	fmt.Println(headMsg)
	for _, pkgname := range packagesList {
		pinfo, err := getPackageInfo(pkgname)
		if err != nil {
			continue
		}
		if !caseInsenstiveContains(pinfo.Name, pkg) && !caseInsenstiveContains(pinfo.Desc, pkg) {
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
	searchPtr := flag.Bool("s", false, "Search into database")
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

	if *searchPtr == true {
		if len(flag.Args()) == 1 {
			r2pmSearch(flag.Args()[0])
		} else {
			r2pmSearch("")
		}
		return
	}
}
