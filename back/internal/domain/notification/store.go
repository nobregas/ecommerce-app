package notification

import (
	"database/sql"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetMyNotifications(userID int) (*[]types.Notification, error) {
	return nil, nil
}

func (s *Store) GetNotifications() (*[]types.Notification, error) {
	return nil, nil
}

func (s *Store) GetNotificationByID(notificationID int) (*types.Notification, error) {
	return nil, nil
}

func (s *Store) CreateNotification(payload *types.CreateNotificationPayload) (*types.Notification, error) {
	return nil, nil
}

func (s *Store) UpdateNotification(payload *types.UpdateNotificationPayload, notificationID int) (*types.Notification, error) {
	return nil, nil
}

func (s *Store) DeleteNotification(notificationID int) error {
	return nil
}
