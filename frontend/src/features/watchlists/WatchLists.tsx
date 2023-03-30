import type { FC } from 'react';
import { useState, useEffect, useMemo } from 'react';

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

  const [currentWatchList, setCurrentWatchList] = useState<number | null>(null);

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
        setCurrentWatchList(newWatchLists[0].id);
      }
    }
  }, [dataWatchLists]);

  const currentItem = useMemo(() => {
    if (currentWatchList === null) return null;

    const findItem = watchLists.find(({ id }) => id === currentWatchList);

    if (findItem === undefined) return null;

    return findItem;
  }, [currentWatchList, watchLists]);

  return (
    <div className="h-full flex flex-col">
      <div className="h-10">
        <WatchListsDropdown
          currentItem={currentItem}
          items={watchLists}
          onClick={(item) => setCurrentWatchList(item.id)}
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

            if (watchLists.length === 0)
              setCurrentWatchList(newWatchLists[0].id);

            setWatchLists(newWatchLists);

            // TODO: Call api to add watchLists
          }}
        />
      </div>

      <div className="h-auto overflow-auto scroll-smooth">
        {/* {currentItem !== null &&
          currentItem.tickers.map((item, idx) => (
            <WatchListCard
              key={item}
              symbol={item}
              price={31.05}
              description="Chứng khoán bản việt"
              changePrice={-0.35}
              changePercent={-1.11}
              isLight={idx % 2 === 0}
              onClick={() => setSymbol(item)}
            />
          ))}
        <WatchListCard
          symbol="FTS"
          price={21.65}
          description="Chứng khoán FPT"
          changePrice={-0.35}
          changePercent={-1.59}
          onClick={() => setSymbol('FTS')}
        /> */}
      </div>
    </div>
  );
};

export default WatchLists;
