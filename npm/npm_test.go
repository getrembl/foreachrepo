package npm

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"io/ioutil"
	"path/filepath"
	"os"
)

const SAMPLE_PACKAGE_CONTENT = `{
  "name": "express-middleware",
  "version": "1.4.0",
  "description": "Set of middlewares for Chauffeur-Privé",
  "keywords": [
    "express",
    "middleware",
    "chauffeur-privé"
  ],
  "main": "index.js",
  "dependencies": {
    "@chauffeur-prive/i18n": "^1.2.0",
    "accept-language-parser": "^1.3.0",
    "bunyan": "~1.8.1"
  },
  "devDependencies": {
    "chai": "~3.5.0",
    "eslint-config-cp": "transcovo/eslint-config-cp#1.1.0"
  },
  "scripts": {
    "test": "npm run lint && npm run coverage"
  },
  "repository": {
    "type": "git",
    "url": "git+https://github.com/transcovo/express-middlewares.git"
  },
  "engines": {
    "node": ">=4.2.2",
    "npm": ">=2.14.7"
  },
  "bugs": {
    "url": "https://github.com/transcovo/express-middlewares/issues"
  },
  "homepage": "https://github.com/transcovo/express-middlewares#readme",
  "author": "Chauffeur-Privé",
  "license": "Apache-2.0"
}
`

const EXPECTED_UPDATED_PACKAGE_CONTENT = `{
  "name": "express-middleware",
  "version": "1.4.0",
  "description": "Set of middlewares for Chauffeur-Privé",
  "keywords": [
    "express",
    "middleware",
    "chauffeur-privé"
  ],
  "main": "index.js",
  "dependencies": {
    "@chauffeur-prive/i18n": "^1.2.0",
    "accept-language-parser": "^1.3.0",
    "bunyan": "1.8.2"
  },
  "devDependencies": {
    "chai": "~3.5.0",
    "eslint-config-cp": "transcovo/eslint-config-cp#1.1.0"
  },
  "scripts": {
    "test": "npm run lint && npm run coverage"
  },
  "repository": {
    "type": "git",
    "url": "git+https://github.com/transcovo/express-middlewares.git"
  },
  "engines": {
    "node": ">=4.2.2",
    "npm": ">=2.14.7"
  },
  "bugs": {
    "url": "https://github.com/transcovo/express-middlewares/issues"
  },
  "homepage": "https://github.com/transcovo/express-middlewares#readme",
  "author": "Chauffeur-Privé",
  "license": "Apache-2.0"
}
`

func TestGetUpdatedPackageContentSuccess(t *testing.T) {
	updated, err := GetUpdatedPackageContent(SAMPLE_PACKAGE_CONTENT, "bunyan", "1.8.2")
	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_UPDATED_PACKAGE_CONTENT, updated)
}

func TestGetUpdatedPackageContentNotFound(t *testing.T) {
	updated, err := GetUpdatedPackageContent(SAMPLE_PACKAGE_CONTENT, "not-a-dependency", "1.8.2")
	assert.Equal(t, "", updated)
	assert.NotNil(t, err)
	errMessage := err.Error()
	assert.Contains(t, errMessage, "not-a-dependency")
	assert.Contains(t, errMessage, "not found")
}

func TestGetUpdatedPackageContentUpToDate(t *testing.T) {
	updated, err := GetUpdatedPackageContent(SAMPLE_PACKAGE_CONTENT, "bunyan", "~1.8.1")
	assert.Equal(t, "", updated)
	assert.NotNil(t, err)
	errMessage := err.Error()
	assert.Contains(t, errMessage, "bunyan")
	assert.Contains(t, errMessage, "up to date")
}

func TestUpdatePackage(t *testing.T) {
	dir, mkdir_err := ioutil.TempDir("", "")
	if mkdir_err != nil {
		panic(mkdir_err)
	}
	defer os.RemoveAll(dir)
	packageFile := filepath.Join(dir, "package.json")
	fileWriteErr := ioutil.WriteFile(packageFile, []byte(SAMPLE_PACKAGE_CONTENT), 0644)
	if fileWriteErr != nil {
		panic(fileWriteErr)
	}
	updateErr := UpdatePackage(dir, "bunyan", "1.8.2")
	assert.Nil(t, updateErr)
	bytes, err := ioutil.ReadFile(packageFile)
	if err != nil {
		panic(err)
	}
	updatedPackageContent := string(bytes)
	assert.Equal(t, EXPECTED_UPDATED_PACKAGE_CONTENT, updatedPackageContent)
}
