import { ApiResponse, ListApiResponse, apiSlice } from './api';

export type CertItem = {
  id: number;
  email: string;
  domains: string[];
  status: number;
  createtime: number;
};

export type CertHostingItem = {
  id: number;
  cdn_id: number;
  cdn: string;
  cert_id: number;
  status: number;
};

export type CertHostingQueryItem = {
  id: number;
  domain: string;
  is_managed: boolean;
};

export const certApiSlice = apiSlice.injectEndpoints({
  endpoints: (build) => ({
    certList: build.query<ListApiResponse<CertItem>, void>({
      query: () => '/cert',
      providesTags: ['Cert'],
    }),
    certCreate: build.mutation<
      ApiResponse,
      { email: string; domains: string[] }
    >({
      query(args) {
        return {
          url: '/cert',
          method: 'POST',
          body: args,
        };
      },
      invalidatesTags: ['Cert'],
    }),
    certDelete: build.mutation<ApiResponse, number>({
      query(id) {
        return {
          url: `/cert/${id}`,
          method: 'DELETE',
        };
      },
      invalidatesTags: ['Cert'],
    }),
    certHostigList: build.query<ListApiResponse<CertHostingItem>, void>({
      query: () => '/cert/hosting',
      providesTags: ['CertHosting'],
    }),
    certHostingAdd: build.mutation<
      ApiResponse,
      {
        email: string;
        type: 1 | 2;
        cdn_id?: number;
        provider_id?: number;
        config?: Record<string, string>;
      }
    >({
      query(args) {
        return {
          url: '/cert/hosting',
          method: 'POST',
          body: args,
        };
      },
      invalidatesTags: ['CertHosting'],
    }),
    certHostingDelete: build.mutation<ApiResponse, number>({
      query(id) {
        return {
          url: `/cert/hosting/${id}`,
          method: 'DELETE',
        };
      },
      invalidatesTags: ['CertHosting'],
    }),
    certHostingQuery: build.query<
      ApiResponse<{ list: CertHostingQueryItem[] }>,
      void
    >({
      query() {
        return `/cert/hosting/query`;
      },
    }),
  }),
});

export const {
  useCertListQuery,
  useCertCreateMutation,
  useCertDeleteMutation,
  useCertHostigListQuery,
  useCertHostingAddMutation,
  useCertHostingDeleteMutation,
  useCertHostingQueryQuery,
} = certApiSlice;
