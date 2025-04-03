import { useState } from 'react';
import authService from '@/service/authService';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { router } from 'expo-router';
import { useCartStore } from '@/store/cardBadgeStore';

export const useAuth = () => {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const { fetchCartItems } = useCartStore();
  
    const handleAuth = async (
      action: 'login' | 'register',
      payload: any,
      successRedirect: string = '/(tabs)'
    ) => {
      setLoading(true);
      setError(null);
      
      try {
        let response;
        if (action === 'login') {
          response = await authService.login(payload);
          await AsyncStorage.setItem('authToken', response.token);
          
          await fetchCartItems();
          
          router.dismissAll();
          router.push('/(tabs)');
        
        } else {
          response = await authService.register(payload);
          router.dismissAll();
          router.push('/(auth)/signin');
        }
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };
  
    return { handleAuth, loading, error };
  };