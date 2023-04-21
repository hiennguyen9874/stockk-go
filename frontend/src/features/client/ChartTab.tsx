import type { FC } from 'react';
import { Dispatch, SetStateAction, useEffect } from 'react';
import cx from 'classnames';

import { CharTabItem } from 'components/common/TabItem';
import { AddIcon } from 'components/common/Icon';
import {
  useGetClientsQuery,
  useDeleteClientMutation,
  useCreateClientMutation,
  useUpdateClientMutation,
} from 'app/services/client';

interface ChartTabProps {
  chartIdx: number;
  setChartIdx: Dispatch<SetStateAction<number>>;
  symbol: string;
  setSymbol: Dispatch<SetStateAction<string>>;
}

const ChartTab: FC<ChartTabProps> = ({
  chartIdx,
  setChartIdx,
  setSymbol,
  symbol,
}) => {
  const { data: clients } = useGetClientsQuery();
  const [createClient] = useCreateClientMutation();
  const [deleteClient] = useDeleteClientMutation();
  const [updateClient] = useUpdateClientMutation();

  useEffect(() => {
    if (clients?.data.length) {
      if (clients.data.find((x) => x.id === chartIdx) === undefined) {
        const client = [...clients.data].sort((a, b) =>
          a.id > b.id ? 1 : -1
        )[0];

        setChartIdx(client.id);
        setSymbol(client.current_ticker);
      }
    }
  }, [chartIdx, clients, setChartIdx, setSymbol]);

  useEffect(() => {
    if (clients !== undefined && clients.data.length === 0) {
      (async () => {
        await createClient({
          current_ticker: 'TCB',
          current_resolution: 'D',
        }).unwrap();
      })();
    }
  }, [clients, createClient]);

  useEffect(() => {
    if (clients !== undefined && clients.data.length > 0) {
      const client = clients.data.find((x) => x.id === chartIdx);

      if (client && symbol !== client.current_ticker) {
        updateClient({
          ...client,
          current_ticker: symbol,
        });
      }
    }
  }, [chartIdx, clients, symbol, updateClient]);

  return (
    <div className="w-full h-8 flex flex-row justify-between rounded-sm bg-slate-700">
      <div className="flex flex-row">
        {clients &&
          [...clients.data]
            .sort((a, b) => (a.id > b.id ? 1 : -1))
            .map((client) => (
              <CharTabItem
                key={client.id}
                name={`${client.current_ticker} (${client.current_resolution})`}
                onClick={() => {
                  setChartIdx(client.id);
                  setSymbol(client.current_ticker);
                }}
                isActive={chartIdx === client.id}
                onDelete={async () => {
                  await deleteClient(client.id).unwrap();
                }}
              />
            ))}
      </div>

      <div>
        <button
          type="button"
          className={cx(
            'group flex flex-row items-center rounded-sm px-2 py-2 text-sm font-sans font-normal text-gray-100',
            'rounded-md',
            'hover:bg-slate-600'
          )}
          // eslint-disable-next-line @typescript-eslint/no-misused-promises
          onClick={async (e) => {
            e.preventDefault();
            e.stopPropagation();

            await createClient({
              current_ticker: 'TCB',
              current_resolution: 'D',
            }).unwrap();
          }}
        >
          <AddIcon />
        </button>
      </div>
    </div>
  );
};

export default ChartTab;
