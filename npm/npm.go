package npm

import (
	"regexp"
	"strings"
	"io/ioutil"
	"path/filepath"
	"os"
	"encoding/json"
	"os/exec"
	"log"
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

type NpmListOutput struct {
	Dependencies map[string]map[string]string
}

func ExtractVersions(data NpmListOutput) map[string]string {
	result := make(map[string]string)
	dependencies := data.Dependencies
	for name, properties := range dependencies {
		if version, ok := properties["version"]; ok {
		    result[name] = version
		}
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
		currentVersion := packageContent[dependencyVersionStart:dependencyVersionEnd]
		if currentVersion == version {
			return "", &DependencyUpToDate{packageContent, dependency}
		}
		return packageContent[:dependencyVersionStart] + version + packageContent[dependencyVersionEnd:], nil
	default:
		return "", &InvalidPackageJsonContent{packageContent, dependency}
	}
}

func UpdateDependencies(packageContent string, updates map[string]string) (string, error) {
	for name, version := range updates {
		newPackageContent, err := UpdateDependency(packageContent, name, version)
		if err == nil {
			packageContent = newPackageContent
		} else {
			switch t := err.(type) {
			default:
				return "", err
			case *DependencyUpToDate:
				log.Print("FYI: ", t.Error())
			case *DependencyNotFound:
				log.Print("FYI: ", t.Error())
			}
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

func ExecNpmList(dir string) (map[string]string, error) {
	log.Print("Executing npm list")
	cmd := exec.Command("bash", "-c", "source ~/.nvm/nvm.sh && nvm i 6 >/dev/null && npm i >/dev/null && npm list --depth 0 --json || echo")
	cmd.Dir = dir

	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Print("Npm list outPipe failed: ", err.Error())
		return nil, err
	}

	errPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Print("Npm list errPipe failed: ", err.Error())
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		log.Print("Npm list Start failed: ", err.Error())
		return nil, err
	}


	// pring std and err output after done, but before checking for success
	outBytes, _ := ioutil.ReadAll(outPipe)
	log.Print("Npm list done: <", string(outBytes), ">")

	errBytes, _ := ioutil.ReadAll(errPipe)
	log.Print("Npm list error output: <", string(errBytes), ">")

	err = cmd.Wait()
	if err != nil {
		log.Print("Npm list failed: ", err.Error())
		return nil, err
	}

	versions := ParseNpmListOutput(outBytes)
	log.Print("Parsing successful, effective dependency versions: ", versions)
	return versions, nil
}

func FreezePackage(dir string) error {
	log.Print("Freezing dependencies in: ", dir)
	packageFile := filepath.Join(dir, "package.json")
	if !exists(packageFile) {
		log.Print("No package.json found, aborting")
		return &NoPackageJson{dir}
	}
	bytes, err := ioutil.ReadFile(packageFile)
	if err != nil {
		log.Print("Could not read package.json", err.Error())
		return err
	}
	packageContent := string(bytes)

	Exec(dir, "rm", "-rf", "node_modules")
	Exec(dir, "bash", "-c", "nvm i 6 && npm i")
	versions, err := ExecNpmList(dir)
	if err != nil {
		return err
	}

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
