import { api } from './api';
import type { Response, TickerResponse, TickerSnapshotResponse } from './types';

export const tickerApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getTicker: builder.query<Response<TickerResponse>, string>({
      query: (symbol) => ({
        url: `ticker/${symbol}`,
        method: 'GET',
      }),
      providesTags: (ticker) => [{ type: 'Ticker', id: ticker?.data.id }],
    }),
    getTickerSnapshot: builder.query<Response<TickerSnapshotResponse>, string>({
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

export const { useGetTickerQuery, useGetTickerSnapshotQuery } = tickerApi;

export const {
  endpoints: { getTicker, getTickerSnapshot },
} = tickerApi;
