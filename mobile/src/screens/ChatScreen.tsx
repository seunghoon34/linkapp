import React, { useState, useEffect } from 'react';
import { View, Text, FlatList, TextInput, TouchableOpacity, StyleSheet } from 'react-native';
import { useSelector, useDispatch } from 'react-redux';
import { AppDispatch, RootState } from '../store/store';
import { sendMessage, verifyNFCAndUnlockChatroom } from '../store/slices/matchSlice';
import NFCService from '../services/NFCService';
import { RouteProp } from '@react-navigation/native';
import { RootStackParamList } from '../types/navigation';

type ChatScreenProps = {
  route: RouteProp<RootStackParamList, 'Chat'>;
};

const ChatScreen: React.FC<ChatScreenProps> = ({ route }) => {
  const { chatId: chatroomId } = route.params;
  const dispatch = useDispatch<AppDispatch>();
  const [message, setMessage] = useState('');
  const chatroom = useSelector((state: RootState) => 
    state.match.chatrooms.find(room => room.id === chatroomId)
  );

  const sendMessageHandler = async () => {
    if (message.trim()) {
      try {
        await dispatch(sendMessage({ chatroomId, content: message })).unwrap();
        setMessage('');
      } catch (error) {
        console.error('Error sending message:', error);
      }
    }
  };

  const unlockChatroomHandler = async () => {
    try {
      await NFCService.init();
      await NFCService.readNFC(); // We still read the NFC, but don't pass it to the API
      await dispatch(verifyNFCAndUnlockChatroom({ chatroomId })).unwrap();
    } catch (error: unknown) {
      if (error instanceof Error) {
        console.error('Error unlocking chatroom:', error.message);
      } else {
        console.error('Unknown error unlocking chatroom');
      }
    }
  };

  return (
    <View style={styles.container}>
      {chatroom?.isLocked && (
        <View>
          <TouchableOpacity style={styles.unlockButton} onPress={unlockChatroomHandler}>
            <Text style={styles.unlockButtonText}>Unlock with NFC</Text>
          </TouchableOpacity>
        </View>
      )}
      <FlatList
        data={chatroom?.messages || []}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => (
          <View style={styles.messageContainer}>
            <Text style={styles.sender}>{item.sender}</Text>
            <Text style={styles.messageText}>{item.text}</Text>
          </View>
        )}
      />
      <View style={styles.inputContainer}>
        <TextInput
          style={styles.input}
          value={message}
          onChangeText={setMessage}
          placeholder="Type a message..."
          editable={!chatroom?.isLocked}
        />
        <TouchableOpacity 
          style={[styles.sendButton, chatroom?.isLocked && styles.disabledButton]} 
          onPress={sendMessageHandler}
          disabled={chatroom?.isLocked}
        >
          <Text style={styles.sendButtonText}>Send</Text>
        </TouchableOpacity>
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 10,
  },
  lockedContainer: {
    backgroundColor: '#FFD700',
    padding: 10,
    marginBottom: 10,
    borderRadius: 5,
  },
  lockedText: {
    fontSize: 16,
    textAlign: 'center',
    marginBottom: 10,
  },
  unlockButton: {
    backgroundColor: '#4CAF50',
    padding: 10,
    borderRadius: 5,
  },
  unlockButtonText: {
    color: 'white',
    textAlign: 'center',
    fontWeight: 'bold',
  },
  messageContainer: {
    marginBottom: 10,
  },
  sender: {
    fontWeight: 'bold',
  },
  messageText: {
    fontSize: 16,
  },
  inputContainer: {
    flexDirection: 'row',
    padding: 10,
  },
  input: {
    flex: 1,
    borderColor: 'gray',
    borderWidth: 1,
    borderRadius: 5,
    paddingHorizontal: 10,
    marginRight: 10,
  },
  sendButton: {
    backgroundColor: '#007AFF',
    padding: 10,
    borderRadius: 5,
    justifyContent: 'center',
  },
  disabledButton: {
    backgroundColor: '#ccc',
  },
  sendButtonText: {
    color: 'white',
    fontWeight: 'bold',
  },
});

export default ChatScreen;