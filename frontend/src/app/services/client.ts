import { api } from './api';
import type { Response, ClientResponse } from './types';

export const clientApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getClients: builder.query<Response<ClientResponse[]>, void>({
      query: () => ({
        url: 'client',
        method: 'GET',
      }),
      providesTags: (result) =>
        result
          ? [
              ...result.data.map(({ id }) => ({
                type: 'Client' as const,
                id,
              })),
              { type: 'Client', id: 'LIST' },
            ]
          : [{ type: 'Client', id: 'LIST' }],
    }),
  }),
});

export const { useGetClientsQuery } = clientApi;

export const {
  endpoints: { getClients },
} = clientApi;
