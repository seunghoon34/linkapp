import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import * as api from '../../services/api';

interface MatchState {
  isSearching: boolean;
  currentMatch: any | null;
  chatrooms: any[];
  status: 'idle' | 'loading' | 'succeeded' | 'failed';
  error: string | null;
}

const initialState: MatchState = {
  isSearching: false,
  currentMatch: null,
  chatrooms: [],
  status: 'idle',
  error: null,
};

export const startSearching = createAsyncThunk(
  'match/startSearching',
  async (_, { getState }) => {
    const { auth } = getState() as { auth: { token: string } };
    const response = await api.startSearching(auth.token);
    return response;
  }
);

export const stopSearching = createAsyncThunk(
  'match/stopSearching',
  async (_, { getState }) => {
    const { auth } = getState() as { auth: { token: string } };
    const response = await api.stopSearching(auth.token);
    return response;
  }
);

export const findMatch = createAsyncThunk(
  'match/findMatch',
  async (_, { getState }) => {
    const { auth } = getState() as { auth: { token: string } };
    const response = await api.findMatch(auth.token);
    return response;
  }
);

export const respondToLink = createAsyncThunk(
  'match/respondToLink',
  async ({ linkId, accept }: { linkId: string; accept: boolean }, { getState }) => {
    const { auth } = getState() as { auth: { token: string } };
    const response = await api.respondToLink(auth.token, linkId, accept);
    return response;
  }
);

export const sendMessage = createAsyncThunk(
  'match/sendMessage',
  async ({ chatroomId, content }: { chatroomId: string; content: string }, { getState }) => {
    const { auth } = getState() as { auth: { token: string } };
    const response = await api.sendMessage(auth.token, chatroomId, content);
    return response;
  }
);

export const verifyNFCAndUnlockChatroom = createAsyncThunk(
  'match/verifyNFCAndUnlockChatroom',
  async ({ chatroomId }: { chatroomId: string }, { getState }) => {
    const { auth } = getState() as { auth: { token: string } };
    const response = await api.verifyNFCAndUnlockChatroom(auth.token, chatroomId);
    return response;
  }
);

const matchSlice = createSlice({
  name: 'match',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(startSearching.fulfilled, (state) => {
        state.isSearching = true;
        state.status = 'succeeded';
      })
      .addCase(stopSearching.fulfilled, (state) => {
        state.isSearching = false;
        state.status = 'succeeded';
      })
      .addCase(findMatch.fulfilled, (state, action) => {
        state.currentMatch = action.payload;
        state.status = 'succeeded';
      })
      .addCase(respondToLink.fulfilled, (state, action) => {
        if (action.payload.chatroom) {
          state.chatrooms.push(action.payload.chatroom);
        }
        state.currentMatch = null;
        state.status = 'succeeded';
      })
      .addCase(sendMessage.fulfilled, (state, action) => {
        const chatroom = state.chatrooms.find(room => room.id === action.payload.chatroomId);
        if (chatroom) {
          chatroom.messages.push(action.payload);
        }
        state.status = 'succeeded';
      })
      .addCase(verifyNFCAndUnlockChatroom.fulfilled, (state, action) => {
        const chatroom = state.chatrooms.find(room => room.id === action.payload.chatroomId);
        if (chatroom) {
          chatroom.isLocked = false;
        }
        state.status = 'succeeded';
      });
  },
});

export default matchSlice.reducer;