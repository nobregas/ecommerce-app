import AsyncStorage from "@react-native-async-storage/async-storage";
import axios from "axios";

const api = axios.create({
    baseURL: "http://192.168.1.100:8080/api/v1/",
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

export interface LoginPayload {
    email: string;
    password: string;
}

export interface RegisterPayload {
    fullName: string;
    email: string;
    cpf: string;
    password: string;
}

class AuthService {
    async login(payload: LoginPayload) {
        try {
            const response = await api.post('/login', {
                Email: payload.email,
                Password: payload.password,
            });
            return response.data;
        } catch (error) {
            throw this.handleError(error);
        }
    }

    async register(payload: RegisterPayload) {
        try {
            console.log(payload)
            const response = await api.post('/register', {
                FullName: payload.fullName,
                Email: payload.email,
                Cpf: payload.cpf,
                Password: payload.password,
            });
            console.log("1")
            return response.data;
        } catch (error) {
            //throw this.handleError(error);
        }
    }

    private handleError(error: any) {
        if (error.response) {
            throw new Error(error.response.data.message || 'Request error');
        } else if (error.request) {
            throw new Error('No response from server');
        } else {
            throw new Error('Request configuration error');
        }
    }
};

export default new AuthService();