import type { FC } from 'react';
import { memo } from 'react';

const EditIcon: FC = () => (
  <svg
    data-icon="edit"
    width="16"
    height="16"
    viewBox="0 0 16 16"
    fill="currentColor"
  >
    <desc>edit</desc>
    <path
      d="M3.25 10.26l2.47 2.47 6.69-6.69-2.46-2.48-6.7 6.7zM.99 14.99l3.86-1.39-2.46-2.44-1.4 3.83zm12.25-14c-.48 0-.92.2-1.24.51l-1.44 1.44 2.47 2.47 1.44-1.44c.32-.32.51-.75.51-1.24.01-.95-.77-1.74-1.74-1.74z"
      fillRule="evenodd"
    />
  </svg>
);

export default memo(EditIcon);
