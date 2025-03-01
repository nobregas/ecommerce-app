import { CartItemType, CategoryType, NotificationType, ProductStrType, ProductType } from '@/types/type';
import axios from 'axios';
import { get } from 'react-native/Libraries/TurboModule/TurboModuleRegistry';

const API_BASE_URL = 'http://10.0.2.2:8000';

const api = axios.create({
    baseURL: API_BASE_URL,
});


export const getProductDetails = async (id: string, productType: ProductStrType) => {
    const endPoint = productType === "regular"
        ? `/products/${id}`
        : `/saleProducts/${id}`

    const response = await api.get(endPoint);
    return response.data;
}

export const getProducts = async (): Promise<ProductType[]> => {
    const response = await api.get('/products');
    return response.data;
}

export const getCategories = async (): Promise<CategoryType[]> => {
    const response = await api.get('/categories');
    return response.data;
}

export const getSaleProducts = async (): Promise<ProductType[]> => {
    const response = await api.get('/saleProducts');
    return response.data;
}

export const getNotifications = async (): Promise<NotificationType[]> => {
    const response = await api.get('/notifications');
    return response.data;
}

export const getCartItems = async (): Promise<CartItemType[]> => {
    const response = await api.get('/cart');
    return response.data;
}