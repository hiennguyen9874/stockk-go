import { api } from './api';
import type { Response, TickerResponse } from './types';

export const tickerApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getTicker: builder.query<Response<TickerResponse>, string>({
      query: (symbol) => ({
        url: `ticker/${symbol}`,
        method: 'GET',
      }),
      providesTags: (ticker) => [{ type: 'Ticker', id: ticker?.data.id }],
    }),
  }),
});

export const { useGetTickerQuery } = tickerApi;

export const {
  endpoints: { getTicker },
} = tickerApi;
