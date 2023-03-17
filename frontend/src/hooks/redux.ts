import { TypedUseSelectorHook, useDispatch, useSelector } from 'react-redux';

import type { RootStateType, RootDispatchType } from 'store';

// Use throughout your app instead of plain `useDispatch` and `useSelector`
// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types
export const useAppDispatch = () => useDispatch<RootDispatchType>();
export const useAppSelector: TypedUseSelectorHook<RootStateType> = useSelector;
