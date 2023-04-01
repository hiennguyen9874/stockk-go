import type { FC } from 'react';
import { ReactNode } from 'react';
import { Navigate } from 'react-router-dom';

import { useAppSelector } from 'app/hooks';

import { selectIsAuthenticated } from './authSlice';

interface ProtectedRouteProps {
  to: string;
  children: ReactNode;
}

const ProtectedRoute: FC<ProtectedRouteProps> = ({ to, children }) => {
  const isAuthenticated = useAppSelector(selectIsAuthenticated);

  return (
    <>
      {!isAuthenticated && <Navigate to={to} />}
      {isAuthenticated && children}
    </>
  );
};

export default ProtectedRoute;
