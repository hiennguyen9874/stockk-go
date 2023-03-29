import type { FC } from 'react';
import { useState } from 'react';

import TVChartContainer from 'components/common/TVChartContainer';
import { CharTabItem, SliderTabItem } from 'components/common/TabItem';
import { API_DATAFEED_URL, API_STORAGE_URL } from 'configs/api-server';
import { SliderDropdown } from 'components/common/Dropdown';
import { Item } from 'components/common/Dropdown/SliderDropdown';
import { WatchListCard } from 'components/common/Card';

const Chart: FC = () => {
  const [symbol, setSymbol] = useState('TCB');
  const [chartIdx, setChartIdx] = useState<number>(0);
  const [currentWatchList, setCurrentWatchList] = useState<Item>({
    id: 0,
    name: 'BDS',
  });

  return (
    <div className="w-screen h-screen">
      <div className="h-full flex p-1.5 flex-row bg-slate-900">
        <div className="w-full h-full p-0 rounded-md flex flex-col justify-center items-center">
          <div className="w-full h-8 flex flex-row rounded-sm bg-slate-700">
            <CharTabItem
              name="VCI"
              onClick={() => setChartIdx(0)}
              isActive={chartIdx === 0}
            />
          </div>
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

        <div className="w-96 ml-1 bg-slate-800 border-solid border-x-2 rounded-md border-slate-900">
          <div className="h-full flex flex-row">
            <div className="h-full mr-auto truncate">
              <div className="h-full flex flex-col">
                <div className="h-10">
                  <SliderDropdown
                    currentItem={currentWatchList}
                    items={[
                      {
                        id: 0,
                        name: 'BDS',
                      },
                      {
                        id: 1,
                        name: 'Ngân hàng',
                      },
                      {
                        id: 2,
                        name: 'Chứng khoán',
                      },
                      {
                        id: 3,
                        name: 'Thép',
                      },
                    ]}
                    onChange={(item) => setCurrentWatchList(item)}
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
            </div>

            <div className="w-6 bg-slate-900">
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
