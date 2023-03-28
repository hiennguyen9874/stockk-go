import type { FC } from 'react';
import { memo } from 'react';

const RemoveIcon: FC = () => (
  <svg
    data-icon="cross"
    width="16"
    height="16"
    viewBox="0 0 16 16"
    fill="currentColor"
  >
    <desc>cross</desc>
    <path
      d="M9.41 8l3.29-3.29c.19-.18.3-.43.3-.71a1.003 1.003 0 00-1.71-.71L8 6.59l-3.29-3.3a1.003 1.003 0 00-1.42 1.42L6.59 8 3.3 11.29c-.19.18-.3.43-.3.71a1.003 1.003 0 001.71.71L8 9.41l3.29 3.29c.18.19.43.3.71.3a1.003 1.003 0 00.71-1.71L9.41 8z"
      fillRule="evenodd"
    />
  </svg>
);

export default memo(RemoveIcon);
