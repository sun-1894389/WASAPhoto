package api

import (
	"Wasa-Photo-1894389/service/api/reqcontext"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Funzione che rimuove un commento da una foto
func (rt *_router) deleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Imposto l'header "Content-Type" della risposta a "application/json".
	// Estraggo il token di autenticazione dall'header "Authorization" della richiesta HTTP.
	w.Header().Set("Content-Type", "application/json")
	requestingUserId := extractBearer(r.Header.Get("Authorization"))

	// Controllo se l'user è loggato(se non lo è error:403)
	if isNotLogged(requestingUserId) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Verifica se l'utente che effettua la richiesta è stato bandito dal proprietario della foto.
	banned, err := rt.db.BannedUserCheck(
		User{IdUser: requestingUserId}.ToDatabase(),
		User{IdUser: ps.ByName("id")}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("post-comment/db.BannedUserCheck: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if banned {
		// Utente bannato,non posso commentare(403)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Converto l'identificatore della foto da stringa a int64.
	photo_id_64, err := strconv.ParseInt(ps.ByName("photo_id"), 10, 64)
	if err != nil {
		// C'è stato un errore di conversione(error:400)
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("post-comment: failed convert photo_id to int64")
		return
	}

	// Converto l'identificatore del commento da stringa a int64.
	comment_id_64, err := strconv.ParseInt(ps.ByName("comment_id"), 10, 64)
	if err != nil {
		// C'è stato un errore di conversione(error:400)
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("post-comment: failed convert photo_id to int64")
		return
	}

	// Controllo per vedere se l'utente che effettua la richiesta è l'autore del post.
	if ps.ByName("id") == requestingUserId {
		// Chiamo la funzione dal db per rimuovere il commento
		err = rt.db.UncommentPhotoAuthor(
			PhotoId{IdPhoto: photo_id_64}.ToDatabase(),
			CommentId{IdComment: comment_id_64}.ToDatabase())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("post-comment: failed to execute query for insertion")
			return
		}

		// Rispondo con codice 204(la richiesta è andata a buon fine)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Chiamo la funzione dal db per rimuovere il commento(solo l'autore può rimuovere il proprio commento)
	err = rt.db.UncommentPhoto(
		PhotoId{IdPhoto: photo_id_64}.ToDatabase(),
		User{IdUser: requestingUserId}.ToDatabase(),
		CommentId{IdComment: comment_id_64}.ToDatabase())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("post-comment: failed to execute query for insertion")
		return
	}

	// Rispondo con codice 204(la richiesta è andata a buon fine)
	w.WriteHeader(http.StatusNoContent)
}
