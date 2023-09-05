package api

import (
	"Wasa-Photo-1894389/service/api/reqcontext"
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
)

// Function that serves the requested photo
func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	http.ServeFile(w, r,
		filepath.Join(photoFolder, ps.ByName("id"), "photos", ps.ByName("photo_id")))

}
