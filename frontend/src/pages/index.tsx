import { createBrowserRouter } from 'react-router-dom';

import Home from './Home';
import Chart from './Chart';

const router = createBrowserRouter([
  {
    path: '/',
    element: <Home />,
  },
  {
    path: '/chart',
    element: <Chart />,
  },
]);

export default router;
