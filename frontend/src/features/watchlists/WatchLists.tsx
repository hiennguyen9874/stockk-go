import type { FC } from 'react';
import { useState, useEffect } from 'react';

import { WatchListsDropdown } from 'components/common/Dropdown';
import { Item } from 'components/common/Dropdown/WatchListsDropdown';
import { WatchListCard } from 'components/common/Card';
import { useGetWatchListsQuery } from 'app/services/watchlists';

interface WatchListProps {
  setSymbol: (symbol: string) => void;
}

const WatchLists: FC<WatchListProps> = ({ setSymbol }) => {
  const { data: dataWatchLists } = useGetWatchListsQuery();

  const [watchLists, setWatchLists] = useState<
    {
      id: number;
      name: string;
      tickers: string[];
    }[]
  >([]);

  const [currentWatchList, setCurrentWatchList] = useState<Item | null>(
    watchLists.length === 0 ? null : watchLists[0]
  );

  useEffect(() => {
    if (dataWatchLists === undefined || dataWatchLists.data.length === 0) {
      setWatchLists([]);
    } else {
      const newWatchLists = dataWatchLists.data.map(
        ({ id, name, tickers }) => ({
          id,
          name,
          tickers,
        })
      );

      setWatchLists(newWatchLists);
      if (newWatchLists.length > 0) {
        setCurrentWatchList(newWatchLists[0]);
      }
    }
  }, [dataWatchLists]);

  return (
    <div className="h-full flex flex-col">
      <div className="h-10">
        <WatchListsDropdown
          currentItem={currentWatchList}
          items={watchLists}
          onClick={(item) => setCurrentWatchList(item)}
          onEdit={(editId, editValue) =>
            setWatchLists(
              watchLists.map(({ id, name, tickers }) => ({
                id,
                name: id === editId ? editValue : name,
                tickers,
              }))
            )
          }
          onDelete={(deleteId) => {
            const newWatchLists = watchLists.filter(
              ({ id }) => id !== deleteId
            );

            if (newWatchLists.length === 0) setCurrentWatchList(null);

            setWatchLists(newWatchLists);

            // TODO: Call api to remove watchList
          }}
          onAdd={() => {
            const newWatchLists = [
              ...watchLists,
              {
                id: Math.max(...watchLists.map(({ id }) => id)) + 1,
                name: 'Danh mục mới',
                tickers: [],
              },
            ];

            // TODO: Call api to remove watchLists

            if (watchLists.length === 0) setCurrentWatchList(newWatchLists[0]);

            setWatchLists(newWatchLists);

            // TODO: Call api to add watchLists
          }}
        />
      </div>

      <div className="h-auto overflow-auto scroll-smooth">
        <WatchListCard
          symbol="VCI"
          price={31.05}
          description="Chứng khoán bản việt"
          changePrice={-0.35}
          changePercent={-1.11}
          isLight
          onClick={() => setSymbol('VCI')}
        />
        <WatchListCard
          symbol="FTS"
          price={21.65}
          description="Chứng khoán FPT"
          changePrice={-0.35}
          changePercent={-1.59}
          onClick={() => setSymbol('FTS')}
        />
        <WatchListCard
          symbol="HCM"
          price={24.5}
          description="Chứng khoán Thành Phố Hồ Chí Minh, Hồ Chí Minh"
          changePrice={0}
          changePercent={0}
          isLight
          onClick={() => setSymbol('HCM')}
        />
      </div>
    </div>
  );
};

export default WatchLists;
