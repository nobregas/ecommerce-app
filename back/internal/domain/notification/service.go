package notification

import "github.com/nobregas/ecommerce-mobile-back/internal/shared/types"

type Service struct {
	notificationStore types.NotificationStore
}

func NewNotificationService(notificationStore types.NotificationStore) *Service {
	return &Service{notificationStore: notificationStore}
}

func (s *Service) GetMyNotifications(userID int) *[]types.Notification {
	return nil
}

func (s *Service) GetNotifications() *[]types.Notification {
	return nil
}

func (s *Service) GetNotificationByID(notificationID int) *types.Notification {
	return nil
}

func (s *Service) CreateNotification(payload *types.CreateNotificationPayload) *types.Notification {
	return nil
}

func (s *Service) DeleteNotification(notificationID int) {

}
