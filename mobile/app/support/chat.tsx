import React, { useRef, useState, useEffect } from 'react';
import {
  StyleSheet,
  Text,
  View,
  FlatList,
  TouchableOpacity,
  TextInput,
  KeyboardAvoidingView,
  Platform,
  ActivityIndicator,
  Keyboard,
} from 'react-native';
import { Stack, useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { Colors } from '@/constants/Colors';
import { useChatStore } from '@/store/chatStore';
import { ChatMessageType } from '@/service/aiChatService';
import Animated, { FadeInDown, FadeInUp } from 'react-native-reanimated';
import { SafeAreaView } from 'react-native-safe-area-context';

const ChatSupportScreen = () => {
  const [message, setMessage] = useState('');
  const { messages, sendMessage, isLoading, error, clearError } = useChatStore();
  const flatListRef = useRef<FlatList>(null);
  const router = useRouter();
  const [keyboardVisible, setKeyboardVisible] = useState(false);

  // Monitor keyboard visibility
  useEffect(() => {
    const keyboardDidShowListener = Keyboard.addListener(
      'keyboardDidShow',
      () => {
        setKeyboardVisible(true);
        flatListRef.current?.scrollToEnd({ animated: true });
      }
    );
    const keyboardDidHideListener = Keyboard.addListener(
      'keyboardDidHide',
      () => {
        setKeyboardVisible(false);
      }
    );

    return () => {
      keyboardDidShowListener.remove();
      keyboardDidHideListener.remove();
    };
  }, []);

  const handleSendMessage = async () => {
    if (message.trim() === '' || isLoading) return;
    
    const trimmedMessage = message.trim();
    setMessage('');
    await sendMessage(trimmedMessage);
    
    // Scroll to bottom after sending message
    setTimeout(() => {
      flatListRef.current?.scrollToEnd({ animated: true });
    }, 100);
  };

  const renderMessage = ({ item, index }: { item: ChatMessageType; index: number }) => {
    const isUser = item.role === 'user';
    
    return (
      <Animated.View 
        entering={FadeInUp.delay(index * 100).duration(300)}
        style={[
          styles.messageContainer,
          isUser ? styles.userMessage : styles.assistantMessage
        ]}
      >
        <Text style={[styles.messageText, isUser ? styles.userMessageText : styles.assistantMessageText]}>
          {item.content}
        </Text>
        <Text style={styles.timestamp}>
          {new Date(item.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
        </Text>
      </Animated.View>
    );
  };

  return (
    <SafeAreaView style={styles.container} edges={['bottom']}>
      <Stack.Screen
        options={{
          headerTitle: 'Suporte ao Cliente',
          headerTitleAlign: 'center',
          headerLeft: () => (
            <TouchableOpacity onPress={() => router.back()} style={styles.backButton}>
              <Ionicons name="arrow-back" size={24} color={Colors.black} />
            </TouchableOpacity>
          ),
        }}
      />

      {error && (
        <Animated.View 
          entering={FadeInDown.duration(300)} 
          style={styles.errorContainer}
        >
          <Text style={styles.errorText}>{error}</Text>
          <TouchableOpacity onPress={clearError}>
            <Ionicons name="close-circle" size={20} color={Colors.red} />
          </TouchableOpacity>
        </Animated.View>
      )}

      <FlatList
        ref={flatListRef}
        data={messages}
        keyExtractor={(item) => item.id}
        renderItem={renderMessage}
        contentContainerStyle={[
            styles.messagesList, 
            keyboardVisible && Platform.OS === 'android' ? { paddingBottom: 100 } : {}
          ]}
        showsVerticalScrollIndicator={false}
        onLayout={() => flatListRef.current?.scrollToEnd({ animated: false })}
      />

        
        <View style={styles.inputContainer}>
          <TextInput
            style={styles.input}
            placeholder="Digite sua mensagem..."
            placeholderTextColor={Colors.gray}
            value={message}
            onChangeText={setMessage}
            multiline
            maxLength={500}
            onSubmitEditing={handleSendMessage}
            blurOnSubmit={false}
          />
          <TouchableOpacity 
            style={styles.sendButton} 
            onPress={handleSendMessage}
            disabled={isLoading || message.trim() === ''}
          >
            {isLoading ? (
              <ActivityIndicator size="small" color={Colors.white} />
            ) : (
              <Ionicons name="send" size={20} color={Colors.white} />
            )}
          </TouchableOpacity>
        </View>
      
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: Colors.background,
  },
  backButton: {
    marginLeft: 10,
  },
  errorContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    backgroundColor: Colors.transparentPrimary,
    marginHorizontal: 16,
    marginVertical: 8,
    paddingHorizontal: 12,
    paddingVertical: 8,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: Colors.red,
  },
  errorText: {
    color: Colors.red,
    flex: 1,
    fontSize: 13,
  },
  messagesList: {
    paddingHorizontal: 16,
    paddingBottom: Platform.OS === 'android' ? 20 : 16,
  },
  messageContainer: {
    maxWidth: '80%',
    marginVertical: 6,
    padding: 12,
    borderRadius: 16,
    elevation: 1,
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.1,
    shadowRadius: 1,
  },
  userMessage: {
    alignSelf: 'flex-end',
    backgroundColor: Colors.primary,
    borderBottomRightRadius: 4,
  },
  assistantMessage: {
    alignSelf: 'flex-start',
    backgroundColor: Colors.white,
    borderBottomLeftRadius: 4,
  },
  messageText: {
    fontSize: 15,
    lineHeight: 20,
  },
  userMessageText: {
    color: Colors.white,
  },
  assistantMessageText: {
    color: Colors.black,
  },
  timestamp: {
    fontSize: 10,
    marginTop: 4,
    alignSelf: 'flex-end',
    opacity: 0.7,
    color: Colors.extraLightGray,
  },
  inputContainer: {
    flexDirection: 'row',
    paddingHorizontal: 16,
    paddingVertical: Platform.OS === 'android' ? 8 : 12,
    backgroundColor: Colors.white,
    alignItems: 'center',
    borderTopWidth: 1,
    borderTopColor: Colors.extraLightGray,
    marginBottom: 55,
  },
  input: {
    flex: 1,
    backgroundColor: Colors.background,
    borderRadius: 20,
    paddingHorizontal: 16,
    paddingVertical: Platform.OS === 'android' ? 8 : 10,
    marginRight: 12,
    maxHeight: 80,
    color: Colors.black,
  },
  sendButton: {
    backgroundColor: Colors.primary,
    width: Platform.OS === 'android' ? 36 : 40,
    height: Platform.OS === 'android' ? 36 : 40,
    borderRadius: Platform.OS === 'android' ? 18 : 20,
    justifyContent: 'center',
    alignItems: 'center',
  },
});

export default ChatSupportScreen; 