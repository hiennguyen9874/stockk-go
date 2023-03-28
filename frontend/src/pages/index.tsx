import { createBrowserRouter } from 'react-router-dom';

import ProtectedRoute from 'features/auth/ProtectedRoute';
import PublicRoute from 'features/auth/PublicRoute';

import Home from './Home';
import SignUp from './SignUp';
import Chart from './Chart';

const router = createBrowserRouter([
  {
    path: '/',
    element: (
      <PublicRoute to="/chart">
        <Home />
      </PublicRoute>
    ),
  },
  {
    path: '/signup',
    element: (
      <PublicRoute to="/chart">
        <SignUp />
      </PublicRoute>
    ),
  },
  {
    path: '/chart',
    element: (
      <ProtectedRoute to="/">
        <Chart />
      </ProtectedRoute>
    ),
  },
]);

export default router;
