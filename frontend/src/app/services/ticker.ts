import { api } from './api';
import type { Response, TickerResponse, TickerSnapshotResponse } from './types';

export const tickerApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getTicker: builder.query<Response<TickerResponse>, string>({
      query: (symbol) => ({
        url: `ticker/${symbol}`,
        method: 'GET',
      }),
      providesTags: (result) => [{ type: 'Ticker', id: result?.data.id }],
    }),
    searchBySymbol: builder.query<Response<TickerResponse[]>, string>({
      query: (symbol) => ({
        url: 'ticker/search',
        method: 'POST',
        params: {
          symbol,
        },
      }),
      providesTags: (result) =>
        result
          ? [
              ...result.data.map(({ id }) => ({
                type: 'Ticker' as const,
                id,
              })),
              { type: 'Ticker', id: 'LIST' },
            ]
          : [{ type: 'Ticker', id: 'LIST' }],
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

export const {
  useGetTickerQuery,
  useGetTickerSnapshotQuery,
  useSearchBySymbolQuery,
} = tickerApi;

export const {
  endpoints: { getTicker, getTickerSnapshot, searchBySymbol },
} = tickerApi;
