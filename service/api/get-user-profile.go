package api

import (
	"Wasa-Photo-1894389/service/api/reqcontext"
	"Wasa-Photo-1894389/service/database"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Funzione che ritrova tutte le info necessarie del profilo
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Estraggo l'id dell'utente che fa la richiesta e l'id dell'utente richiesto
	requestingUserId := extractBearer(r.Header.Get("Authorization"))
	requestedUser := ps.ByName("id")

	var followers []database.User
	var following []database.User
	var photos []database.Photo

	// Controlla se l'utente che effettua la richiesta è bannato dall'utente richiesto, utilizzando una funzione del database.
	// Se l'utente che effettua la richiesta è bannato, la funzione restituisce un codice di stato HTTP 403 (Forbidden) e termina.
	userBanned, err := rt.db.BannedUserCheck(User{IdUser: requestingUserId}.ToDatabase(),
		User{IdUser: requestedUser}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile/db.BannedUserCheck/userBanned: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if userBanned {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Controllo se il profilo richiesto è stato bannato dal richiedente. Se si rispondo StatusPartialContent
	requestedProfileBanned, err := rt.db.BannedUserCheck(User{IdUser: requestedUser}.ToDatabase(),
		User{IdUser: requestingUserId}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile/db.BannedUserCheck/requestedProfileBanned: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if requestedProfileBanned {
		w.WriteHeader(http.StatusPartialContent)
		return
	}

	// Controllo se l'utente esiste
	userExists, err := rt.db.CheckUser(User{IdUser: requestedUser}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile/db.CheckUser: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !userExists {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Recupero la lista dei followers dell'utente richiesto
	followers, err = rt.db.GetFollowers(User{IdUser: requestedUser}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile/db.GetFollowers: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Recupero la lista dei utenti seguiti dall'utente richiesto
	following, err = rt.db.GetFollowing(User{IdUser: requestedUser}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile/db.GetFollowing: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Recupera la lista delle foto dell'utente richiesto dal database.
	photos, err = rt.db.GetPhotosList(User{IdUser: requestingUserId}.ToDatabase(), User{IdUser: requestedUser}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile/db.GetPhotosList: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Recupero il nickname dell'utente richiesto
	nickname, err := rt.db.GetNickname(User{IdUser: requestedUser}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile/db.GetNickname: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Imposta il codice di stato della risposta HTTP come 200 (OK) e invia un oggetto
	// JSON che rappresenta il profilo completo dell'utente richiesto come corpo della risposta.
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(CompleteProfile{
		Name:      requestedUser,
		Nickname:  nickname,
		Followers: followers,
		Following: following,
		Posts:     photos,
	})

}
