package api

import (
	"Wasa-Photo-1894389/service/api/reqcontext"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Funzione per rimuovere il like da una foto
func (rt *_router) deleteLike(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Estraggo l'ID dell'autore della foto e dell'ID dell'utente che effettua la richiesta dalla richiesta HTTP.
	photoAuthor := ps.ByName("id")
	requestingUserId := extractBearer(r.Header.Get("Authorization"))

	// Controllo se l'user è loggato(se non lo è error:403)
	if isNotLogged(requestingUserId) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Controllo se è l'utente stesso a volersi togliere like(ho messo che non è possibile farlo)
	if photoAuthor == requestingUserId {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Controllo che il richiedente non è stato bannato dal photo owner
	banned, err := rt.db.BannedUserCheck(
		User{IdUser: requestingUserId}.ToDatabase(),
		User{IdUser: photoAuthor}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("post-comment/db.BannedUserCheck: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if banned {
		// User was banned by owner, can't post the comment
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Conversione dell'ID della foto da stringa a int64.
	photoIdInt, err := strconv.ParseInt(ps.ByName("photo_id"), 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("remove-like/ParseInt: error converting photo_id to int64")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Chiamata della funzione UnlikePhoto del database(like-db) per rimuovere il "like".
	err = rt.db.UnlikePhoto(
		PhotoId{IdPhoto: photoIdInt}.ToDatabase(),
		User{IdUser: requestingUserId}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("remove-like/db.UnlikePhoto: error executing insert query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}
