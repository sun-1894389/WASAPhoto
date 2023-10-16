package database

// Database fuction che permette a un utente(banner) di bannarne un'altro(banned)
func (db *appdbimpl) BanUser(banner User, banned User) error {

	_, err := db.c.Exec("INSERT INTO banned_users (banner,banned) VALUES (?, ?)", banner.IdUser, banned.IdUser)
	if err != nil {
		return err
	}

	return nil
}

// Database fuction che rimuovere un utente(banned)dalla lista dei banned di un'altro utente(banner)
func (db *appdbimpl) UnbanUser(banner User, banned User) error {

	_, err := db.c.Exec("DELETE FROM banned_users WHERE (banner = ? AND banned = ?)", banner.IdUser, banned.IdUser)
	if err != nil {
		return err
	}

	return nil
}

// [Util] Data base function per controllare se un utente è stato bannato.
// Restituisco 'true' se è banned, sennò 'false'
func (db *appdbimpl) BannedUserCheck(requestingUser User, targetUser User) (bool, error) {
	// Utilizza il metodo QueryRow per eseguire una query SQL SELECT COUNT(*) che 
	// conta quante volte l'utente requestingUser appare nella tabella banned_users come utente bannato da targetUser
	var cnt int
	err := db.c.QueryRow("SELECT COUNT(*) FROM banned_users WHERE banned = ? AND banner = ?",
		requestingUser.IdUser, targetUser.IdUser).Scan(&cnt)

	if err != nil {
		// Count always returns a row thanks to COUNT(*), so this situation should not happen
		return true, err
	}

	// If the counter is 1 then the user was banned
	if cnt == 1 {
		return true, nil
	}
	return false, nil
}
