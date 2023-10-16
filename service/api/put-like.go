package api

import (
	"Wasa-Photo-1894389/service/api/reqcontext"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Funzione permette a un utente di mettere like ad una foto 
func (rt *_router) putLike(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Estraggo l'ID dell'autore della foto, l'ID dell'utente che sta effettuando la richiesta 
	// e l'ID del "like" dai parametri della richiesta HTTP.
	photoAuthor := ps.ByName("id")
	requestingUserId := extractBearer(r.Header.Get("Authorization"))
	pathLikeId := ps.ByName("like_id")

	// Controlla se l'user è loggato
	if isNotLogged(requestingUserId) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// User is trying to like his/her photo
	if photoAuthor == requestingUserId {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the requesting user wasn't banned by the photo owner
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

	// Controlla se l'ID del "like" nella richiesta corrisponde all'ID dell'utente che effettua la richiesta. 
	// Se non corrisponde, invia una risposta con stato "Bad Request" e termina la funzione.
	if pathLikeId != requestingUserId {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Converte l'ID della foto da stringa a int64.(in caso di errore:500)
	photo_id_64, err := strconv.ParseInt(ps.ByName("photo_id"), 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("put-like: error converting path param photo_id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Chiama una funzione del database per aggiungere il "like" alla foto specificata.
	err = rt.db.LikePhoto(
		PhotoId{IdPhoto: photo_id_64}.ToDatabase(),
		User{IdUser: pathLikeId}.ToDatabase())
	if err != nil {
		// ctx.Logger.WithError(err).Error("put-like: error executing insert query")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}
