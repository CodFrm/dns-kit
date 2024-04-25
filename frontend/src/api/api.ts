import { Message } from '@arco-design/web-react';
import { PayloadAction, isRejectedWithValue } from '@reduxjs/toolkit';
import { Action, Dispatch, Middleware, MiddlewareAPI } from 'redux';

export interface ApiResponse<T = any> {
  code: number;
  msg: string;
  data: T;
}

export type ListApiResponse<T = any> = ApiResponse<{
  list: T[];
  total: number;
}>;

export const rtkQueryErrorHandler: Middleware =
  (api: MiddlewareAPI) =>
  (next: Dispatch<Action>) =>
  (action: PayloadAction<{ status: number; data: any; error: string }>) => {
    if (isRejectedWithValue(action)) {
      console.log('pl', action.payload);
      if (typeof action.payload.status == 'string') {
        Message.error('请求失败: ' + action.payload.error);
      } else if (action.payload.status >= 400) {
        if (action.payload.status >= 500) {
          Message.error('系统错误');
        } else {
          Message.error(action.payload.data?.msg);
        }
      }
    }
    return next(action);
  };
