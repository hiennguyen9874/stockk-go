/* eslint-disable @typescript-eslint/no-unsafe-assignment */
import { createApi, fetchBaseQuery, retry } from '@reduxjs/toolkit/query/react';
import type {
  BaseQueryFn,
  FetchArgs,
  FetchBaseQueryError,
} from '@reduxjs/toolkit/query';

import type { RootState } from 'app/store';

import type { TokenResponse } from './types';

// Create our baseQuery instance
const baseQuery = fetchBaseQuery({
  baseUrl: 'https://stockk.dscilab.site:20007/api',
  credentials: 'include',
  prepareHeaders: (headers, { getState }) => {
    // By default, if we have a token in the store, let's use that for authenticated requests
    const { token } = (getState() as RootState).auth;
    if (token) {
      headers.set('authorization', `Bearer ${token.accessToken}`);
    }
    return headers;
  },
});

const baseQueryWithReauth: BaseQueryFn<
  string | FetchArgs,
  unknown,
  FetchBaseQueryError
> = async (args, api, extraOptions) => {
  let result = await baseQuery(args, api, extraOptions);
  if (result.error && result.error.status === 401) {
    const { token } = (api.getState() as RootState).auth;

    const headers = token
      ? {
          authorization: token.refreshToken,
        }
      : undefined;

    // try to get a new token
    const refreshResult = await baseQuery(
      {
        url: 'auth/refresh',
        method: 'GET',
        headers,
      },
      api,
      extraOptions
    );

    if (refreshResult.data) {
      // store the new token
      api.dispatch({
        type: 'auth/setTokenUser',
        payload: {
          token: {
            accessToken: (refreshResult.data as TokenResponse).access_token,
            refreshToken: (refreshResult.data as TokenResponse).refresh_token,
          },
          user: {
            id: (refreshResult.data as TokenResponse).user.id,
            name: (refreshResult.data as TokenResponse).user.name,
            email: (refreshResult.data as TokenResponse).user.email,
            isActive: (refreshResult.data as TokenResponse).user.is_active,
            isSuperUser: (refreshResult.data as TokenResponse).user
              .is_superuser,
            verified: (refreshResult.data as TokenResponse).user.verified,
          },
        },
      });
      // retry the initial query
      result = await baseQuery(args, api, extraOptions);
    } else {
      api.dispatch({
        type: 'auth/logout',
        payload: undefined,
      });
    }
  }
  return result;
};

const baseQueryWithRetry = retry(baseQueryWithReauth, { maxRetries: 6 });

/**
 * Create a base API to inject endpoints into elsewhere.
 * Components using this API should import from the injected site,
 * in order to get the appropriate types,
 * and to ensure that the file injecting the endpoints is loaded
 */
export const api = createApi({
  /**
   * `reducerPath` is optional and will not be required by most users.
   * This is useful if you have multiple API definitions,
   * e.g. where each has a different domain, with no interaction between endpoints.
   * Otherwise, a single API definition should be used in order to support tag invalidation,
   * among other features
   */
  reducerPath: 'splitApi',
  /**
   * A bare bones base query would just be `baseQuery: fetchBaseQuery({ baseUrl: '/' })`
   */
  baseQuery: baseQueryWithRetry,
  /**
   * Tag types must be defined in the original API definition
   * for any tags that would be provided by injected endpoints
   */
  tagTypes: ['Auth', 'User', 'WatchLists'],
  /**
   * This api has endpoints injected in adjacent files,
   * which is why no endpoints are shown below.
   * If you want all endpoints defined in the same file, they could be included here instead
   */
  endpoints: () => ({}),
});

export const enhancedApi = api.enhanceEndpoints({
  endpoints: () => ({
    getPost: () => 'test',
  }),
});
