import { ApiResponse, ListApiResponse, apiSlice } from './api';

export interface DomainItem {
  id: number;
  provider_name: string;
  name: string;
}

export interface RecordExtraField {
  key: string;
  title: string;
  field_type: 'text' | 'number' | 'switch' | 'select';
  options?: string[];
  default?: any;
}

export type RecordType = 'A' | 'AAAA' | 'CNAME' | 'TXT';

export interface RecordItem {
  id: string;
  type: RecordType;
  name: string;
  value: string;
  ttl: number;
  extra: { [key: string]: any };
}

export interface QueryDomainItem {
  provider_id: number;
  provider_name: string;
  domain_id: string;
  domain: string;
  is_managed: boolean;
}

export const domainApiSlice = apiSlice.injectEndpoints({
  endpoints: (build) => ({
    domainList: build.query<ListApiResponse, void>({
      query: () => '/domain',
      providesTags: ['Domain'],
    }),
    domainQuery: build.query<ApiResponse<{ items: QueryDomainItem[] }>, void>({
      query: () => '/domain/query',
      providesTags: ['Domain'],
    }),
    domainAdd: build.mutation<
      ApiResponse,
      { provider_id: number; domain_id: string; domain: string }
    >({
      query(params) {
        return {
          url: '/domain',
          method: 'POST',
          body: params,
        };
      },
      invalidatesTags: ['Domain'],
    }),
    domainDelete: build.mutation<ApiResponse,  number >({
      query(id) {
        return {
          url: `/domain/${id}`,
          method: 'DELETE',
        };
      },
      invalidatesTags: ['Domain', 'Record'],
    }),
    recordList: build.query<
      ListApiResponse<RecordItem, { extra_fields: RecordExtraField[] }>,
      void
    >({
      query: () => '/domain/record',
      providesTags: ['Record'],
    }),
  }),
});

export const {
  useDomainListQuery,
  useRecordListQuery,
  useDomainQueryQuery,
  useDomainAddMutation,
  useDomainDeleteMutation,
} = domainApiSlice;
