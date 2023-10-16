package api

import (
	"Wasa-Photo-1894389/service/api/reqcontext"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Function that deletes a photo (this includes comments and likes)
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Estraggo il token di autenticazione bearer dall'header "Authorization"
	// della richiesta HTTP e l'ID della foto dal parametro del percorso.
	bearerAuth := extractBearer(r.Header.Get("Authorization"))
	photoIdStr := ps.ByName("photo_id")

	// Verifica se l'utente che effettua la richiesta è valido
	valid := validateRequestingUser(ps.ByName("id"), bearerAuth)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Converte l'ID della foto da una stringa a un intero a 64 bit
	photoInt, err := strconv.ParseInt(photoIdStr, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-delete/ParseInt: error converting photoId to int")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Chiama una funzione per rimuovere la foto dal database.
	err = rt.db.RemovePhoto(
		User{IdUser: bearerAuth}.ToDatabase(),
		PhotoId{IdPhoto: photoInt}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-delete/RemovePhoto: error coming from database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Ottiene il percorso della cartella che contiene la foto da eliminare.
	pathPhoto, err := getUserPhotoFolder(bearerAuth)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-delete/getUserPhotoFolder: error with directories")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Rimuovo il file dalla cartella delle foto
	err = os.Remove(filepath.Join(pathPhoto, photoIdStr))
	if err != nil {
		// Error occurs if the file doesn't exist, but for idempotency an error won't be raised
		ctx.Logger.WithError(err).Error("photo-delete/os.Remove: photo to be removed is missing")
	}

	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}
