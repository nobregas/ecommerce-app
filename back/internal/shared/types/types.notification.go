package types

import "time"

type NotificationService interface {
	GetMyNotifications(userID int) *[]Notification
	GetNotifications() *[]Notification
	GetNotificationByID(notificationID int) *Notification
	CreateNotification(payload *CreateNotificationPayload, userID int) *Notification
	DeleteNotification(notificationID int)
}

type NotificationStore interface {
	GetMyNotifications(userID int) (*[]Notification, error)
	GetNotifications() (*[]Notification, error)
	GetNotificationByID(notificationID int) (*Notification, error)
	CreateNotification(payload *CreateNotificationPayload, userID int) (*Notification, error)
	DeleteNotification(notificationID int) error
}

type Notification struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"isRead"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateNotificationPayload struct {
	Title   string `json:"title" validate:"required, min=3, max=100"`
	Message string `json:"message" validate:"required"`
}

type UpdateNotificationPayload struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}
