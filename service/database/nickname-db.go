package database

// Database function che restituice il nickname di un utente
func (db *appdbimpl) GetNickname(user User) (string, error) {

	var nickname string

	// Utilizza una query SQL SELECT per cercare il nickname dell'utente nella tabella users utilizzando l'identificativo dell'utente (id_user).
	err := db.c.QueryRow(`SELECT nickname FROM users WHERE id_user = ?`, user.IdUser).Scan(&nickname)
	if err != nil {
		// Error during the execution of the query
		return "", err
	}
	return nickname, nil
}

// Database function che permette di modificare il nickname di un utente (user) con un nuovo nickname (newNickname).
func (db *appdbimpl) ModifyNickname(user User, newNickname Nickname) error {
	// Utilizza una query SQL UPDATE per modificare il nickname dell'utente nella tabella users utilizzando l'identificativo dell'utente (id_user).
	_, err := db.c.Exec(`UPDATE users SET nickname = ? WHERE id_user = ?`, newNickname.Nickname, user.IdUser)
	if err != nil {
		// Error during the execution of the query
		return err
	}
	return nil
}
