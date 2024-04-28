import { ApiResponse, ListApiResponse, apiSlice } from './api';

export interface DomainItem {
  id: number;
  provider_name: string;
  domain: string;
}

export interface RecordExtraField {
  key: string;
  title: string;
  field_type: 'text' | 'number' | 'switch' | 'select';
  options?: string[];
  default?: any;
}

export type RecordType = 'A' | 'AAAA' | 'CNAME' | 'TXT' | 'MX' | 'NS';

export const RecordTypes = ['A', 'AAAA', 'CNAME', 'TXT', 'MX', 'NS'];

export interface RecordItem {
  id?: string;
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
    domainDelete: build.mutation<ApiResponse, number>({
      query(id) {
        return {
          url: `/domain/${id}`,
          method: 'DELETE',
        };
      },
      invalidatesTags: ['Domain', 'Record'],
    }),
    recordList: build.query<
      ApiResponse<{ list: RecordItem[]; extra: RecordExtraField[] }>,
      number
    >({
      query: (id: number) => '/domain/' + id + '/record',
      providesTags: ['Record'],
    }),
    recordCreate: build.mutation<
      ApiResponse,
      { domain_id: number; record: RecordItem }
    >({
      query(params) {
        return {
          url: `/domain/${params.domain_id}/record`,
          method: 'POST',
          body: params.record,
        };
      },
      invalidatesTags: (result, error, params) => {
        if (error) {
          return [];
        }
        return ['Record'];
      },
    }),
    recordUpdate: build.mutation<
      ApiResponse,
      { domain_id: number; record_id: string; record: RecordItem }
    >({
      query(params) {
        return {
          url: `/domain/${params.domain_id}/record/${params.record_id}`,
          method: 'PUT',
          body: params.record,
        };
      },
      invalidatesTags: (result, error, params) => {
        if (error) {
          return [];
        }
        return ['Record'];
      },
    }),
    recordDelete: build.mutation<
      ApiResponse,
      { domain_id: number; record_id: string }
    >({
      query(params) {
        return {
          url: `/domain/${params.domain_id}/record/${params.record_id}`,
          method: 'DELETE',
        };
      },
      invalidatesTags: (result, error, params) => {
        if (error) {
          return [];
        }
        return ['Record'];
      },
    }),
  }),
});

export const {
  useDomainListQuery,
  useRecordListQuery,
  useDomainQueryQuery,
  useDomainAddMutation,
  useDomainDeleteMutation,
  useLazyRecordListQuery,
  useRecordDeleteMutation,
  useRecordCreateMutation,
  useRecordUpdateMutation,
} = domainApiSlice;
