package api

import (
	"Wasa-Photo-1894389/service/api/reqcontext"
	"Wasa-Photo-1894389/service/database"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// This function retrieves all the photos of the people that the user is following
func (rt *_router) getHome(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// imposto il tipo di contenuto della risposta http in json
	// estraggo l'id del'utente dal  token Bearer nell'header di autorizzazione della richiesta HTTP.
	w.Header().Set("Content-Type", "application/json")
	identifier := extractBearer(r.Header.Get("Authorization"))

	// Verifica se è l'utente stesso a vedere la propria home
	valid := validateRequestingUser(ps.ByName("id"), identifier)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Ottengo un elenco di utenti che l'utente sta seguendo dal database.
	followers, err := rt.db.GetFollowing(User{IdUser: identifier}.ToDatabase())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Itera sugli utenti seguiti e recupera le loro foto. Aggiungendo le foto a un elenco.
	var photos []database.Photo
	for _, follower := range followers {

		followerPhoto, err := rt.db.GetPhotosList(
			User{IdUser: identifier}.ToDatabase(),
			User{IdUser: follower.IdUser}.ToDatabase())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for i, photo := range followerPhoto {
			if i >= database.PhotosPerUserHome {
				break
			}
			photos = append(photos, photo)
		}

	}

	// Imposta lo stato della risposta HTTP come 200 OK. Codifica l'elenco di foto in formato JSON e lo invia come corpo della risposta HTTP.
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(photos)
}
