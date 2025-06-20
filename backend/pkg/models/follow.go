package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type FollowModel struct {
	DB *sql.DB
}

func (m *FollowModel) Follow(ctx context.Context, followerID, followedID string) error {
	if followerID == followedID {
		return errors.New("cannot follow yourself")
	}

	now := time.Now().Unix() // 8:35 PM EAT, May 15, 2025 = 1744732500
	id := uuid.New().String()

	// Fetch is_private status of the followed user
	isPrivate, err := m.IsUserProfilePrivate(ctx, followedID)
	if err != nil {
		// IsUserProfilePrivate already returns a specific error for "user not found"
		return err
	}

	status := "accepted"
	if isPrivate {
		status = "pending"
	}

	stmt := `
        INSERT INTO follows (id, follower_id, followed_id, status, created_at)
        VALUES (?, ?, ?, ?, ?)
        ON CONFLICT (follower_id, followed_id) DO NOTHING;
		`
	result, err := m.DB.ExecContext(ctx, stmt, id, followerID, followedID, status, now)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("follow request already exists or was already accepted")
	}

	return nil
}

func (m *FollowModel) Unfollow(ctx context.Context, followerID, followedID string) error {
	stmt := `
        DELETE FROM follows
        WHERE follower_id = ? AND followed_id = ?;
    `
	result, err := m.DB.ExecContext(ctx, stmt, followerID, followedID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("follow does not exist")
	}

	return nil
}

func (m *FollowModel) GetFollowers(ctx context.Context, userID string) ([]User, error) {
	stmt := `
        SELECT u.id, u.email, u.first_name, u.last_name, u.nickname, u.date_of_birth,
               u.about_me, u.avatar_url, u.is_private, u.created_at, u.updated_at
        FROM users u
        JOIN follows f ON u.id = f.follower_id
        WHERE f.followed_id = ? AND f.status = 'accepted';
    `
	rows, err := m.DB.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var followers []User
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Nickname,
			&user.DateOfBirth,
			&user.AboutMe,
			&user.AvatarURL,
			&user.IsPrivate,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		followers = append(followers, user)
	}
	if followers == nil {
		followers = []User{}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return followers, nil
}

func (m *FollowModel) GetFollowing(ctx context.Context, userID string) ([]User, error) {
	stmt := `
        SELECT u.id, u.email, u.first_name, u.last_name, u.nickname, u.date_of_birth,
               u.about_me, u.avatar_url, u.is_private, u.created_at, u.updated_at
        FROM users u
        JOIN follows f ON u.id = f.followed_id
        WHERE f.follower_id = ? AND f.status = 'accepted';
    `
	rows, err := m.DB.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var following []User
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Nickname,
			&user.DateOfBirth,
			&user.AboutMe,
			&user.AvatarURL,
			&user.IsPrivate,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		following = append(following, user)
	}
	// Initialize empty slice if no rows returned
	if following == nil {
		following = []User{}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return following, nil
}

// IsUserProfilePrivate checks if a user's profile is private.
func (m *FollowModel) IsUserProfilePrivate(ctx context.Context, userID string) (bool, error) {
	var isPrivate bool
	query := `SELECT is_private FROM users WHERE id = ?`
	err := m.DB.QueryRowContext(ctx, query, userID).Scan(&isPrivate)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("user not found")
		}
		return false, err
	}
	return isPrivate, nil
}

// AcceptFollowRequest updates a pending follow request to 'accepted'.
// currentUserID is the ID of the user accepting the request (the one who was followed).
// requesterID is the ID of the user who sent the request (the follower).
func (m *FollowModel) AcceptFollowRequest(ctx context.Context, currentUserID, requesterID string) error {
	stmt := `
		UPDATE follows
		SET status = 'accepted'
		WHERE follower_id = ? AND followed_id = ? AND status = 'pending';
	`
	result, err := m.DB.ExecContext(ctx, stmt, requesterID, currentUserID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("follow request not found or already accepted/declined")
	}
	return nil
}

// DeclineFollowRequest removes a pending follow request.
// currentUserID is the ID of the user declining the request (the one who was followed).
// requesterID is the ID of the user who sent the request (the follower).
func (m *FollowModel) DeclineFollowRequest(ctx context.Context, currentUserID, requesterID string) error {
	stmt := `
		DELETE FROM follows
		WHERE follower_id = ? AND followed_id = ? AND status = 'pending';
	`
	result, err := m.DB.ExecContext(ctx, stmt, requesterID, currentUserID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("follow request not found or already actioned")
	}
	return nil
}

// GetPendingFollowRequests retrieves users who have sent a follow request to userID that is still 'pending'.
func (m *FollowModel) GetPendingFollowRequests(ctx context.Context, userID string) ([]User, error) {
	stmt := `
        SELECT u.id, u.email, u.first_name, u.last_name, u.nickname, u.date_of_birth,
               u.about_me, u.avatar_url, u.is_private, u.created_at, u.updated_at
        FROM users u
        JOIN follows f ON u.id = f.follower_id
        WHERE f.followed_id = ? AND f.status = 'pending';
    `
	rows, err := m.DB.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pendingRequests []User
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Nickname,
			&user.DateOfBirth,
			&user.AboutMe,
			&user.AvatarURL,
			&user.IsPrivate,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		pendingRequests = append(pendingRequests, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Initialize empty slice if no rows returned, as per requirement
	if pendingRequests == nil {
		pendingRequests = []User{}
	}

	return pendingRequests, nil
}
