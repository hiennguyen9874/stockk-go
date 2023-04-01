import type { FC } from 'react';
import { memo } from 'react';

const AddIcon: FC = () => (
  <svg
    data-icon="plus"
    width="16"
    height="16"
    viewBox="0 0 16 16"
    fill="currentColor"
  >
    <desc>plus</desc>
    <path
      d="M13 7H9V3c0-.55-.45-1-1-1s-1 .45-1 1v4H3c-.55 0-1 .45-1 1s.45 1 1 1h4v4c0 .55.45 1 1 1s1-.45 1-1V9h4c.55 0 1-.45 1-1s-.45-1-1-1z"
      fillRule="evenodd"
    />
  </svg>
);

export default memo(AddIcon);
