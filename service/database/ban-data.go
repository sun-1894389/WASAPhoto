package database

import (
	"errors"
	"fmt"
)

// BanUser bans a user (banned) by another user (banner)
func (db *appdbimpl) BanUser(bannerUser User, bannedUser User) error {
	query := "INSERT INTO banned_users (banner, banned) VALUES (?, ?)"
	_, err := db.c.Exec(query, bannerUser.IdUser, bannedUser.IdUser)
	if err != nil {
		return fmt.Errorf("failed to ban user: %w", err)
	}
	return nil
}

// UnbanUser removes a user (banned) from the banned list of another user (banner)
func (db *appdbimpl) UnbanUser(bannerUser User, bannedUser User) error {
	query := "DELETE FROM banned_users WHERE banner = ? AND banned = ?"
	_, err := db.c.Exec(query, bannerUser.IdUser, bannedUser.IdUser)
	if err != nil {
		return fmt.Errorf("failed to unban user: %w", err)
	}
	return nil
}

// BannedUserCheck checks if the requesting user was banned by another user.
// It returns true if the user is banned, false otherwise.
func (db *appdbimpl) BannedUserCheck(requestingUser User, targetUser User) (bool, error) {
	query := "SELECT COUNT(*) FROM banned_users WHERE banned = ? AND banner = ?"
	var cnt int
	err := db.c.QueryRow(query, requestingUser.IdUser, targetUser.IdUser).Scan(&cnt)
	if err != nil {
		return false, fmt.Errorf("failed to check banned user: %w", err)
	}
	return cnt == 1, nil
}