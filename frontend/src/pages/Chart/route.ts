import { lazy } from 'react';

export default {
  path: '/chart',
  extract: false,
  publicPage: true,
  onlyForAdmin: false,
  component: lazy(() => import('./index')),
};
