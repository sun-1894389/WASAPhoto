package database

// Database function che recupera la lista degli utenti che hanno messo "mi piace" a una determinata foto.
func (db *appdbimpl) GetLikesList(requestingUser User, requestedUser User, photo PhotoId) ([]CompleteUser, error) {
	// La query SQL seleziona tutti gli utenti che hanno messo "mi piace" alla foto specificata, escludendo gli utenti
	// che hanno bannato l'utente che fa la richiesta (requestingUser) o sono stati bannati da lui.
	rows, err := db.c.Query("SELECT id_user FROM likes WHERE id_photo = ? AND id_user NOT IN (SELECT banned FROM banned_users WHERE banner = ? OR banner = ?) "+
		"AND id_user NOT IN (SELECT banner FROM banned_users WHERE banned = ?)",
		photo.IdPhoto, requestingUser.IdUser, requestedUser.IdUser, requestingUser.IdUser)
	if err != nil {
		return nil, err
	}
	// Wait for the function to finish before closing rows.
	defer func() { _ = rows.Close() }()

	// Read all the users in the resulset (users that liked the photo that didn't ban the requesting user).
	var likes []CompleteUser
	for rows.Next() {
		var user CompleteUser
		err = rows.Scan(&user.IdUser)
		if err != nil {
			return nil, err
		}

		// Get the nickname of the user that liked the photo
		nickname, err := db.GetNickname(User{IdUser: user.IdUser})
		if err != nil {
			return nil, err
		}
		user.Nickname = nickname

		likes = append(likes, user)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return likes, nil
}

// Database function che  permette a un utente di mettere "mi piace" a una foto.
func (db *appdbimpl) LikePhoto(p PhotoId, u User) error {
	// Utilizza una query SQL INSERT per aggiungere un record nella tabella likes, indicando che l'utente u ha messo "mi piace" alla foto p.
	_, err := db.c.Exec("INSERT INTO likes (id_photo,id_user) VALUES (?, ?)", p.IdPhoto, u.IdUser)
	if err != nil {
		return err
	}

	return nil
}

// Database function che permette a un utente (u) di rimuovere il "mi piace" da una foto (p).
func (db *appdbimpl) UnlikePhoto(p PhotoId, u User) error {
	// Utilizza una query SQL DELETE per rimuovere un record dalla tabella likes, indicando che l'utente u ha rimosso il "mi piace" dalla foto p.
	_, err := db.c.Exec("DELETE FROM likes WHERE(id_photo = ? AND id_user = ?)", p.IdPhoto, u.IdUser)
	if err != nil {
		return err
	}

	return nil
}
