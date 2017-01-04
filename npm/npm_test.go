package npm

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"io/ioutil"
	"path/filepath"
	"os"
	"os/exec"
	"log"
	"encoding/json"
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

func TestUpdateDependencySuccess(t *testing.T) {
	updated, err := UpdateDependency(SAMPLE_PACKAGE_CONTENT, "bunyan", "1.8.2")
	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_UPDATED_PACKAGE_CONTENT, updated)
}

func TestUpdateDependencyNotFound(t *testing.T) {
	updated, err := UpdateDependency(SAMPLE_PACKAGE_CONTENT, "not-a-dependency", "1.8.2")
	assert.Equal(t, "", updated)
	assert.NotNil(t, err)
	errMessage := err.Error()
	assert.Contains(t, errMessage, "not-a-dependency")
	assert.Contains(t, errMessage, "not found")
}

const EXPECTED_UPDATED_PACKAGE_CONTENT_MULTIPLE_DEPS = `{
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
    "@chauffeur-prive/i18n": "1.2.0",
    "accept-language-parser": "1.3.0",
    "bunyan": "1.8.2"
  },
  "devDependencies": {
    "chai": "3.5.0",
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

func TestUpdateDependencies(t *testing.T) {
	updates := map[string]string{
		"@chauffeur-prive/i18n": "1.2.0",
		"accept-language-parser": "1.3.0",
		"bunyan": "1.8.2",
		"chai": "3.5.0",
		"eslint-config-cp": "transcovo/eslint-config-cp#1.1.0",
	}
	updated, err := UpdateDependencies(SAMPLE_PACKAGE_CONTENT, updates)
	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_UPDATED_PACKAGE_CONTENT_MULTIPLE_DEPS, updated)
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

const SAMPLE_NPM_LIST_OUTPUT = `{
  "name": "c3po",
  "version": "0.1.0",
  "problems": [
    "peer dep missing: eslint-plugin-react@^5.0.1, required by eslint-config-airbnb@9.0.1"
  ],
  "dependencies": {
    "@chauffeur-prive/mongo-helper": {
      "version": "2.1.0",
      "from": "@chauffeur-prive/mongo-helper@>=2.1.0 <3.0.0",
      "resolved": "https://registry.npmjs.org/@chauffeur-prive/mongo-helper/-/mongo-helper-2.1.0.tgz"
    },
    "express-middleware": {
      "version": "1.5.0",
      "from": "express-middleware@>=1.3.2 <2.0.0",
      "resolved": "https://registry.npmjs.org/express-middleware/-/express-middleware-1.5.0.tgz"
    }
  }
}`

func TestParseNpmListOutput(t *testing.T) {
	versions := ParseNpmListOutput([]byte(SAMPLE_NPM_LIST_OUTPUT))

	assert.Len(t, versions, 2)
	assert.Equal(t, "2.1.0", versions["@chauffeur-prive/mongo-helper"])
	assert.Equal(t, "1.5.0", versions["express-middleware"])
}

func TestFreezePackage(t *testing.T) {
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
	updateErr := FreezePackage(dir)
	assert.Nil(t, updateErr)
	bytes, err := ioutil.ReadFile(packageFile)
	if err != nil {
		panic(err)
	}
	updatedPackageContent := string(bytes)
	assert.NotContains(t, updatedPackageContent, "^")
	assert.NotContains(t, updatedPackageContent, "~")
}
