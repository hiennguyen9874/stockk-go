import { api } from './api';

interface Response<T> {
  data: T;
  is_success: boolean;
}

interface WatchListCreate {
  name: string;
}

interface WatchListResponse {
  id: number;
  created_at: string;
  updated_at: string;
  name: string;
  tickers: string[];
  owner_id: number;
}

export const watchListsApi = api.injectEndpoints({
  endpoints: (builder) => ({
    createWatchList: builder.mutation<WatchListResponse, WatchListCreate>({
      query: ({ name }) => ({
        url: 'watchlist',
        method: 'POST',
        body: {
          name,
          tickers: [],
        },
      }),
    }),
    getWatchLists: builder.query<Response<WatchListResponse[]>, void>({
      query: () => ({
        url: 'watchlist',
        method: 'GET',
      }),
    }),
  }),
});

export const { useCreateWatchListMutation, useGetWatchListsQuery } =
  watchListsApi;

export const {
  endpoints: { createWatchList, getWatchLists },
} = watchListsApi;
