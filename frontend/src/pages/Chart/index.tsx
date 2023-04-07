import type { FC } from 'react';
import { useState } from 'react';

import { API_DATAFEED_URL, API_STORAGE_URL } from 'configs/api-server';
import TVChartContainer from 'components/common/TVChartContainer';
import { SliderTabItem } from 'components/common/TabItem';
import WatchLists from 'features/watchlists/WatchLists';
import ChartTab from 'features/client/ChartTab';

const Chart: FC = () => {
  const [symbol, setSymbol] = useState('TCB');
  const [chartIdx, setChartIdx] = useState<number>(0);

  console.log(symbol);

  return (
    <div className="w-screen h-screen">
      <div className="h-full flex py-1.5 flex-row bg-slate-900">
        <div className="w-full h-full p-0 rounded-md flex flex-col justify-center items-center">
          <ChartTab
            chartIdx={chartIdx}
            setChartIdx={setChartIdx}
            setSymbol={setSymbol}
          />

          <div className="w-full h-full">
            <TVChartContainer
              // datafeed={Datafeed}
              defaultSymbol="TCB"
              symbol={symbol}
              datafeedUrl={API_DATAFEED_URL}
              chartsStorageUrl={API_STORAGE_URL}
              clientId={`${chartIdx}`}
            />
          </div>
        </div>

        <div className="w-[21rem] shrink-0 grow-0 ml-1 bg-slate-800 border-solid border-x-2 rounded-md border-slate-900">
          <div className="h-full w-full flex flex-row">
            <div className="h-full w-full grow truncate">
              <WatchLists setSymbol={(newSymbol) => setSymbol(newSymbol)} />
            </div>

            <div className="h-full w-6 bg-slate-900">
              <div className="flex flex-col">
                <SliderTabItem name="Watch list" isActive />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Chart;
