export function checkLogin() {
  return localStorage.getItem('userStatus') === 'login';
}
