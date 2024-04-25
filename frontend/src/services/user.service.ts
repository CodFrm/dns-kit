import { ApiResponse, apiSlice } from './api';

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

export const userApiSlice = apiSlice.injectEndpoints({
  endpoints: (build) => ({
    currentUser: build.query<ApiResponse<CurrentUserResponse>, void>({
      query: () => '/user/current',
      providesTags: ['User'],
    }),
    login: build.mutation<ApiResponse, void>({
      query: (params) => ({
        url: '/user/login',
        method: 'POST',
        body: params,
      }),

      invalidatesTags: ['User'],
    }),
    logout: build.mutation<ApiResponse<LoginResponse>, void>({
      query() {
        return {
          url: '/user/logout',
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
