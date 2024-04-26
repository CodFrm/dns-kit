import { PayloadAction } from '@reduxjs/toolkit';
import defaultSettings from '../settings.json';
import { createAppSlice } from './hooks';
import { CurrentUserResponse } from '@/services/user.service';

export interface UserInfo {
  username?: string;
  avatar?: string;
  job?: string;
  organization?: string;
  location?: string;
  email?: string;
  permissions: Record<string, string[]>;
}

export interface GlobalState {
  settings?: typeof defaultSettings;
  userInfo?: UserInfo;
  userLoading?: boolean;
}

const initialState: GlobalState = {
  settings: defaultSettings,
  userInfo: {
    permissions: {},
  },
};

export const globalSlice = createAppSlice({
  name: 'global',
  initialState,
  reducers: (create) => ({
    updateSetting: create.reducer(
      (state, action: PayloadAction<typeof defaultSettings>) => {
        state.settings = action.payload;
      },
    ),
    updateUserInfo: create.reducer(
      (state, action: PayloadAction<CurrentUserResponse>) => {
        state.userInfo = {
          username: action.payload.username,
          avatar:
            'https://gw.alipayobjects.com/zos/rmsportal/BiazfanxmamNRoxxVxka.png',
          permissions: {},
        };
        state.userLoading = false;
      },
    ),
    userLogout: create.reducer((state, action) => {
      state.userInfo = undefined;
      localStorage.removeItem('userStatus');
      localStorage.removeItem('refreshToken');
      if (window.location.pathname.replace(/\//g, '') !== 'login') {
        window.location.pathname = '/login';
      }
    }),
    userLoading: create.reducer((state) => {
      state.userLoading = true;
    }),
  }),
  selectors: {
    selectGlobal: (global) => global,
    selectSetting: (global) => global.settings,
    selectUserInfo: (global) => global.userInfo,
  },
});

export const { userLoading, updateUserInfo, userLogout } = globalSlice.actions;

export const { selectGlobal, selectSetting, selectUserInfo } =
  globalSlice.selectors;
