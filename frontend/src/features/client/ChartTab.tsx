import type { FC } from 'react';
import { Dispatch, SetStateAction, useEffect } from 'react';
import cx from 'classnames';

import { CharTabItem } from 'components/common/TabItem';
import { AddIcon } from 'components/common/Icon';
import { useGetClientsQuery } from 'app/services/client';

interface ChartTabProps {
  chartIdx: number;
  setChartIdx: Dispatch<SetStateAction<number>>;
  setSymbol: Dispatch<SetStateAction<string>>;
}

const ChartTab: FC<ChartTabProps> = ({ chartIdx, setChartIdx, setSymbol }) => {
  const { data: clients } = useGetClientsQuery();

  useEffect(() => {
    if (clients) {
      setChartIdx(clients.data[0].id);
      setSymbol(clients.data[0].current_ticker);
    }
  }, [clients, setChartIdx, setSymbol]);

  return (
    <div className="w-full h-8 flex flex-row justify-between rounded-sm bg-slate-700">
      <div className="flex flex-row">
        {clients &&
          clients.data.map((client) => (
            <CharTabItem
              key={client.id}
              name={`${client.current_ticker} (${client.current_resolution})`}
              onClick={() => {
                setChartIdx(client.id);
                setSymbol(client.current_ticker);
              }}
              isActive={chartIdx === client.id}
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
        >
          <AddIcon />
        </button>
      </div>
    </div>
  );
};

export default ChartTab;
