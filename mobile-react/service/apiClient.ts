import axios from "axios";
import AsyncStorage from "@react-native-async-storage/async-storage";

const IPV4 = "10.5.7.125"

const api = axios.create({
  baseURL: `http://${IPV4}:8080/api/v1/`,
  headers: {
    "Content-Type": "application/json",
  },
});

api.interceptors.request.use(async (config) => {
  const token = await AsyncStorage.getItem('authToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response) {
      return Promise.reject(error.response.data);
    } else if (error.request) {
      return Promise.reject({ message: "No response from server" });
    } else {
      return Promise.reject({ message: "Request configuration error" });
    }
  }
);

export default api;