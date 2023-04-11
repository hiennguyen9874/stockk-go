import { useEffect, useState } from 'react';
import cx from 'classnames';

import { SearchIcon } from 'components/common/Icon';
import { Modal } from 'components/common/Modal';
import { useSearchBySymbolQuery } from 'app/services/ticker';

import WatchListItem from './WatchListItem';

interface SearchProps {
  isOpen: boolean;
  onClose: () => void;
  onAddRemove: (symbol: string) => void;
}

const Search = ({ isOpen, onClose, onAddRemove }: SearchProps) => {
  const [inputValue, setInputValue] = useState('');

  const { data: tickers } = useSearchBySymbolQuery(inputValue, {
    skip: inputValue === '',
  });

  useEffect(() => {
    setInputValue('');
  }, [isOpen]);

  return (
    <Modal isOpen={isOpen} onClose={() => onClose()}>
      <div className="bg-slate-800 max-h-96 flex flex-col overflow-hidden">
        <div className="flex flex-row items-center text-white shadow-2xl shadow-slate-900 rounded-md">
          <div className="h-10 w-10 p-2.5 text-gray-400">
            <SearchIcon className="w-full h-full " />
          </div>

          <div className="grow w-full h-full">
            <input
              type="text"
              placeholder="Thêm mã CK vào watchlist"
              value={inputValue}
              onChange={(e) => setInputValue(e.target.value)}
              onKeyDown={(e) => {
                if (e.code === 'Enter') {
                  e.preventDefault();
                  e.stopPropagation();
                  if (inputValue !== '' && tickers && tickers.data.length > 0) {
                    onAddRemove(tickers.data[0].symbol);
                  }
                }
              }}
              className={cx(
                'w-full h-full text-md text-gray-300 leading-10 shadow-none rounded-none bg-transparent uppercase',
                'focus:outline-none',
                'placeholder:text-gray-600 placeholder:normal-case'
              )}
            />
          </div>
        </div>

        <div className="grow h-full overflow-y-auto">
          {inputValue !== '' &&
            tickers &&
            tickers.data.map((ticker, idx) => (
              <WatchListItem
                symbol={ticker.symbol}
                isLight={idx % 2 === 0}
                onSet={() => onAddRemove(ticker.symbol)}
                className="h-16 pl-4 pt-2 pr-4 pb-2"
              />
            ))}
        </div>
      </div>
    </Modal>
  );
};

export default Search;
