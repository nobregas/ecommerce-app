import { create } from 'zustand';
import { createJSONStorage, persist } from 'zustand/middleware';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { aiChatService, ChatMessageType } from '../service/aiChatService';

// Chat state interface
export interface ChatState {
  messages: ChatMessageType[];
  isLoading: boolean;
  error: string | null;
  
  // Actions
  sendMessage: (content: string) => Promise<void>;
  clearMessages: () => void;
  clearError: () => void;
}

// Welcome message that appears when no history exists
const welcomeMessage: ChatMessageType = {
  id: 'welcome',
  content: 'Hello! I am the virtual assistant of ShopX. How can I help you today?',
  role: 'assistant',
  timestamp: new Date(),
};

// Create a custom storage that handles date serialization/deserialization
const customStorage = {
  getItem: async (name: string): Promise<string | null> => {
    try {
      const value = await AsyncStorage.getItem(name);
      if (value === null) return null;
      
      const parsed = JSON.parse(value);
      
      // Add a safety check for the structure
      if (parsed && parsed.state && parsed.state.messages) {
        parsed.state.messages = parsed.state.messages.map((msg: any) => {
          // Ensure timestamp is a valid Date object
          try {
            return {
              ...msg,
              timestamp: new Date(msg.timestamp)
            };
          } catch (err) {
            // If there's an error parsing the date, use current date
            console.warn('Error parsing date, using current date instead:', err);
            return {
              ...msg,
              timestamp: new Date()
            };
          }
        });
      }
      
      return JSON.stringify(parsed);
    } catch (error) {
      console.error('Error getting items from storage:', error);
      return null;
    }
  },
  
  setItem: async (name: string, value: string): Promise<void> => {
    try {
      const parsed = JSON.parse(value);
      
      // Add a safety check for the structure
      if (parsed && parsed.state && parsed.state.messages) {
        parsed.state.messages = parsed.state.messages.map((msg: any) => {
          // Ensure timestamp is stored as ISO string
          try {
            return {
              ...msg,
              timestamp: msg.timestamp instanceof Date ? msg.timestamp.toISOString() : 
                         typeof msg.timestamp === 'string' ? msg.timestamp : 
                         new Date().toISOString()
            };
          } catch (err) {
            console.warn('Error serializing date, using current date instead:', err);
            return {
              ...msg,
              timestamp: new Date().toISOString()
            };
          }
        });
      }
      
      return AsyncStorage.setItem(name, JSON.stringify(parsed));
    } catch (error) {
      console.error('Error setting items to storage:', error);
    }
  },
  
  removeItem: async (name: string): Promise<void> => {
    return AsyncStorage.removeItem(name);
  }
};

// Zustand store for managing chat state with persistence
export const useChatStore = create<ChatState>()(
  persist(
    (set, get) => ({
      messages: [welcomeMessage],
      isLoading: false,
      error: null,

      // Method to send a message
      sendMessage: async (content: string) => {
        try {
          // Add user message to state
          const userMessage: ChatMessageType = {
            id: Date.now().toString(),
            content,
            role: 'user',
            timestamp: new Date(),
          };

          set(state => ({
            messages: [...state.messages, userMessage],
            isLoading: true,
            error: null
          }));

          // Send message to AI service
          const assistantMessage = await aiChatService.sendMessage(content);

          // Add assistant response to state
          set(state => ({
            messages: [...state.messages, assistantMessage],
            isLoading: false
          }));
        } catch (error) {
          // Handle errors
          console.error('Error sending message:', error);
          set({
            isLoading: false,
            error: error instanceof Error ? error.message : 'An error occurred while processing your message'
          });

          // Clear error after 5 seconds
          setTimeout(() => {
            set({ error: null });
          }, 5000);
        }
      },

      // Method to clear all messages
      clearMessages: () => {
        set({
          messages: [welcomeMessage]
        });
      },

      // Method to clear error messages
      clearError: () => {
        set({ error: null });
      }
    }),
    {
      name: 'shopx-chat-storage',
      storage: createJSONStorage(() => customStorage),
      partialize: (state) => ({ messages: state.messages }),
    }
  )
); 