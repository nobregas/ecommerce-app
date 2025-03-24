package notification

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetMyNotifications(userID int) (*[]types.Notification, error) {
	query := `
		SELECT id, userId, title, message, isRead, createdAt
		FROM notifications
		WHERE userId = ?
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("[GetMyNotifications] error getting notifications of user %d: %v", userID, err)
	}
	defer rows.Close()

	notifications := make([]types.Notification, 0)
	for rows.Next() {
		n, err := scanRowsIntoNotification(rows)
		if err != nil {
			return nil, fmt.Errorf("[GetMyNotifications] error scanning rows: %v", err)
		}
		notifications = append(notifications, *n)
	}

	return &notifications, nil
}

func (s *Store) GetNotifications() (*[]types.Notification, error) {
	rows, err := s.db.Query(`
		SELECT id, userId, title, message, isRead, createdAt
		FROM notifications
	`)

	if err != nil {
		return nil, fmt.Errorf("[GetNotifications] error getting notifications: %v", err)
	}
	defer rows.Close()

	notifications := make([]types.Notification, 0)
	for rows.Next() {
		n, err := scanRowsIntoNotification(rows)
		if err != nil {
			return nil, fmt.Errorf("[GetNotifications] error scanning rows: %v", err)
		}

		notifications = append(notifications, *n)
	}

	return &notifications, nil
}

func (s *Store) GetNotificationByID(notificationID int) (*types.Notification, error) {
	query := `
		SELECT id, userId, title, message, isRead, createdAt
		FROM notifications
		WHERE id = ?
	`

	row := s.db.QueryRow(query, notificationID)

	n, err := scanRowIntoNotification(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("[GetNotificationByID] error getting notification of id %d", notificationID)
	}

	return n, nil
}

func (s *Store) CreateNotification(payload *types.CreateNotificationPayload, userID int) (*types.Notification, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("[CreateNotification] error begin transaction: %v", err)
	}
	defer tx.Rollback()

	var userExists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`, userID).Scan(&userExists)
	if err != nil || !userExists {
		return nil, fmt.Errorf("[CreateNotification] error creating notification, user does not exist %d: %v", userID, err)
	}

	res, err := tx.Exec(`
		INSERT INTO notifications(userId, title, message)
		VALUES (?, ?, ?)
	`, userID, payload.Title, payload.Message)
	if err != nil {
		return nil, fmt.Errorf("[CreateNotification] error creating notification: %v", err)
	}

	notificationID, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("[CreateNotification] error get notification ID: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("[CreateNotification] error commit transaction: %v", err)
	}

	return s.GetNotificationByID(int(notificationID))
}

func (s *Store) DeleteNotification(notificationID int) error {
	res, err := s.db.Exec(`
		DELETE FROM notifications 
		WHERE id = ?
	`, notificationID)
	if err != nil {
		return fmt.Errorf("[DeleteNotification] failed to delete notification: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("[DeleteNotification] failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("[DeleteNotification] no rows affected")
	}

	return nil
}

func scanRowsIntoNotification(rows *sql.Rows) (*types.Notification, error) {
	notification := new(types.Notification)

	err := rows.Scan(
		&notification.ID,
		&notification.UserID,
		&notification.Title,
		&notification.Message,
		&notification.IsRead,
		&notification.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("[scanRowsIntoNotification] error scanning rows: %v", err)
	}

	return notification, nil
}

func scanRowIntoNotification(rows *sql.Row) (*types.Notification, error) {
	notification := new(types.Notification)

	err := rows.Scan(
		&notification.ID,
		&notification.UserID,
		&notification.Title,
		&notification.Message,
		&notification.IsRead,
		&notification.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("[scanRowsIntoNotification] error scanning row: %v", err)
	}

	return notification, nil
}
