import { create } from 'zustand';
import orderService, { CreateOrderPayload, OrderHistory, OrderStatus, OrderWithItems, PaymentMethod } from '@/service/orderService';

interface OrderStore {
  orders: OrderHistory[];
  orderDetails: OrderWithItems | null;
  isLoading: boolean;
  error: string | null;
  
  fetchOrders: () => Promise<void>;
  fetchOrderWithItems: (orderId: number) => Promise<void>;
  createOrder: (payload: CreateOrderPayload) => Promise<OrderHistory>;
  updateOrderStatus: (orderId: number, status: OrderStatus) => Promise<void>;
  clearError: () => void;
}

export const useOrderStore = create<OrderStore>((set, get) => ({
  orders: [],
  orderDetails: null,
  isLoading: false,
  error: null,

  fetchOrders: async () => {
    try {
      set({ isLoading: true, error: null });
      const orders = await orderService.getOrders();
      set({ orders, isLoading: false });
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Error fetching orders', isLoading: false });
      console.error('Error fetching orders:', error);
    }
  },

  fetchOrderWithItems: async (orderId: number) => {
    try {
      set({ isLoading: true, error: null });
      const orderDetails = await orderService.getOrderWithItemsById(orderId);
      set({ orderDetails, isLoading: false });
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Error fetching order details', isLoading: false });
      console.error('Error fetching order details:', error);
    }
  },

  createOrder: async (payload: CreateOrderPayload) => {
    try {
      set({ isLoading: true, error: null });
      const newOrder = await orderService.createOrder(payload);
      set(state => ({ 
        orders: [newOrder, ...state.orders],
        isLoading: false 
      }));
      return newOrder;
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Error creating order', isLoading: false });
      console.error('Error creating order:', error);
      throw error;
    }
  },

  updateOrderStatus: async (orderId: number, status: OrderStatus) => {
    try {
      set({ isLoading: true, error: null });
      await orderService.updateOrderStatus(orderId, status);
      
      set(state => ({
        orders: state.orders.map(order => 
          order.id === orderId ? { ...order, status, updatedAt: new Date().toISOString() } : order
        ),
        orderDetails: state.orderDetails?.order.id === orderId ? {
          ...state.orderDetails,
          order: {
            ...state.orderDetails.order,
            status,
            updatedAt: new Date().toISOString()
          }
        } : state.orderDetails,
        isLoading: false
      }));
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Error updating order status', isLoading: false });
      console.error('Error updating order status:', error);
      throw error;
    }
  },

  clearError: () => set({ error: null })
})); 