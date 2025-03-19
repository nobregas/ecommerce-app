import { create } from 'zustand';
import { getCartCount } from '@/service/ApiService';

interface CartStore {
  cartCount: number;
  fetchCartCount: () => Promise<void>;
}

export const useCartStore = create<CartStore>((set) => ({
  cartCount: 0,
  
  fetchCartCount: async () => {
    try {
      const count = await getCartCount();
      set({ cartCount: count });
    } catch (error) {
      console.error('Error fetching cart count:', error);
    }
  }
}));

useCartStore.getState().fetchCartCount();