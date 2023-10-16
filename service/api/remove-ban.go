package api

import (
	"Wasa-Photo-1894389/service/api/reqcontext"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Funzione che rimuove un user dalla lista dei banned di un altro user
func (rt *_router) deleteBan(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Estraggo il token di autenticazione dall'header "Authorization" della richiesta HTTP,
	// l'ID dell'utente e dell'utente da "sbandire"
	bearerToken := extractBearer(r.Header.Get("Authorization"))
	pathId := ps.ByName("id")
	userToUnban := ps.ByName("banned_id")

	// Verifica dell'identità dell'utente che effettua la richiesta.
	valid := validateRequestingUser(pathId, bearerToken)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Controllo per assicurarsi che un utente non stia cercando di "sbandire" se stesso.(error:204)
	// (Un utente non può bannarsi da solo)
	if userToUnban == bearerToken {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Chiamo la funzione UnbanUser dal database per rimuovere l'utente dalla lista dei banned.
	err := rt.db.UnbanUser(
		User{IdUser: pathId}.ToDatabase(),
		User{IdUser: userToUnban}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("remove-ban: error executing delete query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}
