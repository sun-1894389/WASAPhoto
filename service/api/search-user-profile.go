package api

import (
	"Wasa-Photo-1894389/service/api/reqcontext"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Funzione che recupera tutti gli utenti corrispondenti al parametro di query e invia la risposta contenente tutte le corrispondenze.
func (rt *_router) getUsersQuery(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Imposta l'intestazione della risposta per indicare che il tipo di contenuto sarà JSON.
	w.Header().Set("Content-Type", "application/json")

	// Estraggo l'identificatore dell'utente dal token Bearer nella richiesta.
	identifier := extractBearer(r.Header.Get("Authorization"))

	// Se non sono loggato rispondo con error:403
	if identifier == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Estraggo il parametro di query "id" dall'URL della richiesta.
	identificator := r.URL.Query().Get("id")

	// Ricerca dell'utente nel database utilizzando il parametro di query come filtro.
	res, err := rt.db.SearchUser(User{IdUser: identifier}.ToDatabase(), User{IdUser: identificator}.ToDatabase())
	if err != nil {
		// C'è stato un errore, Return an empty json(error 500)
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Database has encountered an error")
		// controllaerrore
		_ = json.NewEncoder(w).Encode([]User{})
		return
	}

	// Imposta un codice di stato 200 OK per la risposta, indicando che tutto è andato bene fino a questo punto.
	w.WriteHeader(http.StatusOK)

	// Invia l'output all'utente. Se non ci sono corrispondenze, invece di restituire un valore null,
	// viene restituito un array JSON vuoto. Altrimenti, restituisce l'array di utenti corrispondenti.
	if len(res) == 0 {
		_ = json.NewEncoder(w).Encode([]User{})
		return
	}
	_ = json.NewEncoder(w).Encode(res)
}
