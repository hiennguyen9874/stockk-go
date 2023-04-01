import { api } from './api';
import type {
  Response,
  WatchListCreate,
  WatchListResponse,
  WatchListUpdate,
} from './types';

export const watchListsApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getWatchLists: builder.query<Response<WatchListResponse[]>, void>({
      query: () => ({
        url: 'watchlist',
        method: 'GET',
      }),
      providesTags: (result) =>
        result
          ? [
              ...result.data.map(({ id }) => ({
                type: 'WatchLists' as const,
                id,
              })),
              { type: 'WatchLists', id: 'LIST' },
            ]
          : [{ type: 'WatchLists', id: 'LIST' }],
    }),
    createWatchList: builder.mutation<
      Response<WatchListResponse>,
      WatchListCreate
    >({
      query: ({ name }) => ({
        url: 'watchlist',
        method: 'POST',
        body: {
          name,
          tickers: [],
        },
      }),
      invalidatesTags: ['WatchLists'],
    }),
    getWatchList: builder.query<Response<WatchListResponse>, number>({
      query: (id) => ({
        url: `watchlist/${id}`,
        method: 'GET',
      }),
      providesTags: (watchlist) => [
        { type: 'WatchLists', id: watchlist?.data.id },
      ],
    }),
    updateWatchList: builder.mutation<
      Response<WatchListResponse>,
      WatchListUpdate
    >({
      query: ({ id, name, tickers }) => ({
        url: `watchlist/${id}`,
        method: 'PUT',
        body: {
          name,
          tickers,
        },
      }),
      invalidatesTags: (watchlist) => [
        { type: 'WatchLists', id: watchlist?.data.id },
      ],
    }),
    deleteWatchList: builder.mutation<Response<WatchListResponse>, number>({
      query: (id) => ({
        url: `watchlist/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: (watchlist) => [
        { type: 'WatchLists', id: watchlist?.data.id },
      ],
    }),
  }),
});

export const {
  useGetWatchListsQuery,
  useCreateWatchListMutation,
  useGetWatchListQuery,
  useUpdateWatchListMutation,
  useDeleteWatchListMutation,
} = watchListsApi;

export const {
  endpoints: {
    getWatchLists,
    createWatchList,
    getWatchList,
    updateWatchList,
    deleteWatchList,
  },
} = watchListsApi;
