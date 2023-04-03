import { api } from './api';
import type { Response, StockSnapshotResponse } from './types';

export const stockSnapshotApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getStockSnapshot: builder.query<Response<StockSnapshotResponse>, string>({
      query: (symbol) => ({
        url: `stocksnapshot/${symbol}`,
        method: 'GET',
      }),
      providesTags: (stockSnapshot) => [
        { type: 'StockSnapshot', id: stockSnapshot?.data.ticker },
      ],
    }),
  }),
});

export const { useGetStockSnapshotQuery } = stockSnapshotApi;

export const {
  endpoints: { getStockSnapshot },
} = stockSnapshotApi;
