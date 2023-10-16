package database

// Database function che recupera la lista degli utenti che seguono l'utente specificato 
func (db *appdbimpl) GetFollowers(requestinUser User) ([]User, error) {
	// Utilizza una query SQL per selezionare tutti gli utenti che seguono l'utente specificato.
	rows, err := db.c.Query("SELECT follower FROM followers WHERE followed = ?", requestinUser.IdUser)
	if err != nil {
		return nil, err
	}
	// Wait for the function to finish before closing rows.
	defer func() { _ = rows.Close() }()

	// Read all the users in the resulset (users that follow the requesting user)
	var followers []User
	for rows.Next() {
		var folower User
		err = rows.Scan(&folower.IdUser)
		if err != nil {
			return nil, err
		}
		followers = append(followers, folower)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return followers, nil
}

// Database function che recupera la lista degli utenti seguiti dall'utente specificato (requestinUser).
func (db *appdbimpl) GetFollowing(requestinUser User) ([]User, error) {
	// Utilizza una query SQL per selezionare tutti gli utenti seguiti dall'utente specificato.
	rows, err := db.c.Query("SELECT followed FROM followers WHERE follower = ?", requestinUser.IdUser)
	if err != nil {
		return nil, err
	}
	// Wait for the function to finish before closing rows.
	defer func() { _ = rows.Close() }()

	// Read all the users in the resulset (users followed by the requesting user)
	var following []User
	for rows.Next() {
		var folowed User
		err = rows.Scan(&folowed.IdUser)
		if err != nil {
			return nil, err
		}
		following = append(following, folowed)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return following, nil
}

// Database function che permette a un utente (follower) di seguire un altro utente (followed).
func (db *appdbimpl) FollowUser(follower User, followed User) error {
	// Utilizza una query SQL INSERT per aggiungere un record nella tabella followers, indicando che l'utente follower segue l'utente followed.
	_, err := db.c.Exec("INSERT INTO followers (follower,followed) VALUES (?, ?)",
		follower.IdUser, followed.IdUser)
	if err != nil {
		return err
	}

	return nil
}

// Database function che permette a un utente (follower) di smettere di seguire un altro utente (followed).
func (db *appdbimpl) UnfollowUser(follower User, followed User) error {
	// Utilizza una query SQL DELETE per rimuovere un record dalla tabella followers,
	// indicando che l'utente follower non segue più l'utente followed.
	_, err := db.c.Exec("DELETE FROM followers WHERE(follower = ? AND followed = ?)",
		follower.IdUser, followed.IdUser)
	if err != nil {
		return err
	}

	return nil
}
