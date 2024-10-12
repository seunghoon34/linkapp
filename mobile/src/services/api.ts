import axios from 'axios';

const API_URL = 'http://your-backend-url.com/api'; // Replace with your actual API URL

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const login = async (email: string, password: string) => {
  const response = await api.post('/login', { email, password });
  return response.data;
};

export const register = async (username: string, email: string, password: string) => {
  const response = await api.post('/register', { username, email, password });
  return response.data;
};

export const getProfile = async (token: string) => {
  const response = await api.get('/profile', {
    headers: { Authorization: `Bearer ${token}` }
  });
  return response.data;
};

export const updateProfile = async (token: string, profileData: any) => {
  const response = await api.put('/profile', profileData, {
    headers: { Authorization: `Bearer ${token}` }
  });
  return response.data;
};

export const getMatches = async (token: string) => {
  const response = await api.get('/matches', {
    headers: { Authorization: `Bearer ${token}` }
  });
  return response.data;
};

export const sendMessage = async (token: string, recipientId: string, content: string) => {
  const response = await api.post('/messages', { recipientId, content }, {
    headers: { Authorization: `Bearer ${token}` }
  });
  return response.data;
};

export const startSearching = async (token: string) => {
  const response = await api.post('/users/start-searching', {}, {
    headers: { Authorization: `Bearer ${token}` }
  });
  return response.data;
};

export const stopSearching = async (token: string) => {
  const response = await api.post('/users/stop-searching', {}, {
    headers: { Authorization: `Bearer ${token}` }
  });
  return response.data;
};

export const findMatch = async (token: string) => {
  const response = await api.get('/users/find-match', {
    headers: { Authorization: `Bearer ${token}` }
  });
  return response.data;
};

export const respondToLink = async (token: string, linkId: string, accept: boolean) => {
  const response = await api.post(`/users/links/${linkId}/respond`, { accept }, {
    headers: { Authorization: `Bearer ${token}` }
  });
  return response.data;
};

export const unlockChatroom = async (token: string, chatroomId: string) => {
  const response = await api.post(`/chatrooms/${chatroomId}/unlock`, {}, {
    headers: { Authorization: `Bearer ${token}` }
  });
  return response.data;
};

export const verifyNFCAndUnlockChatroom = async (token: string, chatroomId: string) => {
  const response = await api.post(`/users/chatrooms/${chatroomId}/nfc-unlock`, {}, {
    headers: { Authorization: `Bearer ${token}` }
  });
  return response.data;
};

export default api;