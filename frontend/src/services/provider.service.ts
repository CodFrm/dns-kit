import { ApiResponse, ListApiResponse, apiSlice } from './api';

export interface ProviderItem {
  id: number;
  name: string;
  platform: string;
}

export type ProviderPlatform = 'tencent' | 'cloudflare';

export const providerApiSlice = apiSlice.injectEndpoints({
  endpoints: (build) => ({
    providerList: build.query<ListApiResponse, void>({
      query: () => '/provider',
      providesTags: ['Provider'],
    }),
    providerEdit: build.mutation<
      ApiResponse,
      {
        id?: number;
        name: string;
        platform: ProviderPlatform;
        secret: { [key: string]: string };
      }
    >({
      query(params) {
        return {
          url: '/provider' + (params.id ? '/' + params.id : ''),
          method: params.id ? 'PUT' : 'POST',
          body: params,
        };
      },
      invalidatesTags: (result, error, params) => {
        if (error) {
          return [];
        }
        return ['Provider'];
      },
    }),
    providerDelete: build.mutation<ApiResponse, number>({
      query(id) {
        return {
          url: '/provider/' + id,
          method: 'DELETE',
        };
      },
      invalidatesTags: ['Provider'],
    }),
  }),
});

export const {
  useProviderListQuery,
  useProviderEditMutation,
  useProviderDeleteMutation,
} = providerApiSlice;
