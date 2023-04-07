import type { FC } from 'react';

import { WatchListCard } from 'components/common/Card';
import {
  useGetTickerQuery,
  useGetTickerSnapshotQuery,
} from 'app/services/ticker';

interface WatchListItemProps {
  symbol: string;
  isLight: boolean;
  onSet: () => void;
  className?: string;
}

const WatchListItem: FC<WatchListItemProps> = ({
  symbol,
  isLight,
  onSet,
  className,
}) => {
  const { data: tickerSnapshot } = useGetTickerSnapshotQuery(symbol, {
    pollingInterval: 10000,
  });
  const { data: ticker } = useGetTickerQuery(symbol);

  return (
    // eslint-disable-next-line react/jsx-no-useless-fragment
    <>
      {tickerSnapshot && ticker && (
        <WatchListCard
          key={symbol}
          symbol={symbol}
          price={tickerSnapshot.data.price}
          description={ticker.data.full_name}
          changePrice={
            tickerSnapshot.data.price - tickerSnapshot.data.ref_price
          }
          changePercent={
            ((tickerSnapshot.data.price - tickerSnapshot.data.ref_price) *
              100) /
            tickerSnapshot.data.ref_price
          }
          isLight={isLight}
          onClick={() => onSet()}
          status={(() => {
            if (tickerSnapshot.data.price === tickerSnapshot.data.ceil_price)
              return 'ceil';

            if (tickerSnapshot.data.price === tickerSnapshot.data.floor_price)
              return 'floor';

            if (tickerSnapshot.data.price < tickerSnapshot.data.ref_price)
              return 'decrease';

            if (tickerSnapshot.data.price > tickerSnapshot.data.ref_price)
              return 'increase';

            return 'reference';
          })()}
          className={className}
        />
      )}
    </>
  );
};

WatchListItem.defaultProps = {
  className: '',
};

export default WatchListItem;
