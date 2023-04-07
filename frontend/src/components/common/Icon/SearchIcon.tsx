import type { FC } from 'react';
import { memo } from 'react';

interface SearchIconProps {
  className?: string;
}

const SearchIcon: FC<SearchIconProps> = ({ className }) => (
  <svg
    data-icon="search"
    width="16"
    height="16"
    viewBox="0 0 16 16"
    fill="currentColor"
    className={className}
  >
    <path
      d="M15.55 13.43l-2.67-2.68a6.94 6.94 0 001.11-3.76c0-3.87-3.13-7-7-7s-7 3.13-7 7 3.13 7 7 7c1.39 0 2.68-.42 3.76-1.11l2.68 2.67a1.498 1.498 0 102.12-2.12zm-8.56-1.44c-2.76 0-5-2.24-5-5s2.24-5 5-5 5 2.24 5 5-2.24 5-5 5z"
      fillRule="evenodd"
    />
  </svg>
);

SearchIcon.defaultProps = {
  className: undefined,
};

export default memo(SearchIcon);
