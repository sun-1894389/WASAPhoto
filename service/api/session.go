package api

import (
	"Wasa-Photo-1894389/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
)

// Funzione che gestice la sessione degli utenti
func (rt *_router) sessionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Imposta l'intestazione della risposta per indicare che il tipo di contenuto sarà JSON.
	w.Header().Set("Content-Type", "application/json")

	// Inizializza una variabile user e tenta di decodificare il corpo della richiesta in questa variabile.
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	// Controlla se c'è stato un errore durante la decodifica o se l'identificatore dell'utente non è valido.
	// In entrambi i casi, risponde con un codice di stato 400 Bad Request.
	if err != nil {
		// The body was not a parseable JSON, reject it
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if !validIdentifier(user.IdUser) {
		// Here we checked the user identifier and we discovered that it's not valid
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Crea l'utente nel database.
	// Se c'è un errore durante la creazione dell'utente nel database (ad es. l'utente esiste già),
	// risponde con un codice di stato 200 OK e restituisce l'utente.  Se c'è un errore durante la codifica della risposta, 
	// risponde con un codice di stato 500 Internal Server Error.
	err = rt.db.CreateUser(user.ToDatabase())
	if err != nil {
		// l'utente esiste gia,viene restituito l'identifier
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("session: can't create response json")
		}
		return
	}

	// Crea una cartella per l'utente.
	err = createUserFolder(user.IdUser, ctx)

	// Controllo se c'è un errore durante la creazione della cartella
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("session: can't create user's photo folder")
		return
	}

	// Se tutto va bene, risponde con un codice di stato 201 Created e restituisce l'utente
	// Se c'è un errore rispondo con error 500
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("session: can't create response json")
		return
	}
}

// Funzione per creare una nuova cartella per un specifico utente
func createUserFolder(identifier string, ctx reqcontext.RequestContext) error {

	// Creao il path media/useridentifier/ dentro il project dir
	path := filepath.Join(photoFolder, identifier)

	// nel path creato aggiungo un subdir "photos"
	err := os.MkdirAll(filepath.Join(path, "photos"), os.ModePerm)
	if err != nil {
		ctx.Logger.WithError(err).Error("session/createUserFolder:: error creating directories for user")
		return err
	}
	return nil
}
