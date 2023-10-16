//go:build !webui

package main

import (
	"net/http"
)

// Questo stub si usa quando non si vuole lanciare il frontend
// registerWebUI is an empty stub because `webui` tag has not been specified.
func registerWebUI(hdl http.Handler) (http.Handler, error) {
	return hdl, nil
}
