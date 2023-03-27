import { createBrowserRouter } from 'react-router-dom';

import Home from './Home';
import SignUp from './SignUp';
import Chart from './Chart';

const router = createBrowserRouter([
  {
    path: '/',
    element: <Home />,
  },
  {
    path: '/signup',
    element: <SignUp />,
  },
  {
    path: '/chart',
    element: <Chart />,
  },
]);

export default router;
