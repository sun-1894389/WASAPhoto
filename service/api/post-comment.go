package api

import (
	"Wasa-Photo-1894389/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// funzione che gestisce l'aggiunta di un commento a una foto e invia una risposta contenente l'ID univoco del commento creato.
func (rt *_router) postComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	//Imposta l'intestazione "Content-Type" della risposta HTTP su "application/json", indicando che la risposta sarà in formato JSON.
	//Estrae l'ID del proprietario della foto e l'ID dell'utente che sta effettuando la richiesta dalle intestazioni della richiesta HTTP.
	w.Header().Set("Content-Type", "application/json")
	photoOwnerId := ps.ByName("id")
	requestingUserId := extractBearer(r.Header.Get("Authorization"))

	//Controlla se l'utente che effettua la richiesta non è autenticato.(errore:403)
	if isNotLogged(requestingUserId) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Controlla se l'utente che effettua la richiesta è stato bannato dal proprietario della foto.
	banned, err := rt.db.BannedUserCheck(
		User{IdUser: requestingUserId}.ToDatabase(),
		User{IdUser: photoOwnerId}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("post-comment/db.BannedUserCheck: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if banned {
		// l'utente bannato non puo commentare
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Copio il body content(comment mandato dal'user) nel comment(Struct)
	var comment Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("post-comment/Decode: failed to decode request body json")
		return
	}

	// Controllo la lunghezza del comment(<=30)
	if len(comment.Comment) > 30 {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("post-comment: comment longer than 30 characters")
		return
	}

	// Convertisco l'id della foto da string a int 64 bit
	photo_id_64, err := strconv.ParseInt(ps.ByName("photo_id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("post-comment/ParseInt: failed convert photo_id to int64")
		return
	}

	// Chiama una funzione del database per creare il commento.
	commentId, err := rt.db.CommentPhoto(
		PhotoId{IdPhoto: photo_id_64}.ToDatabase(),
		User{IdUser: requestingUserId}.ToDatabase(),
		comment.ToDatabase())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("post-comment/db.CommentPhoto: failed to execute query for insertion")
		return
	}

	// Imposta lo stato della risposta su "Created".
	w.WriteHeader(http.StatusCreated)

	// Codifica l'ID univoco del commento creato in formato JSON e lo invia come corpo della risposta.
	err = json.NewEncoder(w).Encode(CommentId{IdComment: commentId})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("post-comment/Encode: failed convert photo_id to int64")
		return
	}
}
