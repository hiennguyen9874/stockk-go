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
}

const WatchListCard: FC<WatchListCardProps> = ({
  symbol,
  price,
  description,
  changePrice,
  changePercent,
  isLight,
  onClick,
}) => {
  return (
    <div
      className={cx(
        'h-14 flex flex-col justify-between px-2 pt-0.5 pb-1',
        'rounded-sm',
        'cursor-pointer',
        'bg-slate-800',
        'text-sm font-sans font-normal',
        'text-gray-100',
        'hover:bg-slate-600',
        {
          'bg-slate-700 hover:bg-slate-600': isLight,
        }
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
            'text-red-500': changePercent < 0,
            'text-green-500': changePercent > 0,
            'text-yellow-500': changePercent === 0,
          })}
        >{`${changePrice.toFixed(2)}/${changePercent.toFixed(2)}%`}</div>
      </div>
    </div>
  );
};

WatchListCard.defaultProps = {
  isLight: false,
};

export default memo(WatchListCard);
