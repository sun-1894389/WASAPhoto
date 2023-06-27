package database

// Database function that retrieves the list of comments of a photo (excluding comments from users that banned the requesting user)
func (db *appdbimpl) GetCompleteCommentsList(requestingUser User, requestedUser User, photo PhotoID) ([]CompleteComment, error) {
	query := `
		SELECT c.id_comment, c.id_photo, c.id_user, c.comment, u.nickname
		FROM comments AS c
		JOIN users AS u ON c.id_user = u.id_user
		WHERE c.id_photo = ? AND c.id_user NOT IN (
			SELECT banned FROM banned_users WHERE banner = ? OR banner = ?
		) AND c.id_user NOT IN (
			SELECT banner FROM banned_users WHERE banned = ?
		)
	`
	rows, err := db.c.Query(query, photo, requestingUser.IdUser, requestedUser.IdUser, requestingUser.IdUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []CompleteComment
	for rows.Next() {
		var comment CompleteComment
		err = rows.Scan(&comment.IdComment, &comment.IdPhoto, &comment.IdUser, &comment.Comment, &comment.Nickname)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return comments, nil
}

// Database function that adds a comment from a user to a photo
func (db *appdbimpl) CommentPhoto(photoID PhotoID, user User, comment Comment) (int64, error) {
	result, err := db.c.Exec("INSERT INTO comments (id_photo, id_user, comment) VALUES (?, ?, ?)", photoID, user.IdUser, comment)
	if err != nil {
		return -1, err
	}
	commentID, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return commentID, nil
}

// Database function that removes a comment from a photo
func (db *appdbimpl) UncommentPhoto(photoID PhotoID, user User, commentID CommentID) error {
	_, err := db.c.Exec("DELETE FROM comments WHERE id_photo = ? AND id_user = ? AND id_comment = ?", photoID, user.IdUser, commentID)
	if err != nil {
		return err
	}
	return nil
}

// Database function that allows the author of a photo to remove a comment from another user on their photo
func (db *appdbimpl) UncommentPhotoAuthor(photoID PhotoID, commentID CommentID) error {
	_, err := db.c.Exec("DELETE FROM comments WHERE id_photo = ? AND id_comment = ?", photoID, commentID)
	if err != nil {
		return err
	}
	return nil
}