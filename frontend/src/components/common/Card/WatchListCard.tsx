import type { FC } from 'react';
import { memo } from 'react';
import cx from 'classnames';

interface WatchListCardProps {
  symbol: string;
  price: number;
  description: string;
  changePrice: number;
  changePercent: number;
  isLight?: boolean;
  onClick: () => void;
  status: 'ceil' | 'floor' | 'increase' | 'decrease' | 'reference';
  className?: string;
}

const WatchListCard: FC<WatchListCardProps> = ({
  symbol,
  price,
  description,
  changePrice,
  changePercent,
  isLight,
  onClick,
  status,
  className,
}) => {
  return (
    <div
      className={cx(
        'w-full h-14 flex flex-col justify-between px-2 pt-0.5 pb-1',
        'rounded-sm',
        'cursor-pointer',
        'text-sm font-sans font-normal',
        'text-gray-100',
        {
          'bg-[#243143] hover:bg-slate-600': isLight,
          'bg-[#1e293b] hover:bg-slate-600': !isLight,
        },
        className
      )}
      role="button"
      aria-hidden="true"
      onClick={() => onClick()}
    >
      <div className="flex flex-row justify-between">
        <div className="text-lg font-bold text-blue-600">{symbol}</div>
        <div className="text-base font-bold">{price}</div>
      </div>

      <div className="flex flex-row justify-between">
        <div className="truncate">{description}</div>
        <div
          className={cx('ml-4 text-sm font-bold', {
            'text-[#ff3747]': status === 'decrease',
            'text-[#00f4b0]': status === 'increase',
            'text-[#fbac20]': status === 'reference',
            'text-[#e683ff]': status === 'ceil',
            'text-[#64baff]': status === 'floor',
          })}
        >{`${changePrice.toFixed(2)}/${changePercent.toFixed(2)}%`}</div>
      </div>
    </div>
  );
};

WatchListCard.defaultProps = {
  isLight: false,
  className: '',
};

export default memo(WatchListCard);
