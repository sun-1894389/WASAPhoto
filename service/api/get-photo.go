package api

import (
	"Wasa-Photo-1894389/service/api/reqcontext"
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
)

// Funzione che restituisce la foto richiesta
// Vengono estratti user_id,photo_id si crea il percorso
// e infine viene servito il file con il metodo http
func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	http.ServeFile(w, r,
		filepath.Join(photoFolder, ps.ByName("id"), "photos", ps.ByName("photo_id")))

}
