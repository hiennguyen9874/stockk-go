import type { FC } from 'react';
import { memo } from 'react';
import cx from 'classnames';

interface CharTabItemProps {
  name: string;
  isActive?: boolean;
  onClick: () => void;
}

const CharTabItem: FC<CharTabItemProps> = ({ name, isActive, onClick }) => {
  return (
    <div
      className={cx(
        'h-full w-28 px-2 text-white flex justify-start items-center cursor-pointer box-border',
        { 'border-b-2 border-b-gray-400': isActive }
      )}
      onClick={onClick}
      aria-hidden="true"
    >
      <div className="h-full w-full flex flex-row justify-between items-center">
        <div className="bg-white text-xs">
          <span>
            <svg
              data-icon="waterfall-chart"
              width="16"
              height="16"
              viewBox="0 0 16 16"
            >
              <path
                d="M8 7c.55 0 1-.45 1-1s-.45-1-1-1-1 .45-1 1 .45 1 1 1zm-4 4h1c.55 0 1-.45 1-1V8c0-.55-.45-1-1-1H4c-.55 0-1 .45-1 1v2c0 .55.45 1 1 1zm7-6c.55 0 1-.45 1-1V3c0-.55-.45-1-1-1s-1 .45-1 1v1c0 .55.45 1 1 1zm4-3h-1c-.55 0-1 .45-1 1v7c0 .55.45 1 1 1h1c.55 0 1-.45 1-1V3c0-.55-.45-1-1-1zm0 10H2V3c0-.55-.45-1-1-1s-1 .45-1 1v10c0 .55.45 1 1 1h14c.55 0 1-.45 1-1s-.45-1-1-1z"
                fillRule="evenodd"
              />
            </svg>
          </span>
        </div>
        <div className="text-sm ml-2 mr-auto">{name}</div>
        <div className="opacity-50 hover:opacity-100">
          <span>
            <svg width="15" height="15" viewBox="0 0 16 16" fill="currentColor">
              <desc>cross</desc>
              <path
                d="M9.41 8l3.29-3.29c.19-.18.3-.43.3-.71a1.003 1.003 0 00-1.71-.71L8 6.59l-3.29-3.3a1.003 1.003 0 00-1.42 1.42L6.59 8 3.3 11.29c-.19.18-.3.43-.3.71a1.003 1.003 0 001.71.71L8 9.41l3.29 3.29c.18.19.43.3.71.3a1.003 1.003 0 00.71-1.71L9.41 8z"
                fillRule="evenodd"
              />
            </svg>
          </span>
        </div>
      </div>
    </div>
  );
};

CharTabItem.defaultProps = {
  isActive: false,
};

export default memo(CharTabItem);
