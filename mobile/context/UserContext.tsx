import React, { createContext, useState, useContext, useEffect, ReactNode } from 'react';
import { UserDTO } from '../service/userService';
import userService from '../service/userService';
import AsyncStorage from '@react-native-async-storage/async-storage';

interface UserContextType {
  user: UserDTO | null;
  isLoading: boolean;
  error: string | null;
  fetchUser: () => Promise<void>;
  updateUser: (userData: Partial<Omit<UserDTO, 'id' | 'createdAt'>>) => Promise<void>;
  logout: () => Promise<void>;
}

// Inicializando o contexto com um valor padrão para evitar undefined
const UserContext = createContext<UserContextType>({
  user: null,
  isLoading: false,
  error: null,
  fetchUser: async () => {},
  updateUser: async () => {},
  logout: async () => {}
});

const USER_STORAGE_KEY = '@user_data';
const USER_TIMESTAMP_KEY = '@user_timestamp';
const CACHE_DURATION = 10 * 60 * 1000; // 10 min 

export const UserProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<UserDTO | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [initialLoadDone, setInitialLoadDone] = useState<boolean>(false);

  useEffect(() => {
    const loadInitialData = async () => {
      setIsLoading(true);
      try {
        const userDataString = await AsyncStorage.getItem(USER_STORAGE_KEY);
        const timestampString = await AsyncStorage.getItem(USER_TIMESTAMP_KEY);
        
        if (userDataString && timestampString) {
          const userData = JSON.parse(userDataString) as UserDTO;
          const timestamp = parseInt(timestampString, 10);
          const now = Date.now();
          
          if (now - timestamp < CACHE_DURATION) {
            setUser(userData);
          } else {
            // Cache expirado, buscar dados atualizados
            await fetchUserData();
          }
        } else {
          // Sem dados armazenados, buscar do servidor
          await fetchUserData();
        }
      } catch (error) {
        console.error('Erro ao carregar dados iniciais do usuário:', error);
        setError('Falha ao carregar dados do usuário');
      } finally {
        setIsLoading(false);
        setInitialLoadDone(true);
      }
    };

    if (!initialLoadDone) {
      loadInitialData();
    }
  }, [initialLoadDone]);

  const fetchUserData = async () => {
    try {
      const userData = await userService.getCurrentUser();
      setUser(userData);
      
      await AsyncStorage.setItem(USER_STORAGE_KEY, JSON.stringify(userData));
      await AsyncStorage.setItem(USER_TIMESTAMP_KEY, Date.now().toString());
    } catch (error) {
      console.error('Erro ao buscar dados do usuário:', error);
      throw error;
    }
  };

  const fetchUser = async () => {
    if (isLoading) return;
    
    setIsLoading(true);
    setError(null);
    
    try {
      await fetchUserData();
    } catch (error) {
      setError('Falha ao buscar dados do usuário');
    } finally {
      setIsLoading(false);
    }
  };

  const updateUser = async (userData: Partial<Omit<UserDTO, 'id' | 'createdAt'>>) => {
    if (isLoading) return;
    
    setIsLoading(true);
    setError(null);
    
    try {
      const updatedUser = await userService.updateUserProfile(userData);
      setUser(updatedUser);
      
      await AsyncStorage.setItem(USER_STORAGE_KEY, JSON.stringify(updatedUser));
      await AsyncStorage.setItem(USER_TIMESTAMP_KEY, Date.now().toString());
    } catch (error) {
      console.error('Erro ao atualizar dados do usuário:', error);
      setError('Falha ao atualizar dados do usuário');
    } finally {
      setIsLoading(false);
    }
  };

  const logout = async () => {
    try {
      await AsyncStorage.removeItem(USER_STORAGE_KEY);
      await AsyncStorage.removeItem(USER_TIMESTAMP_KEY);
      setUser(null);
    } catch (error) {
      console.error('Erro ao fazer logout:', error);
    }
  };

  return (
    <UserContext.Provider value={{ user, isLoading, error, fetchUser, updateUser, logout }}>
      {children}
    </UserContext.Provider>
  );
};

export const useUser = (): UserContextType => {
  const context = useContext(UserContext);
  return context;
}; 