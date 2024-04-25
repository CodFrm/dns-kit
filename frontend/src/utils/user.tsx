import { useLogoutMutation } from '@/api/user';
import { setUserLogout } from '@/store/global';

export function checkLogin() {
  return localStorage.getItem('userStatus') === 'login';
}

export function userLogout() {
  setUserLogout();
  localStorage.removeItem('userStatus');
  localStorage.removeItem("refreshToken");
  if (window.location.pathname.replace(/\//g, '') !== 'login') {
    window.location.pathname = '/login';
  }
}
