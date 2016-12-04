[![CircleCI](https://circleci.com/gh/transcovo/foreachrepo.svg?style=shield)](https://circleci.com/gh/transcovo/foreachrepo)
[![Coverage Status](https://img.shields.io/codecov/c/github/transcovo/foreachrepo.svg)](https://codecov.io/gh/transcovo/foreachrepo)

## Overview

This project is a small utility to do a simple task on many repos. It will create a pull request on all 
the repos where the task needs to be done.

Tasks include:

- Bump an npm dependency wherever it's declared

## Installation

#### 1. Initial setup for go develomement

First, install go and git

```
apt install -qy golang git
```

Then, setup your workspace

```
mkdir ~/go
```

Run and add to your shell's rc file:

```
export GOPATH=~/go
export PATH=$GOPATH/bin:$PATH
```

#### 2. Install the app

First, install the go app (you will need to repeat this when new features are released)

```
go get github.com/transcovo/foreachrepo
go install github.com/transcovo/foreachrepo
```

Configure your Github credentials. The app will need them to use the Github API

Run and add this to your shell's rc file:
```
export GITHUB_USERNAME=<username>
export GITHUB_PASSWORD=<password>
```

You're all set!

## Examples

#### Bump an npm dependency in all projects

The example below will iterate on all projects of `transcovo`, and for every project that has
a `package.json` in which `chpr-metrics` appears, it will create a pull request to set
the version spec to `1.0.0`.

The branch name will be `fixed-chpr-metrics-version` and the commit message will be
`TECH Use fixed version for chpr-metrics`

```
foreachrepo -org transcovo \
            -npm-dep chpr-metrics \
            -npm-dep-ver 1.0.0` \
            -branch fixed-chpr-metrics-version \
            -message "TECH Use fixed version for chpr-metrics"
```
