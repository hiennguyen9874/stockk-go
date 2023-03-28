import { retry } from '@reduxjs/toolkit/query/react';

import { api } from './api';

interface UserSignIn {
  email: string;
  password: string;
}

interface UserResponse {
  id: number;
  name: string;
  email: string;
  created_at: string;
  updated_at: string;
  is_active: boolean;
  is_superuser: boolean;
  verified: boolean;
}

interface TokenResponse {
  access_token: string;
  refresh_token: string;
  token_type: string;
  user: UserResponse;
}

export const authApi = api.injectEndpoints({
  endpoints: (builder) => ({
    login: builder.mutation<TokenResponse, UserSignIn>({
      query: (credentials) => {
        const formData = new FormData();

        formData.append('username', credentials.email);
        formData.append('password', credentials.password);

        return {
          url: 'auth/login',
          method: 'POST',
          body: formData,
        };
      },
      extraOptions: {
        backoff: () => {
          // We intentionally error once on login, and this breaks out of retrying. The next login attempt will succeed.
          retry.fail({ fake: 'error' });
        },
      },
    }),
  }),
});

export const { useLoginMutation } = authApi;

export const {
  endpoints: { login },
} = authApi;
