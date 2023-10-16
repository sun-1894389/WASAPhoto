package database

// Database function che recupera l'elenco delle foto di un utente (targetUser),
// ma solo se l'utente che fa la richiesta (requestingUser) non è stato bannato da targetUser.
func (db *appdbimpl) GetPhotosList(requestingUser User, targetUser User) ([]Photo, error) { // requestinUser User,
	// esegue una query SQL per selezionare tutte le foto di targetUser e le ordina in base alla data in ordine decrescente.
	rows, err := db.c.Query("SELECT * FROM photos WHERE id_user = ? ORDER BY date DESC", targetUser.IdUser)
	if err != nil {
		return nil, err
	}
	// Wait for the function to finish before closing rows
	defer func() { _ = rows.Close() }()

	// Read all the photos in the resulset
	// Per ogni foto, recupera l'elenco completo dei commenti e dei "mi piace" associati a quella foto.
	var photos []Photo
	for rows.Next() {
		var photo Photo
		err = rows.Scan(&photo.PhotoId, &photo.Owner, &photo.Date)
		if err != nil {
			return nil, err
		}

		comments, err := db.GetCompleteCommentsList(requestingUser, targetUser, PhotoId{IdPhoto: int64(photo.PhotoId)}) // Old: GetCommentsLen
		if err != nil {
			return nil, err
		}
		photo.Comments = comments

		likes, err := db.GetLikesList(requestingUser, targetUser, PhotoId{IdPhoto: int64(photo.PhotoId)}) // Old: GetLikesLen
		if err != nil {
			return nil, err
		}
		photo.Likes = likes

		photos = append(photos, photo)
	}

	if rows.Err() != nil {
		return nil, err
	}
	// Restituisce un elenco di foto con i relativi commenti e "mi piace".
	return photos, nil
}

// Database function che recupera una foto specifica (targetPhoto),
// ma solo se l'utente che fa la richiesta (requestinUser) non è stato bannato dal proprietario della foto.
func (db *appdbimpl) GetPhoto(requestinUser User, targetPhoto PhotoId) (Photo, error) {
	// Utilizza una query SQL SELECT per recuperare la foto nel database.
	var photo Photo
	err := db.c.QueryRow("SELECT * FROM photos WHERE (id_photo = ?) AND id_user NOT IN (SELECT banner FROM banned_user WHERE banned = ?)",
		targetPhoto.IdPhoto, requestinUser.IdUser).Scan(&photo)

	if err != nil {
		return Photo{}, ErrUserBanned
	}

	return photo, nil

}

// Database function che crea una nuova foto nel database e restituisce l'ID univoco della foto.
func (db *appdbimpl) CreatePhoto(p Photo) (int64, error) {
	// Utilizza una query SQL INSERT per inserire la foto nel database.
	res, err := db.c.Exec("INSERT INTO photos (id_user,date) VALUES (?,?)",
		p.Owner, p.Date)

	if err != nil {
		// Error executing query
		return -1, err
	}

	photoId, err := res.LastInsertId()
	if err != nil {
		// Error getting id returned by last db operation (photoId)
		return -1, err
	}

	return photoId, nil
}

/*
Adding the owner is an additional security measure to delete photos that are actually owned
by that user
*/

// Database function rimuove una foto specifica (p) dal database, ma solo se l'utente specificato (owner) è il proprietario della foto.
func (db *appdbimpl) RemovePhoto(owner User, p PhotoId) error {

	_, err := db.c.Exec("DELETE FROM photos WHERE id_user = ? AND id_photo = ? ",
		owner.IdUser, p.IdPhoto)
	if err != nil {
		// Error during the execution of the query
		return err
	}

	return nil
}

// [Util] Database function che verifica se una foto specifica (targetPhoto) esiste nel database.
func (db *appdbimpl) CheckPhotoExistence(targetPhoto PhotoId) (bool, error) {
	// Utilizza una query SQL SELECT COUNT(*) per contare quante volte l'ID della foto appare nel database.
	var rows int
	err := db.c.QueryRow("SELECT COUNT(*) FROM photos WHERE (id_photo = ?)", targetPhoto.IdPhoto).Scan(&rows)
	if err != nil {
		return false, err
	}
	// Restituisce true se la foto esiste, altrimenti false.
	if rows == 0 {
		return false, nil
	}
	return true, nil

}
