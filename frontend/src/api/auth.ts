import { AxiosResponse, AxiosError, AxiosRequestConfig } from 'axios';
import { Store } from 'redux';

import { logout as logoutAuth } from 'store/reducers/auth';

import authHeader, { axiosInstance, User } from './auth-header';

const API_URL = '/api/';

export const register: (
  username: string,
  email: string,
  password: string
) => Promise<AxiosResponse> = (username, email, password) =>
  axiosInstance.post(`${API_URL}user/create/`, {
    username,
    email,
    password,
  });

export const login: (
  username: string,
  password: string
) => Promise<AxiosResponse> = (username, password) =>
  axiosInstance.post(`${API_URL}token/obtain/`, {
    username,
    password,
  });

const removeItemInLocalStorage: (item: string) => void = (item) => {
  const userJson = localStorage.getItem(item);
  if (userJson) {
    localStorage.removeItem(item);
  }
};

export const logout: () => Promise<AxiosResponse> | null = () => {
  const userJson = localStorage.getItem('user');
  if (!userJson) {
    return null;
  }
  // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
  const user: User = JSON.parse(userJson);

  return axiosInstance
    .post(
      `${API_URL}token/blacklist/`,
      {
        refresh_token: user.refresh,
      },
      {
        headers: authHeader(),
      }
    )
    .finally(() => {
      removeItemInLocalStorage('user');
    });
};

export const changePassword: (
  oldPassword: string,
  newPassword: string
) => Promise<AxiosResponse> = (oldPassword, newPassword) =>
  axiosInstance.put(
    `${API_URL}user/changepassword/`,
    {
      old_password: oldPassword,
      new_password: newPassword,
    },
    {
      headers: authHeader(),
    }
  );

// export const interceptorJWT = (store: Store): void => {
//   axiosInstance.interceptors.response.use(
//     (response: AxiosResponse) => response,
//     (error: AxiosError) => {
//       const { dispatch } = store;
//       const originalRequest: AxiosRequestConfig = error.config;

//       if (!error.response) {
//         // Network error => Reload to home page!
//         removeItemInLocalStorage('user');
//         dispatch(logoutAuth());
//         return Promise.reject(error);
//       }

//       if (
//         error.response.status === 401 &&
//         originalRequest.url === `${API_URL}token/refresh/`
//       ) {
//         // Pre-request is refresh token => Reload to home page!
//         removeItemInLocalStorage('user');
//         return Promise.reject(error);
//       }

//       if (
//         error.response.data &&
//         // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
//         error.response.data.code &&
//         // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
//         error.response.data.code === 'user_inactive'
//       ) {
//         // Account is inactivate => Reload to home page!
//         removeItemInLocalStorage('user');
//         return Promise.reject(error);
//       }

//       if (
//         error.response.data &&
//         // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
//         error.response.data.code &&
//         // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
//         error.response.data.code === 'token_not_valid'
//       ) {
//         const userJson = localStorage.getItem('user');

//         if (!userJson) {
//           // Token not valid and refresh token not in local storage => Reload to home page!
//           removeItemInLocalStorage('user');
//           return Promise.reject(error);
//         }

//         // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
//         const user: User = JSON.parse(userJson);

//         // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
//         const tokenParts = JSON.parse(atob(user.refresh.split('.')[1]));

//         const now = Math.ceil(Date.now() / 1000);

//         // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
//         if (tokenParts.exp > now) {
//           // Get new access from refresh token => return new response
//           return axiosInstance
//             .post(`${API_URL}token/refresh/`, {
//               refresh: user.refresh,
//             })
//             .then((response: AxiosResponse) => {
//               // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
//               const { data } = response;
//               localStorage.setItem(
//                 'user',
//                 JSON.stringify({
//                   // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-member-access
//                   access: data.access,
//                   // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-member-access
//                   refresh: data.refresh,
//                   username: user.username,
//                   isStaff: user.isStaff,
//                 })
//               );
//               originalRequest.headers = authHeader();
//               return axiosInstance(originalRequest);
//             });
//         }

//         // Else remove token in local storage => Reload to home page!
//         removeItemInLocalStorage('user');
//         return Promise.reject(error);
//       }

//       if (
//         error.response.data &&
//         // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
//         error.response.data.code === 'user_not_found'
//       ) {
//         // User not found in database => Reload to home page!
//         removeItemInLocalStorage('user');
//       }

//       return Promise.reject(error);
//     }
//   );
// };
