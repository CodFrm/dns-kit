import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import { ApiResponse } from './api';
import { UserInfo } from '@/store/global';

export type LoginResponse = {
  token: {
    access_token: string;
    refresh_token: string;
    expire: number;
    refresh_expire: number;
  };
  username: string;
};

export type CurrentUserResponse = {
  username: string;
};

export const userApiSlice = createApi({
  baseQuery: fetchBaseQuery({
    baseUrl: '/api/v1/user',
  }),
  reducerPath: 'api',
  // Tag types are used for caching and invalidation.
  tagTypes: ['User'],
  endpoints: (build) => ({
    currentUser: build.query<ApiResponse<CurrentUserResponse>, void>({
      query: () => '/current',
      providesTags: ['User'],
    }),
    login: build.mutation<ApiResponse, void>({
      query: (params) => ({
        url: '/login',
        method: 'POST',
        body: params,
      }),
      
      invalidatesTags: ['User'],
    }),
    logout: build.mutation<ApiResponse<LoginResponse>, void>({
      query() {
        return {
          url: '/logout',
          method: 'DELETE',
        };
      },
      invalidatesTags: ['User'],
    }),
  }),
});

export const {
  useCurrentUserQuery,
  useLazyCurrentUserQuery,
  useLogoutMutation,
  useLoginMutation,
} = userApiSlice;
