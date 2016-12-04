package npm

import (
	"regexp"
	"strings"
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

func UpdateDependencyVersion(packageContent string, dependency string, version string) (string, error) {
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
