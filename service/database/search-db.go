package database

// Database function che permette di cercare utenti in base a un parametro fornito.
// Ogni partial macth viene incluso nei risultati,restituendo una lista di macthing users
func (db *appdbimpl) SearchUser(searcher User, userToSearch User) ([]CompleteUser, error) {
	// La funzione esegue una query SQL SELECT per cercare tutti gli utenti il cui id_user o nickname corrisponde parzialmente 
	// al parametro fornito (userToSearch.IdUser). L'uso del simbolo % dopo userToSearch.IdUser nella query indica una ricerca di tipo "LIKE",
	// che restituirà tutti gli utenti che hanno un id_user o un nickname che inizia con il valore di userToSearch.IdUser.
	rows, err := db.c.Query("SELECT * FROM users WHERE ((id_user LIKE ?) OR (nickname LIKE ?)) AND id_user NOT IN (SELECT banner FROM banned_users WHERE banned = ?)",
		userToSearch.IdUser+"%", userToSearch.IdUser+"%", searcher.IdUser)
	if err != nil {
		return nil, err
	}
	// Wait for the function to finish before closing rows.
	defer func() { _ = rows.Close() }()

	// Read all the users in the resulset.
	var res []CompleteUser
	for rows.Next() {
		var user CompleteUser
		err = rows.Scan(&user.IdUser, &user.Nickname)
		if err != nil {
			return nil, err
		}
		res = append(res, user)
	}
	// Restituisce un elenco di utenti (CompleteUser) che corrispondono al parametro di ricerca fornito e che non hanno bannato l'utente searcher.
    // Se si verifica un errore durante l'esecuzione della query o la lettura dei risultati, viene restituito un errore.
	if rows.Err() != nil {
		return nil, err
	}

	return res, nil
}
