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

export type RecordItem = {
  id: string;
  type: RecordType;
  name: string;
  value: string;
  ttl: number;
  extra: { [key: string]: any };
};

export const domainApiSlice = apiSlice.injectEndpoints({
  endpoints: (build) => ({
    list: build.query<ListApiResponse, void>({
      query: () => '/domain',
      providesTags: ['Domain'],
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

export const { useListQuery } = domainApiSlice;
