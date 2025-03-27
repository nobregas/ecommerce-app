import api from "./apiClient";

import { formatDateTime } from "@/utils/sharedFunctions";

export interface Notification {
  id: number;
  userId: number;
  title: string;
  message: string;
  isRead: boolean;
  createdAt: string;
}

class NotificationService {
  async getNotifications(): Promise<Notification[]> {
    try {
      const response = await api.get(`/notification/my`);
      return response.data.map((notification: Notification) => ({
        ...notification,
        createdAt: formatDateTime(notification.createdAt)
      }));
    } catch (error) {
      throw this.handleError(error);
    }
  }

//   async getNotificationById(id: number): Promise<Notification> {
//     try {
//       const response = await api.get(`/notifications/${id}`);
//       return response.data;
//     } catch (error) {
//       throw this.handleError(error);
//     }
//   }

//   async markAsRead(id: number): Promise<void> {
//     try {
//       await api.patch(`/notifications/${id}`, { isRead: true });
//     } catch (error) {
//       throw this.handleError(error);
//     }
//   }

//   async markAllAsRead(userId: number): Promise<void> {
//     try {
//       await api.patch(`/notifications/read-all`, { userId });
//     } catch (error) {
//       throw this.handleError(error);
//     }
//   }

  private handleError(error: any): Error {
    const defaultMessage = "Failed to fetch notifications";
    if (error.message) {
      return new Error(error.message || defaultMessage);
    }
    return new Error(defaultMessage);
  }
}

export default new NotificationService();
