import { createSlice } from '@reduxjs/toolkit';

import type { RootState } from 'app/store';
import { login } from 'app/services/auth';

const stateLocalStorage = localStorage.getItem('user');

export interface IUser {
  id: number;
  name: string;
  email: string;
  isActive: boolean;
  isSuperUser: boolean;
  verified: boolean;
}

export interface IToken {
  accessToken: string;
  refreshToken: string;
}

type AuthState = {
  token: IToken | null;
  user: IUser | null;
  isAuthenticated: boolean;
};

const emptyState: AuthState = {
  token: null,
  user: null,
  isAuthenticated: false,
} as AuthState;

// eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
const initialState = (
  stateLocalStorage ? JSON.parse(stateLocalStorage) : emptyState
) as AuthState;

const slice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    logout: () => initialState,
  },
  extraReducers: (builder) => {
    builder
      // .addMatcher(login.matchPending, (state, action) => {
      //   console.log('pending', action);
      // })
      .addMatcher(login.matchFulfilled, (state, action) => {
        state.token = {
          accessToken: action.payload.access_token,
          refreshToken: action.payload.refresh_token,
        };
        state.user = {
          id: action.payload.user.id,
          name: action.payload.user.name,
          email: action.payload.user.email,
          isActive: action.payload.user.is_active,
          isSuperUser: action.payload.user.is_superuser,
          verified: action.payload.user.verified,
        };
        state.isAuthenticated = true;

        localStorage.setItem('user', JSON.stringify(state));
      })
      .addMatcher(login.matchRejected, (state, action) => {
        state.token = null;
        state.user = null;
        state.isAuthenticated = false;

        localStorage.removeItem('user');
      });
  },
});

export const { logout } = slice.actions;

export default slice.reducer;

export const selectCurrentUser = (state: RootState) => state.auth.user;

export const selectIsAuthenticated = (state: RootState) =>
  state.auth.isAuthenticated;
