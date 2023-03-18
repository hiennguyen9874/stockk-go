import debounce from 'lodash/debounce';

import type {
  OnReadyCallback,
  DatafeedConfiguration,
  ResolutionString,
  SearchSymbolsCallback,
  ResolveCallback,
  SymbolResolveExtension,
  SearchSymbolResultItem,
  LibrarySymbolInfo,
  PeriodParams,
  HistoryCallback,
  Bar,
  SubscribeBarsCallback,
  ErrorCallback,
} from 'charting_library';
import * as tcbsService from 'api/tcbs';

const exchangeToValue: {
  [key: string]: string;
} = {
  '0': 'HOSE',
  '1': 'HNX',
  '3': 'UPCOM',
};

const configurationData: DatafeedConfiguration = {
  supported_resolutions: ['D', 'W', 'M'] as ResolutionString[],
  exchanges: [
    { value: '', name: 'All Exchanges', desc: '' },
    {
      // `exchange` argument for the `searchSymbols` method, if a user selects this exchange
      value: 'HOSE',
      // filter name
      name: 'HOSE',
      // full exchange name displayed in the filter popup
      desc: 'HOSE',
    },
    {
      // `exchange` argument for the `searchSymbols` method, if a user selects this exchange
      value: 'HNX',
      // filter name
      name: 'HNX',
      // full exchange name displayed in the filter popup
      desc: 'HNX',
    },
    {
      // `exchange` argument for the `searchSymbols` method, if a user selects this exchange
      value: 'UPCOM',
      // filter name
      name: 'UPCOM',
      // full exchange name displayed in the filter popup
      desc: 'UPCOM',
    },
  ],
  symbols_types: [
    { name: 'All types', value: '' },
    {
      name: 'Stock',
      value: 'stock',
    },
    {
      name: 'Crypto',
      value: 'crypto',
    },
  ],
  supports_marks: true,
  supports_timescale_marks: true,
  supports_time: true,
};

let allSymbols: {
  [key: string]: SearchSymbolResultItem;
} = {};

export default {
  onReady: (callback: OnReadyCallback): void => {
    console.log('[onReady]: Method call');
    setTimeout(() => callback(configurationData), 0);
  },
  searchSymbols: (
    userInput: string,
    exchange: string,
    symbolType: string,
    onResult: SearchSymbolsCallback
  ): void => {
    console.log('[searchSymbols]: Method call');

    if (userInput !== '') {
      debounce(
        // eslint-disable-next-line @typescript-eslint/no-shadow
        (userInput: string) =>
          tcbsService
            .searchSymbolsByKey(userInput)
            .then(({ data }) => {
              data.data.forEach((symbol) => {
                if (Object.keys(exchangeToValue).includes(symbol.exchange)) {
                  allSymbols = {
                    ...allSymbols,
                    [symbol.name]: {
                      symbol: symbol.name,
                      full_name: symbol.value,
                      description: symbol.value,
                      exchange: exchangeToValue[symbol.exchange],
                      ticker: symbol.name,
                      type: symbol.type,
                    },
                  };
                }
              });

              onResult(
                data.data
                  .filter((symbol) => {
                    // const isExchangeValid =
                    //   exchange === '' ||
                    //   (Object.keys(exchangeToValue).includes(symbol.exchange) &&
                    //     exchangeToValue[symbol.exchange].toLowerCase() ===
                    //       exchange.toLowerCase());
                    const isExchangeValid = true;

                    const isFullSymbolContainsInput =
                      symbol.name
                        .toLowerCase()
                        .indexOf(userInput.toLowerCase()) !== -1;

                    const isTypeValid =
                      symbolType === '' || symbolType === symbol.type;

                    return (
                      isExchangeValid &&
                      isFullSymbolContainsInput &&
                      isTypeValid
                    );
                  })
                  .map((symbol) => ({
                    symbol: symbol.name,
                    full_name: symbol.value,
                    description: symbol.value,
                    exchange: exchangeToValue[symbol.exchange],
                    ticker: symbol.name,
                    type: symbol.type,
                  }))
              );
            })
            .catch((error) => {
              console.log(error);
            }),
        300
      )(userInput);
    }
  },
  resolveSymbol: (
    symbolName: string,
    onResolve: ResolveCallback,
    onError: ErrorCallback,
    extension?: SymbolResolveExtension
  ): void => {
    console.log('[resolveSymbol]: Method call', symbolName);
    const symbolItem = Object.values(allSymbols).find(
      (symbol) => symbol.symbol === symbolName
    );

    if (symbolItem === undefined) {
      console.log('[resolveSymbol]: Cannot resolve symbol', symbolName);
      onError('cannot resolve symbol');
      return;
    }

    const symbolInfo: LibrarySymbolInfo = {
      ticker: symbolItem.ticker,
      name: symbolItem.symbol,
      description: symbolItem.description,
      type: symbolItem.type,
      full_name: symbolItem.full_name,
      exchange: symbolItem.exchange,
      listed_exchange: symbolItem.exchange,
      session: '0930-1630',
      timezone: 'Asia/Ho_Chi_Minh',
      minmov: 1,
      minmove2: 0,
      pricescale: 100,
      has_intraday: false,
      has_no_volume: false,
      // has_weekly_and_monthly: true,
      supported_resolutions:
        configurationData.supported_resolutions ??
        (['D'] as ResolutionString[]),
      // volume_precision: 2,
      // data_status: 'pulsed',
      format: 'price',
    };

    setTimeout(() => {
      console.log('[resolveSymbol]: Symbol resolved', symbolName);
      onResolve(symbolInfo);
    }, 0);
  },
  getBars: (
    symbolInfo: LibrarySymbolInfo,
    resolution: ResolutionString,
    periodParams: PeriodParams,
    onResult: HistoryCallback,
    onError: ErrorCallback
  ): void => {
    const { from, to, firstDataRequest } = periodParams;
    console.log('[getBars]: Method call', symbolInfo, resolution, from, to);

    const bars: Bar[] = [];
    setTimeout(() => onResult(bars, { noData: false }), 0);
  },
  subscribeBars: (
    symbolInfo: LibrarySymbolInfo,
    resolution: ResolutionString,
    onTick: SubscribeBarsCallback,
    listenerGuid: string,
    onResetCacheNeededCallback: () => void
  ): void => {
    console.log(
      '[subscribeBars]: Method call with subscriberUID:',
      listenerGuid
    );
  },
  unsubscribeBars: (listenerGuid: string): void => {
    console.log(
      '[unsubscribeBars]: Method call with listenerGuid:',
      listenerGuid
    );
  },
};
