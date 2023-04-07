import { api } from './api';
import type { Response, ClientResponse, ClientCreate } from './types';

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
    deleteClient: builder.mutation<Response<ClientResponse>, number>({
      query: (id) => ({
        url: `client/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: (client) => [{ type: 'Client', id: client?.data.id }],
      extraOptions: { maxRetries: 0 },
    }),
    createClient: builder.mutation<Response<ClientResponse>, ClientCreate>({
      query: (clientCreate) => ({
        url: 'client',
        method: 'POST',
        body: clientCreate,
      }),
      invalidatesTags: ['Client'],
      extraOptions: { maxRetries: 0 },
    }),
  }),
});

export const {
  useGetClientsQuery,
  useDeleteClientMutation,
  useCreateClientMutation,
} = clientApi;

export const {
  endpoints: { getClients, deleteClient, createClient },
} = clientApi;
