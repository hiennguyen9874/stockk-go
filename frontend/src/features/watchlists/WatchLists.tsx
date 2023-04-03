import type { FC } from 'react';
import { useState, useEffect, useMemo } from 'react';

import { WatchListsDropdown } from 'components/common/Dropdown';
import {
  useGetWatchListsQuery,
  useCreateWatchListMutation,
  useDeleteWatchListMutation,
  useUpdateWatchListMutation,
} from 'app/services/watchlists';
import { AddIcon } from 'components/common/Icon';
import WatchListItem from 'features/ticker/WatchListItem';

interface WatchListProps {
  setSymbol: (symbol: string) => void;
}

const WatchLists: FC<WatchListProps> = ({ setSymbol }) => {
  const { data: dataWatchLists, refetch } = useGetWatchListsQuery();
  const [createWatchList] = useCreateWatchListMutation();
  const [deleteWatchList] = useDeleteWatchListMutation();
  const [updateWatchList] = useUpdateWatchListMutation();

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
    <div className="h-full w-full flex flex-col">
      <div className="h-10">
        <WatchListsDropdown
          currentItem={currentItem}
          items={watchLists}
          onClick={(item) => {
            setCurrentWatchList(item.id);
          }}
          onEdit={async (editId, editValue) => {
            let curWatchList: {
              id: number;
              name: string;
              tickers: string[];
            } | null = null;

            setWatchLists(
              watchLists.map(({ id, name, tickers }) => {
                if (id === editId) {
                  curWatchList = {
                    id,
                    name: editValue,
                    tickers,
                  };
                  return curWatchList;
                }

                return {
                  id,
                  name,
                  tickers,
                };
              })
            );

            if (curWatchList !== null)
              await updateWatchList(curWatchList).unwrap();
          }}
          onDelete={async (deleteId) => {
            const newWatchLists = watchLists.filter(
              ({ id }) => id !== deleteId
            );

            if (newWatchLists.length === 0) setCurrentWatchList(null);

            setWatchLists(newWatchLists);

            await deleteWatchList(deleteId).unwrap();

            refetch();
          }}
          onAdd={async () => {
            const newWatchLists = [
              ...watchLists,
              {
                id: Math.max(...watchLists.map(({ id }) => id)) + 1,
                name: 'Danh mục mới',
                tickers: [],
              },
            ];

            if (watchLists.length === 0)
              setCurrentWatchList(newWatchLists[0].id);

            setWatchLists(newWatchLists);

            await createWatchList({
              name: 'Danh mục mới',
            }).unwrap();

            refetch();
          }}
        />
      </div>

      <div className="w-full overflow-y-auto scroll-smooth grow">
        <div className="w-full h-full flex flex-col divide-y divide-white divide-opacity-20">
          <div className="w-full grow">
            {currentItem !== null &&
              currentItem.tickers.map((item, idx) => (
                <WatchListItem
                  key={item}
                  symbol={item}
                  isLight={idx % 2 === 0}
                  onSet={() => setSymbol(item)}
                />
              ))}
          </div>

          <div className="w-full py-1 flex flex-row justify-center">
            <button
              type="button"
              className="group flex flex-row items-center rounded-sm px-2 py-2 text-sm font-sans font-normal text-gray-100"
            >
              <div className="pr-1 pb-[1px]">
                <AddIcon />
              </div>
              Thêm mã mới
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default WatchLists;
