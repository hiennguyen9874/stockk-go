import type { FC } from 'react';
import { memo } from 'react';
import cx from 'classnames';

interface SliderTabItemProps {
  name: string;
  isActive?: boolean;
}

const SliderTabItem: FC<SliderTabItemProps> = ({ name, isActive }) => {
  return (
    <div
      className={cx(
        'w-full h-20 pt-2 mb-1 rounded-tr-md rounded-br-md cursor-pointer bg-cur',
        'text-sm font-sans font-normal',
        'text-gray-100',
        {
          'bg-slate-800': isActive,
        }
      )}
    >
      <div className="transform rotate-90 text-center">
        <div className="whitespace-nowrap">{name}</div>
      </div>
    </div>
  );
};

SliderTabItem.defaultProps = {
  isActive: false,
};

export default memo(SliderTabItem);
