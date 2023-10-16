package api

import (
	"Wasa-Photo-1894389/service/api/reqcontext"
	"bytes"
	"encoding/json"
	"errors"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

// Funzione che gestisce l'upload di una foto
func (rt *_router) postPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")
	auth := extractBearer(r.Header.Get("Authorization"))

	// Verifica l'identità dell'utente che effettua la richiesta
	valid := validateRequestingUser(ps.ByName("id"), auth)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Inizializza una struttura Photo con l'ID dell'utente e la data corrente.
	photo := Photo{
		Owner: auth,
		Date:  time.Now().UTC(),
	}

	// Legge il body della richiesta e verifica se ci sono errori durante la lettura.
	data, err := io.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-upload: error reading body content")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Reimposta il body della richiesta in modo da poterlo leggere di nuovo in seguito
	// Dopo aver letto il body bisogna riassegnare un io.ReadCloser per poterlo rileggere
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	// verifico se il contenuto del body è una immagine png o jpeg(in caso di errore:400 badrequest)
	err = checkFormatPhoto(r.Body, io.NopCloser(bytes.NewBuffer(data)), ctx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("photo-upload: body contains file that is neither jpg or png")
		// controllaerrore
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: IMG_FORMAT_ERROR_MSG})
		return
	}

	// Reimposta nuovamente il corpo della richiesta per poterlo leggere di nuovo.
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	// Chiama una funzione del database per creare un record per la foto e ottenere un ID univoco per essa
	// Se si verifica un errore, risponde con un codice di stato HTTP 500 (Internal Server Error)
	photoIdInt, err := rt.db.CreatePhoto(photo.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-upload: error executing db function call")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Converte l'ID della foto da int64 a string.
	photoId := strconv.FormatInt(photoIdInt, 10)

	// Ottiene il percorso della cartella dove verrà salvata la foto dell'utente utilizzando la funzione definita sotto.
	PhotoPath, err := getUserPhotoFolder(auth)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-upload: error getting user's photo folder")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Crea un nuovo file con l'ID della foto come nome nel percorso ottenuto.(per salvare il contenuto del body)
	out, err := os.Create(filepath.Join(PhotoPath, photoId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("photo-upload: error creating local photo file")
		return
	}

	// Copia il contenuto del corpo della richiesta nel file appena creato. 
	_, err = io.Copy(out, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("photo-upload: error copying body content into file photo")
		return
	}

	// Chiude il file appena creato
	out.Close()

	// Invia una risposta con stato "Created" e un oggetto JSON che rappresenta la foto appena caricata.
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(Photo{
		Comments: nil,
		Likes:    nil,
		Owner:    photo.Owner,
		Date:     photo.Date,
		PhotoId:  int(photoIdInt),
	})

}

// Funzione per controllare se il formato della foto è png o jpeg.Ritorno l'estenzione del formato e un errore 
func checkFormatPhoto(body io.ReadCloser, newReader io.ReadCloser, ctx reqcontext.RequestContext) error {

	_, errJpg := jpeg.Decode(body)
	if errJpg != nil {

		body = newReader
		_, errPng := png.Decode(body)
		if errPng != nil {
			return errors.New(IMG_FORMAT_ERROR_MSG)
		}
		return nil
	}
	return nil
}

// Funzione che restituice il path del photo folder del'utente
func getUserPhotoFolder(user_id string) (UserPhotoFoldrPath string, err error) {

	// Path of the photo dir "./media/user_id/photos/"
	photoPath := filepath.Join(photoFolder, user_id, "photos")

	return photoPath, nil
}
