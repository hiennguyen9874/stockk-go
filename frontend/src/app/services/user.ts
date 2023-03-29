import { api } from './api';

interface Response<T> {
  data: T;
  is_success: boolean;
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
