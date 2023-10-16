package api

import (
	"Wasa-Photo-1894389/service/api/reqcontext"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Funzione per cambiare nickname(gestisce le richieste http put)
func (rt *_router) putNickname(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	pathId := ps.ByName("id")

	// prendo l'ID dell'utente dal percorso della richiesta.
	// Verifica l'identità dell'utente per l'operazione, confrontando l'ID dell'utente con l'ID dell'utente nel token Bearer.
	valid := validateRequestingUser(pathId, extractBearer(r.Header.Get("Authorization")))
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Estraggo il nuovo nickname dal corpo della richiesta e decodifica del JSON
	var nick Nickname
	err := json.NewDecoder(r.Body).Decode(&nick)
	// Se c'è un errore nella decodifica del JSON, si risponde con un codice di stato HTTP 400 (Bad Request).
	if err != nil {
		ctx.Logger.WithError(err).Error("update-nickname: error decoding json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Modifico il nickname dell'utente nel database,usando una funzione del db(nickname.db).
	err = rt.db.ModifyNickname(
		User{IdUser: pathId}.ToDatabase(),
		nick.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("update-nickname: error executing update query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}
