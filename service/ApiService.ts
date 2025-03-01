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

export const addToCart = async (product: ProductType) => {
    try {
        const cartItems = await getCartItems()
        const existingItem = cartItems.find(item => item.id === product.id)

        if (existingItem) {
            await api.patch(`/cart/${existingItem.id}`, { quantity: existingItem.quantity + 1 })  
        } else {
            // verify if image is an array or a string and then get the first image or use a placeholder
            const imageUrl = typeof product.images === "string"
            ? product.images
            : Array.isArray(product.images) && product.images.length > 0
            ? product.images[0]
            : "https://via.placeholder.com/150";

            await api.post('/cart', {
                id: product.id,
                title: product.title,
                price: product.price,
                quantity: 1,
                image: product.images || product.images[0]
            })
        }

        return "Item added to cart successfully";
    } catch (error) {
        console.error(error)
        throw error
    }
}