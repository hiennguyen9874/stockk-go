import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface MessageStateType {
  message?: string;
}

const initialState: MessageStateType = {};

const messageSlice = createSlice({
  name: 'message',
  initialState,
  reducers: {
    setMessage: (state, { payload }: PayloadAction<string>) => ({
      ...state,
      message: payload,
    }),
    clearMessage: (state) => ({
      ...state,
      message: '',
    }),
  },
});

export const { setMessage, clearMessage } = messageSlice.actions;

export default messageSlice.reducer;
