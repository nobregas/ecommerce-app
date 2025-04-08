import React, { useState, useEffect, useRef } from 'react';
import {
  View,
  Text,
  TouchableOpacity,
  StyleSheet,
  ScrollView,
  KeyboardAvoidingView,
  Platform,
  SafeAreaView,
  ActivityIndicator,
  Keyboard
} from 'react-native';
import Animated, { SlideInDown } from 'react-native-reanimated';
import { aiChatService } from '@/service/aiChatService';
import { ChatMessageType } from '@/service/aiChatService';
import { Colors } from '@/constants/Colors';
import InputField from '@/components/InputField';
import { Stack } from 'expo-router';
import { useKeyboard } from '@react-native-community/hooks';
import { Dimensions } from 'react-native';

const ChatScreen = () => {
  const [messages, setMessages] = useState<ChatMessageType[]>([]);
  const [inputText, setInputText] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const scrollViewRef = useRef<ScrollView>(null);
  const keyboard = useKeyboard();
  const { height } = Dimensions.get('window');

  const sendMessage = async (messageContent: string) => {
    if (!messageContent.trim()) return;

    const userMessage: ChatMessageType = {
      id: Date.now().toString(),
      content: messageContent,
      role: 'user',
      timestamp: new Date()
    };

    setMessages(prev => [...prev, userMessage]);
    setInputText('');
    Keyboard.dismiss();

    try {
      setIsLoading(true);
      const aiResponse = await aiChatService.sendMessage(messageContent);

      if (aiResponse) {
        setMessages(prev => [...prev, aiResponse]);
      }
    } catch (error) {
      console.error('Error sending message:', error);
      const errorMessage: ChatMessageType = {
        id: Date.now().toString(),
        content: 'Desculpe, ocorreu um erro. Por favor, tente novamente.',
        role: 'assistant',
        timestamp: new Date()
      };
      setMessages(prev => [...prev, errorMessage]);
    } finally {
      setIsLoading(false);
    }
  };

  const initialWelcomeMessage: ChatMessageType = {
    id: 'welcome',
    content: 'Olá! Sou a assistente virtual da ShopX. Como posso ajudar você hoje?',
    role: 'assistant',
    timestamp: new Date()
  };

  const MessageBubble = ({ message }: { message: ChatMessageType }) => (
    <View style={[
      styles.messageBubble,
      message.role === 'user' ? styles.userBubble : styles.aiBubble
    ]}>
      <Text style={message.role === 'user' ? styles.userText : styles.aiText}>
        {message.content}
      </Text>
      <Text style={styles.timestamp}>
        {message.timestamp.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
      </Text>
    </View>
  );

  return (
    <>
      <Stack.Screen
        options={{ headerShown: true, headerTransparent: true, headerTitleAlign: 'center', title: 'support' }}
      />
      <SafeAreaView style={styles.container}>
        <KeyboardAvoidingView
          behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
          style={styles.flex}
          keyboardVerticalOffset={Platform.OS === 'ios' ? 90 : 0}
        >
          <View style={styles.header}>
            <View>
              <Text style={styles.headerTitle}>Suporte ShopX</Text>
              <View style={styles.headerStatus}>
                <View style={styles.statusIndicator} />
                <Text style={styles.statusText}>Online agora</Text>
              </View>
            </View>
          </View>

          <ScrollView
            ref={scrollViewRef}
            contentContainerStyle={[styles.messagesContainer, {
              paddingBottom: keyboard.keyboardShown ? height * 0.4 : 140
            }]}
            onContentSizeChange={() => scrollViewRef.current?.scrollToEnd({ animated: true })}
            keyboardDismissMode="interactive"
          >
            {messages.length === 0 ? (
              <MessageBubble message={initialWelcomeMessage} />
            ) : (
              messages.map((message) => (
                <MessageBubble key={message.id} message={message} />
              ))
            )}

            {isLoading && (
              <View style={[styles.messageBubble, styles.aiBubble]}>
                <ActivityIndicator size="small" color={Colors.primary} />
              </View>
            )}
          </ScrollView>

          <Animated.View
            entering={SlideInDown.duration(500)}
            style={[
              styles.footer,
              { bottom: keyboard.keyboardShown ? (Platform.OS === 'ios' ? keyboard.keyboardHeight - 30 : 20) : 0 }
            ]}
          >
            <InputField
              placeholder="Digite sua mensagem..."
              value={inputText}
              onChangeText={setInputText}
              style={styles.footerInput}
              multiline
              onSubmitEditing={() => sendMessage(inputText)}
              blurOnSubmit={false}
            />
            <TouchableOpacity
              style={[styles.footerButton, !inputText && styles.disabledButton]}
              onPress={() => sendMessage(inputText)}
              disabled={!inputText || isLoading}
            >
              <Text style={styles.footerButtonText}>Enviar</Text>
            </TouchableOpacity>
          </Animated.View>

        </KeyboardAvoidingView>
      </SafeAreaView>
    </>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: Colors.semiTransparentWhite,
  },
  flex: {
    flex: 1,
  },
  header: {
    padding: 16,
    borderBottomWidth: 1,
    borderBottomColor: Colors.extraLightGray,
    backgroundColor: Colors.solidWhite,
  },
  headerTitle: {
    fontSize: 20,
    fontWeight: '600',
    color: Colors.black,
    fontFamily: Platform.OS === 'ios' ? 'System' : 'Roboto',
  },
  headerStatus: {
    flexDirection: 'row',
    alignItems: 'center',
    marginTop: 4,
  },
  statusIndicator: {
    width: 8,
    height: 8,
    borderRadius: 4,
    backgroundColor: Colors.primary,
    marginRight: 8,
  },
  statusText: {
    color: Colors.gray,
    fontSize: 14,
  },
  messagesContainer: {
    padding: 16,
    paddingBottom: 24,
  },
  messageBubble: {
    maxWidth: '80%',
    borderRadius: 18,
    padding: 12,
    marginBottom: 12,
    shadowColor: Colors.black,
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.1,
    shadowRadius: 2,
    elevation: 2,
  },
  userBubble: {
    alignSelf: 'flex-end',
    backgroundColor: Colors.primary,
    borderBottomRightRadius: 4,
  },
  aiBubble: {
    alignSelf: 'flex-start',
    backgroundColor: Colors.white,
    borderBottomLeftRadius: 4,
  },
  userText: {
    color: Colors.white,
    fontSize: 16,
  },
  aiText: {
    color: Colors.black,
    fontSize: 16,
  },
  timestamp: {
    fontSize: 10,
    color: Colors.lightGray,
    marginTop: 4,
    alignSelf: 'flex-end',
  },
  inputContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: 12,
    backgroundColor: Colors.solidWhite,
    borderTopWidth: 1,
    borderTopColor: Colors.extraLightGray,
  },
  chatInput: {
    flex: 1,
    minHeight: 40,
    maxHeight: 120,
    borderRadius: 20,
    marginRight: 12,
    backgroundColor: Colors.background,
    borderColor: Colors.extraLightGray,
    paddingVertical: 8,
    paddingHorizontal: 16,
  },
  sendButton: {
    paddingVertical: 8,
    paddingHorizontal: 16,
    borderRadius: 18,
    backgroundColor: Colors.extraLightGray,
  },
  activeSendButton: {
    backgroundColor: Colors.primary,
  },
  sendButtonText: {
    color: Colors.white,
    fontWeight: '500',
  },
  footer: {
    position: 'absolute',
    left: 0,
    right: 0,
    padding: 12,
    backgroundColor: Colors.white,
    borderTopWidth: 1,
    borderTopColor: Colors.extraLightGray,
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    shadowColor: Colors.black,
    shadowOffset: { width: 0, height: -2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 5,
  },
  footerInput: {
    flex: 1,
    minHeight: 40,
    maxHeight: 100,
    backgroundColor: Colors.background,
    borderRadius: 15,
    paddingHorizontal: 15,
    paddingVertical: 10,
    fontSize: 16,
    textAlignVertical: 'center',
    marginTop: 12
  },
  footerButton: {
    backgroundColor: Colors.primary,
    borderRadius: 20,
    paddingHorizontal: 20,
    height: 40,
    justifyContent: 'center',
    alignItems: 'center',
  },

  footerButtonText: {
    color: Colors.white,
    fontWeight: '500',
    fontSize: 16,
  },
  disabledButton: {
    backgroundColor: Colors.lightGray,
  },
});

export default ChatScreen;