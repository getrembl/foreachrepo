package npm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const samplePackageContent = `{
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

func TestUpdateDependencyVersionSuccess(t *testing.T) {
	updated, err := UpdateDependencyVersion(samplePackageContent, "bunyan", "1.8.2")
	assert.Nil(t, err)
	expectedUpdated := `{
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
	assert.Equal(t, expectedUpdated, updated)
}

func TestUpdateDependencyVersionNotFound(t *testing.T) {
	updated, err := UpdateDependencyVersion(samplePackageContent, "not-a-dependency", "1.8.2")
	assert.Equal(t, "", updated)
	assert.NotNil(t, err)
	errMessage := err.Error()
	assert.Contains(t, errMessage, "not-a-dependency")
	assert.Contains(t, errMessage, "not found")
}

func TestUpdateDependencyVersionUpToDate(t *testing.T) {
	updated, err := UpdateDependencyVersion(samplePackageContent, "bunyan", "~1.8.1")
	assert.Equal(t, "", updated)
	assert.NotNil(t, err)
	errMessage := err.Error()
	assert.Contains(t, errMessage, "bunyan")
	assert.Contains(t, errMessage, "up to date")
}
