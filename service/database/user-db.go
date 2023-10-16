package database

// Database function che recupera lo "stream" di un utente, che consiste nelle foto delle persone seguite dall'utente.
func (db *appdbimpl) GetStream(user User) ([]Photo, error) {
	// Esegue una query SQL SELECT per selezionare tutte le foto degli utenti seguiti
	// dall'utente specificato e le ordina in base alla data di caricamento in ordine decrescente.
	rows, err := db.c.Query(`SELECT * FROM photos WHERE id_user IN (SELECT followed FROM followers WHERE follower = ?) ORDER BY date DESC`,
		user.IdUser)
	if err != nil {
		return nil, err
	}
	// Wait for the function to finish before closing rows
	defer func() { _ = rows.Close() }()

	// Read all the users in the resulset
	var res []Photo
	for rows.Next() {
		var photo Photo
		err = rows.Scan(&photo.PhotoId, &photo.Owner, &photo.Date) //  &photo.Comments, &photo.Likes,
		if err != nil {
			return nil, err
		}
		res = append(res, photo)
	}

	if rows.Err() != nil {
		return nil, err
	}
	// Restituisce una slice di Photo che rappresenta lo stream dell'utente.
	return res, nil
}

// Database function che aggiunge un nuovo utente al database durante la registrazione.
func (db *appdbimpl) CreateUser(u User) error {
	// Esegue una query SQL per inserire un nuovo utente nel database con un ID utente e un soprannome 
	_, err := db.c.Exec("INSERT INTO users (id_user,nickname) VALUES (?, ?)",
		u.IdUser, u.IdUser)

	if err != nil {
		return err
	}

	return nil
}

// Database function controlla se un utente esiste nel database.
func (db *appdbimpl) CheckUser(targetUser User) (bool, error) {
	//  Esegue una query SQL per contare il numero di righe nella tabella degli utenti che corrispondono all'ID utente specificato.
	var cnt int
	err := db.c.QueryRow("SELECT COUNT(*) FROM users WHERE id_user = ?",
		targetUser.IdUser).Scan(&cnt)

	if err != nil {
		// Count always returns a row thanks to COUNT(*), so this situation should not happen
		return true, err
	}

	// If the counter is 1 then the user exists
	if cnt == 1 {
		return true, nil
	}
	return false, nil
}
