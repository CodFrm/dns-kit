import auth, { AuthParams } from '@/utils/authentication';
import { useEffect, useMemo, useState } from 'react';
import { Route, Routes as ReactRoutes, Navigate } from 'react-router-dom';
import Login from './pages/login';
import PageLayout, { getFlattenRoutes } from './layout';
import { useSelector } from 'react-redux';
import React from 'react';
import Exception403 from './pages/exception/403';
import {
  GlobalState,
  updateUserInfo,
  userLoading,
  userLogout,
} from './store/global';
import { checkLogin } from './utils/user';
import { useLazyCurrentUserQuery } from './services/user.service';
import { useAppDispatch } from './store/hooks';

export type IRoute = AuthParams & {
  name: string;
  key: string;
  // 当前页是否展示面包屑
  breadcrumb?: boolean;
  children?: IRoute[];
  // 当前路由是否渲染菜单项，为 true 的话不会在菜单中显示，但可通过路由地址访问。
  ignore?: boolean;
};

export const routes: IRoute[] = [
  {
    name: 'menu.dashboard',
    key: 'dashboard',
    children: [
      {
        name: 'menu.dashboard.workplace',
        key: 'dashboard/workplace',
      },
    ],
  },
  {
    name: 'CDN管理',
    key: 'example',
  },
  {
    name: '域名解析',
    key: 'example',
    children: [
      {
        name: '我的域名',
        key: 'example',
      },
      {
        name: '域名监控',
        key: 'example',
      },
    ],
  },
  {
    name: '证书管理',
    key: 'example',
    children: [
      {
        name: '证书签发',
        key: 'example',
      },
      {
        name: '证书托管',
        key: 'example',
      },
    ],
  },
  {
    name: '厂商管理',
    key: 'provider',
  },
  {
    name: '系统设置',
    key: 'example',
  },
];

export const getName = (path: string, routes) => {
  return routes.find((item) => {
    const itemPath = `/${item.key}`;
    if (path === itemPath) {
      return item.name;
    } else if (item.children) {
      return getName(path, item.children);
    }
  });
};

export const generatePermission = (role: string) => {
  const actions = role === 'admin' ? ['*'] : ['read'];
  const result = {};
  routes.forEach((item) => {
    if (item.children) {
      item.children.forEach((child) => {
        result[child.name] = actions;
      });
    }
  });
  return result;
};

const useRoute = (userPermission): [IRoute[], string] => {
  const filterRoute = (routes: IRoute[], arr = []): IRoute[] => {
    if (!routes.length) {
      return [];
    }
    for (const route of routes) {
      const { requiredPermissions, oneOfPerm } = route;
      let visible = true;
      if (requiredPermissions) {
        visible = auth({ requiredPermissions, oneOfPerm }, userPermission);
      }

      if (!visible) {
        continue;
      }
      if (route.children && route.children.length) {
        const newRoute = { ...route, children: [] };
        filterRoute(route.children, newRoute.children);
        if (newRoute.children.length) {
          arr.push(newRoute);
        }
      } else {
        arr.push({ ...route });
      }
    }

    return arr;
  };

  const [permissionRoute, setPermissionRoute] = useState(routes);

  useEffect(() => {
    const newRoutes = filterRoute(routes);
    setPermissionRoute(newRoutes);
  }, [JSON.stringify(userPermission)]);

  const defaultRoute = useMemo(() => {
    const first = permissionRoute[0];
    if (first) {
      const firstRoute = first?.children?.[0]?.key || first.key;
      return firstRoute;
    }
    return '';
  }, [permissionRoute]);

  return [permissionRoute, defaultRoute];
};

export const Routes = () => {
  const [currentUser, { data }] = useLazyCurrentUserQuery();
  const { settings, userInfo } = useSelector((state: GlobalState) => state);
  const [routes, defaultRoute] = useRoute(userInfo?.permissions);
  const flattenRoutes = useMemo(() => getFlattenRoutes(routes) || [], [routes]);
  const dispatch = useAppDispatch();

  function fetchUserInfo() {
    dispatch(userLoading());
    currentUser()
      .unwrap()
      .then((res) => {
        dispatch(updateUserInfo(res.data));
      })
      .catch((e) => {
        if (e.status == 401) {
          dispatch(userLogout());
        }
      });
  }

  useEffect(() => {
    if (checkLogin()) {
      fetchUserInfo();
    } else if (window.location.pathname.replace(/\//g, '') !== 'login') {
      window.location.pathname = '/login';
    }
  }, []);

  return (
    <React.Fragment>
      <ReactRoutes>
        <Route path="/login" element={<Login />} />
        <Route path="/" element={<PageLayout />}>
          {flattenRoutes.map((route, index) => {
            return (
              <Route
                key={index}
                path={`${route.key}`}
                element={<route.component />}
              />
            );
          })}
          <Route index element={<Navigate to={`${defaultRoute}`} />} />
          <Route path="*" element={<Exception403 />} />
        </Route>
      </ReactRoutes>
    </React.Fragment>
  );
};

export default useRoute;
