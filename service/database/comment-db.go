package database

// Questa funzione recupera la lista completa dei commenti di una foto, escludendo i commenti degli utenti che hanno bannato l'utente richiedente.
func (db *appdbimpl) GetCompleteCommentsList(requestingUser User, requestedUser User, photo PhotoId) ([]CompleteComment, error) {
	// Utilizza una query SQL per selezionare tutti i commenti della foto specificata,
	// escludendo quelli degli utenti che hanno bannato l'utente richiedente o l'utente richiesto.
	rows, err := db.c.Query("SELECT * FROM comments WHERE id_photo = ? AND id_user NOT IN (SELECT banned FROM banned_users WHERE banner = ? OR banner = ?) "+
		"AND id_user NOT IN (SELECT banner FROM banned_users WHERE banned = ?)",
		photo.IdPhoto, requestingUser.IdUser, requestedUser.IdUser, requestingUser.IdUser)
	if err != nil {
		return nil, err
	}

	// Wait for the function to finish before closing rows
	defer func() { _ = rows.Close() }()

	// Read all the comments in the resulset (comments of the photo with authors that didn't ban the requesting user).
	var comments []CompleteComment
	for rows.Next() {
		var comment CompleteComment
		err = rows.Scan(&comment.IdComment, &comment.IdPhoto, &comment.IdUser, &comment.Comment)
		if err != nil {
			return nil, err
		}

		// Get the nickname of the user that commented
		nickname, err := db.GetNickname(User{IdUser: comment.IdUser})
		if err != nil {
			return nil, err
		}
		comment.Nickname = nickname

		comments = append(comments, comment)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return comments, nil
}

// Database function per aggiungere un commento di un user ad una foto
func (db *appdbimpl) CommentPhoto(p PhotoId, u User, c Comment) (int64, error) {
	//Utilizza una query SQL INSERT per inserire il commento nel database.
	res, err := db.c.Exec("INSERT INTO comments (id_photo,id_user,comment) VALUES (?, ?, ?)",
		p.IdPhoto, u.IdUser, c.Comment)
	if err != nil {
		// Error executing query
		return -1, err
	}

	commentId, err := res.LastInsertId()
	if err != nil {
		// Error getting id returned by last db operation (commentId)
		return -1, err
	}
	// Restituisce l'ID del commento appena inserito.
	return commentId, nil
}

/*
Technically, given the structure of the db, it wouldn't be necessary to have the
id_user to remove a comment, but it is used to make sure that whoever is requesting
the removal is the author of the latter.
Similarly for the id_photo part, except this time we want to make sure that if the url
is not valid but that comment exists for the given user, it won't be deleted.
*/

// Database function che rimuove il commento di un utente dalla foto
func (db *appdbimpl) UncommentPhoto(p PhotoId, u User, c CommentId) error {
	// Utilizza una query SQL DELETE per rimuovere il commento specificato dal database.
	_, err := db.c.Exec("DELETE FROM comments WHERE (id_photo = ? AND id_user = ? AND id_comment = ?)",
		p.IdPhoto, u.IdUser, c.IdComment)
	if err != nil {
		return err
	}

	return nil
}

// data base function permette all'autore di una foto di rimuovere un commento dalla sua foto.
// A differenza della funzione UncommentPhoto, questa funzione non richiede l'ID dell'utente,
// poiché si presume che l'autore della foto abbia il diritto di rimuovere qualsiasi commento dalla sua foto.
func (db *appdbimpl) UncommentPhotoAuthor(p PhotoId, c CommentId) error {
	// Utilizza una query SQL DELETE per rimuovere il commento specificato dal database.
	_, err := db.c.Exec("DELETE FROM comments WHERE (id_photo = ? AND id_comment = ?)",
		p.IdPhoto, c.IdComment)
	if err != nil {
		return err
	}

	return nil
}
