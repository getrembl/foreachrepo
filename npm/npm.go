package npm

import (
	"regexp"
	"strings"
	"io/ioutil"
	"path/filepath"
	"os"
)

const SPACES = "\\s*"
const SEMICOLON = ":"
const VERSION_PATTERN = "[^\"]+"

func quote(str string) string {
	return "\"" + str + "\""
}

func group(str string) string {
	return "(" + str + ")"
}

type DependencyNotFound struct {
	packageContent string
	dependency     string
}

func (D *DependencyNotFound) Error() string {
	return "Dependency " + D.dependency + " not found in package.json"
}

type DependencyUpToDate struct {
	packageContent string
	dependency     string
}

func (D *DependencyUpToDate) Error() string {
	return "Dependency " + D.dependency + " is already up to date"
}

type InvalidPackageJsonContent struct {
	packageContent string
	dependency     string
}

func (D *InvalidPackageJsonContent) Error() string {
	return "Invalid package.json content " + D.packageContent
}

func GetUpdatedPackageContent(packageContent string, dependency string, version string) (string, error) {
	elements := []string{
		SPACES,
		quote(group(regexp.QuoteMeta(dependency))),
		SPACES,
		SEMICOLON,
		SPACES,
		quote(group(VERSION_PATTERN)),
	}
	pattern := strings.Join(elements, "")

	r, _ := regexp.Compile(pattern)

	matches := r.FindAllStringSubmatchIndex(packageContent, -1)

	switch len(matches) {
	case 0:
		return "", &DependencyNotFound{packageContent, dependency}
	case 1:
		match := matches[0]
		dependencyVersionStart := match[4]
		dependencyVersionEnd := match[5]
		currentVersion := packageContent[dependencyVersionStart:dependencyVersionEnd]
		if currentVersion == version {
			return "", &DependencyUpToDate{packageContent, dependency}
		}
		return packageContent[:dependencyVersionStart] + version + packageContent[dependencyVersionEnd:], nil
	default:
		return "", &InvalidPackageJsonContent{packageContent, dependency}
	}
}

// Exists reports whether the named file or directory exists.
func exists(name string) bool {
    if _, err := os.Stat(name); err != nil {
    if os.IsNotExist(err) {
                return false
            }
    }
    return true
}

type NoPackageJson struct {
	dir string
}

func (D *NoPackageJson) Error() string {
	return "No package.json found in " + D.dir
}

func UpdatePackage(dir string, dependency string, version string) error {
	packageFile := filepath.Join(dir, "package.json")
	if !exists(packageFile) {
		return &NoPackageJson{dir}
	}
	bytes, err := ioutil.ReadFile(packageFile)
	if err != nil {
		return err
	}
	packageContent := string(bytes)
	updatedPackageContent, err := GetUpdatedPackageContent(packageContent, dependency, version)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(packageFile, []byte(updatedPackageContent), 0644)
	if err != nil {
		return err
	}
	return nil
}
