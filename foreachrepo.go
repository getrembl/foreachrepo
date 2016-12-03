package main

import (
	"github.com/transcovo/foreachrepo/foreachrepo"
	"net/http"
)

func main() {
	foreachrepo.GetReposList(http.DefaultClient, "transcovo")
}
