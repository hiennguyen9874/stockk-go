import { configureStore, AnyAction } from '@reduxjs/toolkit';
import { ThunkAction } from 'redux-thunk';

// import { interceptorJWT } from 'api/auth';

import rootReducer from './reducers';

const store = configureStore({ reducer: rootReducer, devTools: false });

export type RootStateType = ReturnType<typeof store.getState>;

export type RootDispatchType = typeof store.dispatch;

// export type DispatchFunctionType = ThunkDispatch<
//   RootStateType,
//   undefined,
//   AnyAction
// >;

export type RootThunkType<ReturnType = void> = ThunkAction<
  ReturnType,
  RootStateType,
  unknown,
  AnyAction
>;

// interceptorJWT(store);

export default store;
