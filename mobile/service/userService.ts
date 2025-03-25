import api from "./apiClient";

export interface UserDTO {
  id: number;
  fullName: string;
  email: string;
  cpf: string;
  profileImg: string | null;
  createdAt: string;
}

class UserService {
  async getCurrentUser(): Promise<UserDTO> {
    try {
      const response = await api.get("/me");
      const userData = response.data;
      
      return {
        ...userData,
        createdAt: this.formatDateTime(userData.createdAt)
      };
    } catch (error) {
      throw this.handleError(error);
    }
  }

  async updateUserProfile(userData: Partial<Omit<UserDTO, 'id' | 'createdAt'>>): Promise<UserDTO> {
    try {
      const response = await api.put("/me", userData);
      return {
        ...response.data,
        createdAt: this.formatDateTime(response.data.createdAt)
      };
    } catch (error) {
      throw this.handleError(error);
    }
  }

  private formatDateTime(dateTimeStr: string): string {
    try {
      const date = new Date(dateTimeStr);
      
      if (isNaN(date.getTime())) {
        return dateTimeStr;
      }
      
      const day = date.getDate().toString().padStart(2, '0');
      const month = (date.getMonth() + 1).toString().padStart(2, '0');
      const year = date.getFullYear();
      
      return `${day}/${month}/${year}`;
    } catch (error) {
      return dateTimeStr;
    }
  }

  private handleError(error: any): Error {
    const defaultMessage = "Failed to get user data";
    if (error.message) {
      return new Error(error.message || defaultMessage);
    }
    return new Error(defaultMessage);
  }
}

export default new UserService(); 