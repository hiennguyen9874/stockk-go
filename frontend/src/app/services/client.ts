import { api } from './api';
import type {
  Response,
  ClientResponse,
  ClientCreate,
  ClientUpdate,
} from './types';

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
    updateClient: builder.mutation<Response<ClientResponse>, ClientUpdate>({
      query: (clientUpdate) => ({
        url: `client/${clientUpdate.id}`,
        method: 'PUT',
        body: {
          current_ticker: clientUpdate.current_ticker,
          current_resolution: clientUpdate.current_resolution,
        },
      }),
      invalidatesTags: (client) => [{ type: 'Client', id: client?.data.id }],
    }),
  }),
});

export const {
  useGetClientsQuery,
  useDeleteClientMutation,
  useCreateClientMutation,
  useUpdateClientMutation,
} = clientApi;

export const {
  endpoints: { getClients, deleteClient, createClient, updateClient },
} = clientApi;
