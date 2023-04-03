import type { FC } from 'react';
import {
  useLayoutEffect,
  useState,
  useEffect,
  useRef,
  memo,
  useMemo,
} from 'react';

import type {
  ChartingLibraryWidgetOptions,
  IChartingLibraryWidget,
  ResolutionString,
} from 'charting_library/charting_library';
import { widget } from 'charting_library';

interface TVChartContainerProps {
  symbol?: ChartingLibraryWidgetOptions['symbol'];
  defaultSymbol?: ChartingLibraryWidgetOptions['symbol'];
  interval?: ChartingLibraryWidgetOptions['interval'];
  // BEWARE: no trailing slash is expected in feed URL
  datafeedUrl?: string;
  libraryPath?: ChartingLibraryWidgetOptions['library_path'];
  chartsStorageUrl?: ChartingLibraryWidgetOptions['charts_storage_url'];
  chartsStorageApiVersion?: ChartingLibraryWidgetOptions['charts_storage_api_version'];
  clientId?: ChartingLibraryWidgetOptions['client_id'];
  userId?: ChartingLibraryWidgetOptions['user_id'];
  fullscreen?: ChartingLibraryWidgetOptions['fullscreen'];
  autosize?: ChartingLibraryWidgetOptions['autosize'];
  studiesOverrides?: ChartingLibraryWidgetOptions['studies_overrides'];
  // container?: ChartingLibraryWidgetOptions["container"];
  // datafeed: ChartingLibraryWidgetOptions['datafeed'];
  timezone?: ChartingLibraryWidgetOptions['timezone'];
}

const TVChartContainer: FC<TVChartContainerProps> = ({
  symbol,
  defaultSymbol,
  datafeedUrl,
  interval,
  libraryPath,
  chartsStorageUrl,
  chartsStorageApiVersion,
  clientId,
  userId,
  fullscreen,
  autosize,
  studiesOverrides,
  // datafeed,
  timezone,
}) => {
  const containerRef = useRef<HTMLDivElement | null>(null);
  const tvWidgetRef = useRef<IChartingLibraryWidget | null>(null);
  const [isReady, setIsReady] = useState(false);

  useEffect(() => {
    if (containerRef.current !== null) {
      const widgetOptions: ChartingLibraryWidgetOptions = {
        symbol: defaultSymbol as string,
        // BEWARE: no trailing slash is expected in feed URL
        // tslint:disable-next-line:no-any
        // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-call, @typescript-eslint/no-unsafe-member-access, @typescript-eslint/no-explicit-any
        datafeed: new (window as any).Datafeeds.UDFCompatibleDatafeed(
          datafeedUrl
        ),
        // datafeed,
        interval: interval || ('D' as ResolutionString),
        container: containerRef.current,
        library_path: libraryPath as string,
        locale: 'vi',
        disabled_features: ['use_localstorage_for_settings'],
        enabled_features: ['study_templates'],
        load_last_chart: true,
        charts_storage_url: chartsStorageUrl,
        charts_storage_api_version: chartsStorageApiVersion,
        client_id: clientId,
        user_id: userId,
        fullscreen,
        autosize,
        studies_overrides: studiesOverrides,
        auto_save_delay: 5,
        theme: 'Dark',
        timezone,
        settings_adapter: {
          setValue: (key: string, value: string): void => {
            localStorage.setItem(key, value);
          },
          removeValue: (key: string): void => {
            localStorage.removeItem(key);
          },
        },
      };

      // eslint-disable-next-line new-cap, no-multi-assign, @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-call, @typescript-eslint/no-unsafe-member-access
      tvWidgetRef.current = new widget(widgetOptions);

      tvWidgetRef.current.onChartReady(() => {
        setIsReady(true);

        tvWidgetRef.current?.subscribe('onAutoSaveNeeded', () => {
          tvWidgetRef.current?.saveChartToServer(
            // eslint-disable-next-line @typescript-eslint/no-empty-function
            () => {},
            // eslint-disable-next-line @typescript-eslint/no-empty-function
            () => {},
            {
              defaultChartName: 'Default layout',
            }
          );
        });

        tvWidgetRef.current?.getSavedCharts((chartRecords) => {
          if (tvWidgetRef.current === null) {
            return;
          }

          if (chartRecords.length === 0) {
            tvWidgetRef.current.activeChart().applyStudyTemplate(
              tvWidgetRef.current.activeChart().createStudyTemplate({
                saveSymbol: true,
              })
            );
          }
        });
      });
    }

    return () => {
      if (tvWidgetRef.current !== null) {
        tvWidgetRef.current.remove();
        tvWidgetRef.current = null;
        setIsReady(false);
      }
    };
  }, [
    autosize,
    chartsStorageApiVersion,
    chartsStorageUrl,
    clientId,
    datafeedUrl,
    fullscreen,
    interval,
    libraryPath,
    studiesOverrides,
    defaultSymbol,
    timezone,
    userId,
  ]);

  useLayoutEffect(() => {
    if (!tvWidgetRef.current) return;

    if (!isReady) return;

    const chart = tvWidgetRef.current?.activeChart();

    if (chart && symbol) {
      // eslint-disable-next-line @typescript-eslint/no-empty-function
      chart.setSymbol(symbol, () => {});
    }
  }, [isReady, symbol]);

  return useMemo(
    () => <div ref={containerRef} className="w-full h-full" />,
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [
      autosize,
      chartsStorageApiVersion,
      chartsStorageUrl,
      clientId,
      datafeedUrl,
      fullscreen,
      interval,
      libraryPath,
      studiesOverrides,
      defaultSymbol,
      timezone,
      userId,
    ]
  );
};

TVChartContainer.defaultProps = {
  symbol: 'AAPL',
  defaultSymbol: 'AAPL',
  interval: 'D' as ResolutionString,
  datafeedUrl: 'https://demo_feed.tradingview.com',
  libraryPath: '/charting_library/',
  chartsStorageUrl: 'https://saveload.tradingview.com',
  chartsStorageApiVersion: '1.1',
  clientId: 'tradingview.com',
  userId: 'public_user_id',
  fullscreen: false,
  autosize: true,
  studiesOverrides: {},
  timezone: 'Asia/Ho_Chi_Minh',
};

export default memo(TVChartContainer);
