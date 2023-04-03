import type { FC } from 'react';

import { WatchListCard } from 'components/common/Card';
import { useGetStockSnapshotQuery } from 'app/services/stocksnapshot';
import { useGetTickerQuery } from 'app/services/ticker';

interface WatchListItemProps {
  symbol: string;
  isLight: boolean;
  onSet: () => void;
}

const WatchListItem: FC<WatchListItemProps> = ({ symbol, isLight, onSet }) => {
  const { data: stockSnapshot } = useGetStockSnapshotQuery(symbol);
  const { data: ticker } = useGetTickerQuery(symbol);

  return (
    // eslint-disable-next-line react/jsx-no-useless-fragment
    <>
      {stockSnapshot && ticker && (
        <WatchListCard
          key={symbol}
          symbol={symbol}
          price={stockSnapshot.data.price}
          description={ticker.data.full_name}
          changePrice={stockSnapshot.data.price - stockSnapshot.data.ref_price}
          changePercent={
            ((stockSnapshot.data.price - stockSnapshot.data.ref_price) * 100) /
            stockSnapshot.data.ref_price
          }
          isLight={isLight}
          onClick={() => onSet()}
          status={(() => {
            if (stockSnapshot.data.price === stockSnapshot.data.ceil_price)
              return 'ceil';

            if (stockSnapshot.data.price === stockSnapshot.data.floor_price)
              return 'floor';

            if (stockSnapshot.data.price < stockSnapshot.data.ref_price)
              return 'decrease';

            if (stockSnapshot.data.price > stockSnapshot.data.ref_price)
              return 'increase';

            return 'reference';
          })()}
        />
      )}
    </>
  );
};

export default WatchListItem;
