import { create } from 'zustand';

type FavoritesStore = {
  favorites: Set<number>;
  toggleFavorite: (productId: number, state?: boolean) => void;
  isFavorite: (productId: number) => boolean;
  syncFromServer: (favorites: number[]) => void;
};

export const useFavoritesStore = create<FavoritesStore>((set, get) => ({
  favorites: new Set<number>(),
  
  toggleFavorite: (productId, state) => {
    set(({ favorites }) => {
      const newFavorites = new Set(favorites);
      const shouldAdd = state ?? !newFavorites.has(productId);
      
      if (shouldAdd) {
        newFavorites.add(productId);
      } else {
        newFavorites.delete(productId);
      }
      
      return { favorites: newFavorites };
    });
  },

  isFavorite: (productId) => get().favorites.has(productId),
  
  syncFromServer: (favorites) => {
    set({ favorites: new Set(favorites) });
  },
}));