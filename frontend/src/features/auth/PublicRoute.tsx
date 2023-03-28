import type { FC } from 'react';
import { ReactNode } from 'react';
import { Navigate } from 'react-router-dom';

import { useAppSelector } from 'app/hooks';

import { selectIsAuthenticated } from './authSlice';

interface PublicRouteProps {
  to: string;
  children: ReactNode;
}

const PublicRoute: FC<PublicRouteProps> = ({ to, children }) => {
  const isAuthenticated = useAppSelector(selectIsAuthenticated);

  return (
    <>
      {isAuthenticated && <Navigate to={to} />}
      {!isAuthenticated && children}
    </>
  );
};

export default PublicRoute;
