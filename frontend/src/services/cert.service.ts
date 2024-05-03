import { ApiResponse, ListApiResponse, apiSlice } from './api';

export type CertItem = {
  id: number;
  email: string;
  domains: string[];
  status: number;
  createtime: number;
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
  }),
});

export const { useCertListQuery, useCertCreateMutation } = certApiSlice;
