package notification

import (
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/apperrors"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/types"
	"github.com/nobregas/ecommerce-mobile-back/internal/shared/utils"
)

type Service struct {
	notificationStore types.NotificationStore
	userStore         types.UserStore
}

func NewNotificationService(notificationStore types.NotificationStore) *Service {
	return &Service{notificationStore: notificationStore}
}

func (s *Service) GetMyNotifications(userID int) *[]types.Notification {
	_, err := s.userStore.GetUserByID(userID)
	if err != nil {
		panic(apperrors.NewEntityNotFound("user", userID))
		return nil
	}

	notifications, err := s.notificationStore.GetMyNotifications(userID)
	if err != nil {
		panic(err)
		return nil
	}

	return notifications
}

func (s *Service) GetNotifications() *[]types.Notification {
	notifications, err := s.notificationStore.GetNotifications()
	if err != nil {
		panic(err)
		return nil
	}

	return notifications
}

func (s *Service) GetNotificationByID(notificationID int) *types.Notification {
	notification, err := s.notificationStore.GetNotificationByID(notificationID)
	if err != nil {
		panic(apperrors.NewEntityNotFound("notification", notificationID))
		return nil
	}

	return notification
}

func (s *Service) CreateNotification(payload *types.CreateNotificationPayload, userID int) *types.Notification {
	if err := utils.Validate.Struct(payload); err != nil {
		panic(apperrors.NewValidationError("invalid payload", err.Error()))
		return nil
	}

	createdNotification, err := s.notificationStore.CreateNotification(payload, userID)
	if err != nil {
		panic(err)
		return nil
	}

	return createdNotification
}

func (s *Service) DeleteNotification(notificationID int) {
	_, err := s.notificationStore.GetNotificationByID(notificationID)
	if err != nil {
		panic(apperrors.NewEntityNotFound("Notification", notificationID))
	}

	if err := s.notificationStore.DeleteNotification(notificationID); err != nil {
		panic(err)
	}
}
