import { ApiResponse, ListApiResponse, apiSlice } from './api';

export type CDNItem = {
  id: number;
  provider_name: string;
  domain: string;
};

export type QueryCDNItem = {
  provider_id: number;
  provider_name: string;
  id: string;
  domain: string;
  is_managed: boolean;
};

export const cdnApiSlice = apiSlice.injectEndpoints({
  endpoints: (build) => ({
    cdnList: build.query<ListApiResponse<CDNItem>, void>({
      query: () => '/cdn',
      providesTags: ['CDN'],
    }),
    cdnQuery: build.query<ApiResponse<{ items: QueryCDNItem[] }>, void>({
      query: () => '/cdn/query',
      providesTags: ['CDN'],
    }),
    cdnAdd: build.mutation<
      ApiResponse,
      { provider_id: number; id: string; domain: string }
    >({
      query(params) {
        return {
          url: '/cdn',
          method: 'POST',
          body: params,
        };
      },
      invalidatesTags: ['CDN'],
    }),
    cdnDelete: build.mutation<ApiResponse, number>({
      query(id) {
        return {
          url: `/cdn/${id}`,
          method: 'DELETE',
        };
      },
      invalidatesTags: ['CDN'],
    }),
  }),
});

export const {
  useCdnListQuery,
  useCdnQueryQuery,
  useCdnAddMutation,
  useCdnDeleteMutation,
} = cdnApiSlice;
