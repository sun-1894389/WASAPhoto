package api

import (
	"Wasa-Photo-1894389/service/api/reqcontext"
	"Wasa-Photo-1894389/service/database"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Funzione per mettere nella lista dei follow di un utente il follow di un'altro utente
func (rt *_router) putFollow(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Estrae l'ID dell'utente da seguire e l'ID dell'utente che sta effettuando 
	// la richiesta dalle intestazioni e dai parametri della richiesta HTTP.
	userToFollowId := ps.ByName("id")
	requestingUserId := extractBearer(r.Header.Get("Authorization"))

	// un utente non si puo followare da solo (error 404)
	if requestingUserId == userToFollowId {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Controlla se l'ID del follower nella richiesta corrisponde all'ID dell'utente che effettua la richiesta.
	if ps.ByName("follower_id") != requestingUserId {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Controllo se chi mette follow non è stato bannato dal followato
	banned, err := rt.db.BannedUserCheck(
		database.User{IdUser: requestingUserId},
		database.User{IdUser: userToFollowId})
	if err != nil {
		ctx.Logger.WithError(err).Error("post-comment/rt.db.BannedUserCheck: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if banned {
		// User was banned, can't perform the follow action
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Viene aggiunto il follower usando la funzione dal database
	err = rt.db.FollowUser(
		User{IdUser: requestingUserId}.ToDatabase(),
		User{IdUser: userToFollowId}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("put-follow: error executing insert query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}
