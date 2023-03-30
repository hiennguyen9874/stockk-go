import { api } from './api';
import type { Response, UserResponse } from './types';

export const userApi = api.injectEndpoints({
  endpoints: (builder) => ({
    userMe: builder.query<Response<UserResponse>, void>({
      query: () => ({
        url: 'user/me',
        method: 'GET',
      }),
    }),
  }),
});

export const { useUserMeQuery } = userApi;

export const {
  endpoints: { userMe },
} = userApi;
