package database

// Database function that retrieves the list of followers of a user
func (db *appdbimpl) GetFollowers(requestingUser User) ([]User, error) {
	rows, err := db.c.Query("SELECT follower FROM followers WHERE followed = ?", requestingUser.IdUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []User
	for rows.Next() {
		var follower User
		err = rows.Scan(&follower.IdUser)
		if err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return followers, nil
}

// Database function that retrieves the list of users followed by the user
func (db *appdbimpl) GetFollowing(requestingUser User) ([]User, error) {
	rows, err := db.c.Query("SELECT followed FROM followers WHERE follower = ?", requestingUser.IdUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var following []User
	for rows.Next() {
		var followed User
		err = rows.Scan(&followed.IdUser)
		if err != nil {
			return nil, err
		}
		following = append(following, followed)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return following, nil
}

// Database function that adds a follower to a user
func (db *appdbimpl) FollowUser(follower User, followed User) error {
	_, err := db.c.Exec("INSERT INTO followers (follower, followed) VALUES (?, ?)", follower.IdUser, followed.IdUser)
	if err != nil {
		return err
	}
	return nil
}

// Database function that removes a follower from a user
func (db *appdbimpl) UnfollowUser(follower User, followed User) error {
	_, err := db.c.Exec("DELETE FROM followers WHERE follower = ? AND followed = ?", follower.IdUser, followed.IdUser)
	if err != nil {
		return err
	}
	return nil
}
