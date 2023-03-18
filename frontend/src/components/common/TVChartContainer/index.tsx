import { useEffect, useRef } from "react";

import type {
  ChartingLibraryWidgetOptions,
  IChartingLibraryWidget,
  ResolutionString,
} from "charting_library/charting_library";
import { widget } from "charting_library";

import "./index.css";

interface TVChartContainerProps {
  symbol?: ChartingLibraryWidgetOptions["symbol"];
  interval?: ChartingLibraryWidgetOptions["interval"];
  // BEWARE: no trailing slash is expected in feed URL
  datafeedUrl?: string;
  libraryPath?: ChartingLibraryWidgetOptions["library_path"];
  chartsStorageUrl?: ChartingLibraryWidgetOptions["charts_storage_url"];
  chartsStorageApiVersion?: ChartingLibraryWidgetOptions["charts_storage_api_version"];
  clientId?: ChartingLibraryWidgetOptions["client_id"];
  userId?: ChartingLibraryWidgetOptions["user_id"];
  fullscreen?: ChartingLibraryWidgetOptions["fullscreen"];
  autosize?: ChartingLibraryWidgetOptions["autosize"];
  studiesOverrides?: ChartingLibraryWidgetOptions["studies_overrides"];
  // container?: ChartingLibraryWidgetOptions["container"];
  datafeed: ChartingLibraryWidgetOptions["datafeed"];
  timezone?: ChartingLibraryWidgetOptions["timezone"];
}

function TVChartContainer({
  symbol,
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
  datafeed,
  timezone,
}: TVChartContainerProps): JSX.Element {
  const containerRef = useRef<HTMLDivElement | null>(null);
  const tvWidgetRef = useRef<IChartingLibraryWidget | null>(null);

  useEffect(() => {
    if (containerRef.current !== null) {
      const widgetOptions: ChartingLibraryWidgetOptions = {
        symbol: symbol as string,
        // BEWARE: no trailing slash is expected in feed URL
        // tslint:disable-next-line:no-any
        // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-call, @typescript-eslint/no-unsafe-member-access, @typescript-eslint/no-explicit-any
        datafeed: new (window as any).Datafeeds.UDFCompatibleDatafeed(
          datafeedUrl
        ),
        // datafeed,
        interval: interval || ("D" as ResolutionString),
        container: containerRef.current,
        library_path: libraryPath as string,
        locale: "vi",
        disabled_features: ["use_localstorage_for_settings"],
        enabled_features: ["study_templates"],
        charts_storage_url: chartsStorageUrl,
        charts_storage_api_version: chartsStorageApiVersion,
        client_id: clientId,
        user_id: userId,
        fullscreen,
        autosize,
        studies_overrides: studiesOverrides,
        auto_save_delay: 1,
        theme: "Dark",
        timezone: timezone,
      };

      // eslint-disable-next-line new-cap, no-multi-assign, @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-call, @typescript-eslint/no-unsafe-member-access
      tvWidgetRef.current = new widget(widgetOptions);

      tvWidgetRef.current.onChartReady(() => {
        if (tvWidgetRef.current !== null) {
          tvWidgetRef.current.headerReady().then(() => {
            if (tvWidgetRef.current !== null) {
              const button = tvWidgetRef.current.createButton();
              button.setAttribute(
                "title",
                "Click to show a notification popup"
              );
              button.classList.add("apply-common-tooltip");
              button.addEventListener("click", () => {
                if (tvWidgetRef.current !== null) {
                  tvWidgetRef.current.showNoticeDialog({
                    title: "Notification",
                    body: "TradingView Charting Library API works correctly",
                    callback: () => {
                      console.log("Noticed!");
                    },
                  });
                }
              });
              button.innerHTML = "Check API";
            }
          });
        }
      });
    }

    return () => {
      if (tvWidgetRef.current !== null) {
        tvWidgetRef.current.remove();
        tvWidgetRef.current = null;
      }
    };
  });

  return <div ref={containerRef} className="TVChartContainer" />;
}

TVChartContainer.defaultProps = {
  symbol: "AAPL",
  interval: "D" as ResolutionString,
  datafeedUrl: "https://demo_feed.tradingview.com",
  libraryPath: "/charting_library/",
  chartsStorageUrl: "https://saveload.tradingview.com",
  chartsStorageApiVersion: "1.1",
  clientId: "tradingview.com",
  userId: "public_user_id",
  fullscreen: false,
  autosize: true,
  studiesOverrides: {},
  timezone: "Asia/Ho_Chi_Minh"
};

export default TVChartContainer;
