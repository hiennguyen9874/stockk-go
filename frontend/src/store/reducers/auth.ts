import { createSlice, PayloadAction } from '@reduxjs/toolkit';

import { User } from 'api/auth-header';

const userJson = localStorage.getItem('user');

export interface AuthStateType {
  isLoggedIn: boolean;
  user: User | null;
}

const initialState: AuthStateType = userJson
  ? // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
    { isLoggedIn: true, user: JSON.parse(userJson) }
  : { isLoggedIn: false, user: null };

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    registerSuccess: (state) => ({
      ...state,
      isLoggedIn: true,
    }),
    registerFailure: (state) => ({
      ...state,
      isLoggedIn: false,
    }),
    loginSuccess: (state, { payload }: PayloadAction<User>) => ({
      ...state,
      isLoggedIn: true,
      user: payload,
    }),
    loginFailure: (state) => ({
      ...state,
      isLoggedIn: false,
      user: null,
    }),
    logout: (state) => ({
      ...state,
      isLoggedIn: false,
      user: null,
    }),
  },
});

export const {
  registerSuccess,
  registerFailure,
  loginSuccess,
  loginFailure,
  logout,
} = authSlice.actions;

export default authSlice.reducer;
