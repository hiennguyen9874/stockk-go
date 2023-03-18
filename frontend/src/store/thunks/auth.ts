/* eslint-disable @typescript-eslint/no-unsafe-argument */
/* eslint-disable @typescript-eslint/no-unsafe-call */
/* eslint-disable @typescript-eslint/no-unsafe-member-access */
/* eslint-disable @typescript-eslint/no-unsafe-assignment */
import { RootThunkType } from 'store';
import { User } from 'api/auth-header';
import * as AuthService from 'api/auth';
import {
  registerSuccess,
  registerFailure,
  loginSuccess,
  loginFailure,
  logout as logoutAuth,
} from 'store/reducers/auth';
import { setMessage } from 'store/reducers/message';

export const register =
  (username: string, email: string, password: string): RootThunkType =>
  async (dispatch) => {
    try {
      const response = await AuthService.register(username, email, password);
      dispatch(registerSuccess());
      dispatch(setMessage(response.data.message));
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (error: any) {
      const message =
        (error.response &&
          error.response.data &&
          error.response.data.message) ||
        error.message ||
        error.toString();

      dispatch(registerFailure());
      dispatch(setMessage(message));
    }
  };

export const login =
  (
    username: string,
    password: string,
    callBack?: (isSuccess: boolean) => void
  ): RootThunkType =>
  async (dispatch) => {
    try {
      const response = await AuthService.login(username, password);

      if (response.data && response.data.access && response.data.refresh) {
        const dataUser: User = {
          access: response.data.access,
          refresh: response.data.refresh,
          username: response.data.username,
          isStaff: response.data.is_staff,
        };

        if (dataUser && dataUser.access) {
          localStorage.setItem('user', JSON.stringify(dataUser));
        }

        if (callBack) {
          callBack(true);
        }

        dispatch(loginSuccess(dataUser));
      } else {
        const message = response && response.data && response.data.messages;

        if (callBack) {
          callBack(false);
        }

        dispatch(loginFailure());

        dispatch(setMessage(message));
      }
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (error: any) {
      const message =
        (error.response &&
          error.response.data &&
          error.response.data.message) ||
        error.message ||
        error.toString();

      if (callBack) {
        callBack(false);
      }

      dispatch(loginFailure());
      dispatch(setMessage(message));
    }
  };

export const logout = (): RootThunkType => async (dispatch) => {
  try {
    const logoutFunc = AuthService.logout();
    if (logoutFunc) {
      await logoutFunc;
    }
    dispatch(logoutAuth());
  } catch (error) {
    dispatch(logoutAuth());
  }
};
