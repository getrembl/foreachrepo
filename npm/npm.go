package npm

import (
	"regexp"
	"strings"
	"io/ioutil"
	"path/filepath"
	"os"
	"encoding/json"
	"os/exec"
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

type InvalidPackageJsonContent struct {
	packageContent string
	dependency     string
}

func (D *InvalidPackageJsonContent) Error() string {
	return "Invalid package.json content " + D.packageContent
}

type NpmListOutput struct {
	Dependencies map[string]map[string]string
}

func ExtractVersions(data NpmListOutput) map[string]string {
	result := make(map[string]string)
	dependencies := data.Dependencies
	for name, properties := range dependencies {
		result[name] = properties["version"]
	}
	return result
}

func ParseNpmListOutput(output []byte) map[string]string {
	var data NpmListOutput

	json.Unmarshal(output, &data)

	return ExtractVersions(data)
}

func UpdateDependency(packageContent string, dependency string, version string) (string, error) {
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
		return packageContent[:dependencyVersionStart] + version + packageContent[dependencyVersionEnd:], nil
	default:
		return "", &InvalidPackageJsonContent{packageContent, dependency}
	}
}

func UpdateDependencies(packageContent string, updates map[string]string) (string, error) {
	var err error
	for name, version := range updates {
		packageContent, err = UpdateDependency(packageContent, name, version)
		if err != nil {
			return "", err
		}
	}
	return packageContent, nil
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
	updatedPackageContent, err := UpdateDependency(packageContent, dependency, version)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(packageFile, []byte(updatedPackageContent), 0644)
	if err != nil {
		return err
	}
	return nil
}

func Exec(dir string, name string, elements ...string) error {
	cmd := exec.Command(name, elements...)
	cmd.Dir = dir
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func ExecNpmList(dir string) (NpmListOutput, error) {
	cmd := exec.Command("npm", "list", "--depth", "0", "--json")
	cmd.Dir = dir
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	var result NpmListOutput
	if err := json.NewDecoder(stdout).Decode(&result); err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		return nil, err
	}
	return result, nil
}

func FreezePackage(dir string) error {
	packageFile := filepath.Join(dir, "package.json")
	if !exists(packageFile) {
		return &NoPackageJson{dir}
	}
	bytes, err := ioutil.ReadFile(packageFile)
	if err != nil {
		return err
	}
	packageContent := string(bytes)

	Exec(dir, "rm", "-rf", "node_modules")
	Exec(dir, "npm", "i")
	npmListOutput, err := ExecNpmList(dir)
	if err != nil {
		return nil, err
	}

	versions := ExtractVersions(npmListOutput)

	updatedPackageContent, err := UpdateDependencies(packageContent, versions)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(packageFile, []byte(updatedPackageContent), 0644)
	if err != nil {
		return err
	}

	return nil
}
