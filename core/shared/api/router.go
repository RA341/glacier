package api

import (
	"fmt"
	"net/http"
	"strings"
)

func WithSubRouter(parent *http.ServeMux, path string, child http.Handler) {
	if strings.HasSuffix(path, "/") {
		panic(fmt.Sprintf("path must not end with /: %s", path))
	}

	basepath := path + "/"
	parent.Handle(
		basepath,
		http.StripPrefix(path, child),
	)
}
