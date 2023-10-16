//go:build webui
//

package main

import (
	"Wasa-Photo-1894389/webui"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
)

// questa funzione configura l'applicazione per servire l'interfaccia utente web da una directory embedded
// quando l'URL richiesto inizia con /dashboard/, e per gestire tutte le altre richieste.

func registerWebUI(hdl http.Handler) (http.Handler, error) {
	distDirectory, err := fs.Sub(webui.Dist, "dist")
	if err != nil {
		return nil, fmt.Errorf("error embedding WebUI dist/ directory: %w", err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/dashboard/") {
			http.StripPrefix("/dashboard/", http.FileServer(http.FS(distDirectory))).ServeHTTP(w, r)
			return
		}
		hdl.ServeHTTP(w, r)
	}), nil
}
