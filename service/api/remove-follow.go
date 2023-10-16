package api

import (
	"Wasa-Photo-1894389/service/api/reqcontext"
	"Wasa-Photo-1894389/service/database"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Funzione per rimuovere il follow di un utente dalla lista di un'altro utente
func (rt *_router) deleteFollow(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Estraggo il token di autenticazione (Bearer token) dall'header "Authorization" della richiesta HTTP.
    // Estraggo l'ID del vecchio follower e dell'ID del proprietario della foto dai parametri del percorso.
	requestingUserId := extractBearer(r.Header.Get("Authorization"))
	oldFollower := ps.ByName("follower_id")
	photoOwnerId := ps.ByName("id")

	// Controllo per assicurarmi che l'ID del follower nel percorso sia lo stesso del token Bearer.
	valid := validateRequestingUser(oldFollower, requestingUserId)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Controllo per vedere se l'utente sta cercando di "unfollow" se stesso.(Se si,rispondo 204 e termino)
	if photoOwnerId == requestingUserId {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Verifica se l'utente che effettua la richiesta è stato bannato dal proprietario della foto.
	banned, err := rt.db.BannedUserCheck(
		database.User{IdUser: requestingUserId},
		database.User{IdUser: photoOwnerId})
	if err != nil {
		ctx.Logger.WithError(err).Error("post-comment/rt.db.BannedUserCheck: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if banned {
		// User was banned, can't perform the follow action(403)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Chiamata alla funzione UnfollowUser(follow-db) del database per rimuovere il follower.
	err = rt.db.UnfollowUser(
		User{IdUser: oldFollower}.ToDatabase(),
		User{IdUser: photoOwnerId}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("remove-follow: error executing delete query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}
