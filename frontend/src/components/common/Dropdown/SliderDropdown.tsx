/* eslint-disable react/jsx-props-no-spreading */
import type { FC } from 'react';
import { Fragment, memo, useState, useRef } from 'react';
import cx from 'classnames';
import { Menu, Transition } from '@headlessui/react';
import { ChevronDownIcon } from '@heroicons/react/20/solid';

import {
  AddIcon,
  WatchListIcon,
  RemoveIcon,
  EditIcon,
} from 'components/common/Icon';

export interface Item {
  id: number;
  name: string;
}

interface WatchListsProps {
  currentItem: Item;
  items: Item[];
  onChange: (item: Item) => void;
}

interface WatchListItemProps {
  item: Item;
  isActive: boolean;
  isEdit: boolean;
  onClick: () => void;
  onEdit: () => void;
}

const WatchListItem: FC<WatchListItemProps> = ({
  item,
  isActive,
  isEdit,
  onClick,
  onEdit,
}) => {
  const inputRef = useRef<HTMLInputElement | null>(null);

  return (
    <Menu.Item key={item.id}>
      {({ active }) => (
        <div
          className={cx(
            'group flex w-full justify-between items-center rounded-sm px-2 py-2 text-sm font-sans font-normal text-gray-100',
            {
              'bg-slate-600': active || isActive,
            }
          )}
        >
          <button
            type="button"
            className="flex flex-row items-center justify-center mr-auto"
            onClick={() => {
              if (!isEdit) {
                onClick();
              }
            }}
            disabled={!isEdit}
          >
            <div className="pr-1.5 pt-[2px]">
              <WatchListIcon />
            </div>
            <input
              ref={inputRef}
              type="text"
              className="bg-transparent border-none focus:border-none"
              value={item.name}
              disabled={isEdit}
            />
          </button>

          <div className="flex flex-row items-center justify-center">
            <button
              type="button"
              className="mx-1 px-0.5 py-0.5 border-none rounded-sm cursor-pointer shadow-md bg-blue-500"
              onClick={() => {
                inputRef.current?.focus();
                onEdit();
              }}
            >
              <EditIcon />
            </button>
            <button
              type="button"
              className="mx-1 px-0.5 py-0.5 border-none rounded-sm cursor-pointer shadow-md bg-red-500"
            >
              <RemoveIcon />
            </button>
          </div>
        </div>
      )}
    </Menu.Item>
  );
};

const WatchLists: FC<WatchListsProps> = ({ currentItem, items, onChange }) => {
  const [isOpen, setIsOpen] = useState(true);
  const [itemUpdate, setItemUpdate] = useState<number | null>(null);

  return (
    <div className="w-full h-full text-right bg-slate-800">
      <Menu as="div" className="w-full h-full relative inline-block text-left">
        <div className="w-full h-full">
          <Menu.Button
            className={cx(
              'inline-flex w-full h-full justify-between items-center rounded-md px-4 py-2',
              'text-sm font-sans font-normal',
              'text-gray-100',
              'hover:bg-slate-700',
              'focus:outline-none focus-visible:ring-2 focus-visible:ring-white focus-visible:ring-opacity-75'
            )}
          >
            {currentItem.name}
            <ChevronDownIcon
              className={cx(
                'ml-2 -mr-1 h-5 w-5',
                'text-gray-100',
                'hover:text-violet-100'
              )}
              aria-hidden="true"
            />
          </Menu.Button>
        </div>

        <Transition
          as={Fragment}
          show={isOpen}
          enter="transition ease-out duration-100"
          enterFrom="transform opacity-0 scale-95"
          enterTo="transform opacity-100 scale-100"
          leave="transition ease-in duration-75"
          leaveFrom="transform opacity-100 scale-100"
          leaveTo="transform opacity-0 scale-95"
        >
          <Menu.Items
            static
            className={cx(
              'absolute right-0 mt-1 px-1.5 w-full origin-top-right',
              'divide-y divide-white divide-opacity-20',
              'rounded-md shadow-lg',
              'bg-slate-700',
              'ring-1 ring-black ring-opacity-5',
              'focus:outline-none'
            )}
          >
            <div className="py-1">
              {items.map((item) => (
                <WatchListItem
                  item={item}
                  isActive={currentItem.id === item.id}
                  isEdit={item.id === itemUpdate}
                  onClick={() => onChange(item)}
                  onEdit={() => setItemUpdate(item.id)}
                />
              ))}
            </div>

            <div className="py-1">
              <Menu.Item>
                {({ active }) => (
                  <button
                    type="button"
                    className={cx(
                      'group flex flex-row w-full items-center rounded-sm px-2 py-2 text-sm font-sans font-normal text-gray-100',
                      {
                        'bg-slate-600': active,
                      }
                    )}
                    onClick={() => {}}
                  >
                    <div className="pr-1 pt-[2px]">
                      <AddIcon />
                    </div>
                    Tạo watchlist mới
                  </button>
                )}
              </Menu.Item>
            </div>
          </Menu.Items>
        </Transition>
      </Menu>
    </div>
  );
};

export default memo(WatchLists);
