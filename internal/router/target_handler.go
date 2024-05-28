package router

import (
	"fmt"
	"net/http"
)

func ServerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from the target server! You requested: %s", r.URL.Path)
}