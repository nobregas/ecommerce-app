import { formatDateTime } from "@/utils/sharedFunctions";
import api from "./apiClient";

export interface CartItemType {
    cartId: number;
    productId: number;
    productTitle: string;
    productImage: string;
    quantity: number;
    priceAtAdding: number;
    addedAt: string;
}

class CartService {
    async getCartItems(): Promise<CartItemType[]> {
        try {
            const response = await api.get(`/cart/items`);
            return response.data.map((item: CartItemType) => ({
              ...item,
              addedAt: formatDateTime(item.addedAt)
            }));
          } catch (error) {
            throw this.handleError(error);
          }
    }

    async getCartCount(): Promise<number> {
        try {
            const items = await this.getCartItems();
            return items.reduce((total, item) => total + item.quantity, 0);
        } catch (error) {
            throw this.handleError(error);
        }
    }

    async addToCart(productId: number) {
        try {
            await api.post(`/cart/items/${productId}`);

        } catch (error) {
            throw this.handleError(error);
        }
    }

    async removeOneFromCart(productId: number) {
        try {
            await api.delete(`/cart/items/${productId}`);
        }catch (error) {
            throw this.handleError(error); 
        }
    }

    async removeItemFromCart(productId: number) {
        try {
            await api.delete(`/cart/items/${productId}/remove`);
        } catch (error) {
            throw this.handleError(error);
        }
    }

    async getTotal(): Promise<number> {
        try {
            const response = await api.get(`/cart/total`);
            return response.data;
        } catch (error) {
            throw this.handleError(error);
        }
    }

    async clearCart() {
        try {
            await api.delete(`/cart/clear`);
        } catch (error) {
            throw this.handleError(error);
        }
    }

    private handleError(error: any): Error {
        console.log("Error", error);
        const defaultMessage = "Failed to fetch cart";
        if (error.message) {
          return new Error(error.message || defaultMessage);
        }
        return new Error(defaultMessage);
    }
}

export default new CartService();