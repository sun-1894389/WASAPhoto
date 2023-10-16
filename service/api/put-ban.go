package api

import (
	"Wasa-Photo-1894389/service/api/reqcontext"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Funzione che gestisce l'aggiunta di un utente alla lista degli utenti bannati di un altro utente.
func (rt *_router) putBan(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Estrae l'ID dell'utente che sta effettuando la richiesta, l'ID dell'utente da bannare e l'ID dell'utente 
	// che vuole bannare un altro utente dalle intestazioni e dai parametri della richiesta HTTP.
	pathId := ps.ByName("id")
	pathBannedId := ps.ByName("banned_id")
	requestingUserId := extractBearer(r.Header.Get("Authorization"))

	// Controlla se l'utente che effettua la richiesta ha il permesso di bannare l'utente specificato. 
	// (solo l'owner dell'account puo aggiugnere un banned user nel suo account list)
	valid := validateRequestingUser(pathId, requestingUserId)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Controlla se l'utente sta cercando di bannare se stesso.se si(400:"Bad Request")
	if requestingUserId == pathBannedId {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Chiama una funzione del database per aggiungere l'utente specificato alla lista degli utenti bannati.
	err := rt.db.BanUser(
		User{IdUser: pathId}.ToDatabase(),
		User{IdUser: pathBannedId}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("put-ban/db.BanUser: error executing insert query")

		// C'è stato un errore interno,restituisco(error:500)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Bannare implica anche rimuovere il follow (se ci sta)
	err = rt.db.UnfollowUser(
		User{IdUser: requestingUserId}.ToDatabase(),
		User{IdUser: pathBannedId}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("put-ban/db.UnfollowUser1: error executing insert query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// L'utente bannato viene rimosso l'utente che l'ha bannato
	err = rt.db.UnfollowUser(
		User{IdUser: pathBannedId}.ToDatabase(),
		User{IdUser: requestingUserId}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("put-ban/db.UnfollowUser2: error executing insert query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}
