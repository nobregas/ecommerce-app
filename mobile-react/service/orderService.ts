import api from "./apiClient";
import { formatDateTime } from "@/utils/sharedFunctions";

export enum OrderStatus {
    PENDING = 'PENDING',
    PROCESSING = 'PROCESSING',
    SHIPPED = 'SHIPPED',
    DELIVERED = 'DELIVERED',
    CANCELLED = 'CANCELLED'
}

export enum PaymentMethod {
    CREDIT_CARD = 'CREDIT_CARD',
    DEBIT_CARD = 'DEBIT_CARD',
    PIX = 'PIX'
}

export interface OrderItem {
    orderId: number;
    productId: number;
    quantity: number;
    price: number;
}

export interface OrderHistory {
    id: number;
    userId: number;
    totalAmount: number;
    status: OrderStatus;
    paymentMethod: PaymentMethod;
    paymentId: string;
    createdAt: string;
    updatedAt: string;
}

export interface OrderWithItems {
    order: OrderHistory;
    items: OrderItem[];
}

export interface CreateOrderPayload {
    paymentMethod: PaymentMethod;
    paymentId?: string;
}

class OrderService {
    async createOrder(payload: CreateOrderPayload): Promise<OrderHistory> {
        try {
            const response = await api.post('/orders', payload);
            return {
                ...response.data,
                createdAt: formatDateTime(response.data.createdAt),
                updatedAt: formatDateTime(response.data.updatedAt)
            };
        } catch (error) {
            throw this.handleError(error);
        }
    }

    async getOrders(): Promise<OrderHistory[]> {
        try {
            const response = await api.get('/orders');
            return response.data.map((order: OrderHistory) => ({
                ...order,
                createdAt: formatDateTime(order.createdAt),
                updatedAt: formatDateTime(order.updatedAt)
            }));
        } catch (error) {
            throw this.handleError(error);
        }
    }

    async getOrdersWithItems(): Promise<OrderWithItems[]> {
        try {
            const response = await api.get('/orders?withItems=true');
            return response.data.map((orderWithItems: OrderWithItems) => ({
                ...orderWithItems,
                order: {
                    ...orderWithItems.order,
                    createdAt: formatDateTime(orderWithItems.order.createdAt),
                    updatedAt: formatDateTime(orderWithItems.order.updatedAt)
                }
            }));
        } catch (error) {
            throw this.handleError(error);
        }
    }

    async getOrderById(orderId: number): Promise<OrderHistory> {
        try {
            const response = await api.get(`/orders/${orderId}`);
            return {
                ...response.data,
                createdAt: formatDateTime(response.data.createdAt),
                updatedAt: formatDateTime(response.data.updatedAt)
            };
        } catch (error) {
            throw this.handleError(error);
        }
    }

    async getOrderWithItemsById(orderId: number): Promise<OrderWithItems> {
        try {
            const response = await api.get(`/orders/${orderId}?withItems=true`);
            return {
                ...response.data,
                order: {
                    ...response.data.order,
                    createdAt: formatDateTime(response.data.order.createdAt),
                    updatedAt: formatDateTime(response.data.order.updatedAt)
                }
            };
        } catch (error) {
            throw this.handleError(error);
        }
    }

    async updateOrderStatus(orderId: number, status: OrderStatus): Promise<void> {
        try {
            await api.patch(`/orders/${orderId}/status`, { status });
        } catch (error) {
            throw this.handleError(error);
        }
    }

    private handleError(error: any): Error {
        const defaultMessage = "Failed to perform order operation";
        if (error.response?.data?.error) {
            return new Error(error.response.data.error);
        }
        if (error.message) {
            return new Error(error.message);
        }
        return new Error(defaultMessage);
    }
}

export default new OrderService(); 