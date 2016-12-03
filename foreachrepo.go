package main

import (
	"github.com/transcovo/foreachrepo/repos_explorator"
	"net/http"
)

func main() {
	repos_explorator.GetReposList(http.DefaultClient, "transcovo")
}
