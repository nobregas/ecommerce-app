import { create } from 'zustand';
import cartService, { CartItemType } from '@/service/cartService';

interface CartStore {
  cartItems: CartItemType[];
  total: number;
  cartCount: number;
  fetchCartItems: () => Promise<void>;
  calculateTotal: () => void;
  addToCart: (productId: number) => Promise<void>;
  removeOneFromCart: (productId: number) => Promise<void>;
  removeItemFromCart: (productId: number) => Promise<void>;
}

export const useCartStore = create<CartStore>((set, get) => ({
  cartItems: [],
  total: 0,
  cartCount: 0,

  fetchCartItems: async () => {
    try {
      const items = await cartService.getCartItems();
      set({ 
        cartItems: items,
        cartCount: items.reduce((total, item) => total + item.quantity, 0),
      });
      get().calculateTotal();
    } catch (error) {
      console.error('Error fetching cart items:', error);
    }
  },

  calculateTotal: () => {
    const items = get().cartItems;
    const newTotal = items.reduce(
      (sum, item) => sum + (item.priceAtAdding * item.quantity), 0
    );
    set({ total: newTotal });
  },

  addToCart: async (productId: number) => {
    try {
      await cartService.addToCart(productId);
      await get().fetchCartItems();
    } catch (error) {
      console.error('Error adding to cart:', error);
      throw error;
    }
  },

  removeOneFromCart: async (productId: number) => {
    try {
      await cartService.removeOneFromCart(productId);
      await get().fetchCartItems();
    } catch (error) {
      console.error('Error removing item:', error);
      throw error;
    }
  },

  removeItemFromCart: async (productId: number) => {
    try {
      await cartService.removeItemFromCart(productId);
      await get().fetchCartItems();
    } catch (error) {
      console.error('Error removing item:', error);
      throw error;
    }
  }
}));