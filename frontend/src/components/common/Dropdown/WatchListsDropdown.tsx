/* eslint-disable react/jsx-props-no-spreading */
import type { FC } from 'react';
import { Fragment, memo, useState, useRef, useEffect } from 'react';
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

interface WatchListItemProps {
  item: Item;
  isActive: boolean;
  isEdit: boolean;
  onClick: () => void;
  onEdit: () => void;
  onDelete: () => Promise<void>;
  onChange: (value: string) => Promise<void>;
}

const WatchListItem: FC<WatchListItemProps> = ({
  item,
  isActive,
  isEdit,
  onClick,
  onEdit,
  onDelete,
  onChange,
}) => {
  const inputRef = useRef<HTMLInputElement | null>(null);

  useEffect(() => {
    if (isEdit) {
      inputRef.current?.focus();
    }
  }, [isEdit]);

  return (
    <Menu.Item key={item.id}>
      {({ active }) => (
        <div
          className={cx(
            'w-full flex flex-row items-center rounded-sm text-sm font-sans font-normal text-gray-100',
            {
              'bg-slate-600': active || isActive,
            }
          )}
        >
          <button
            type="button"
            className="flex flex-row items-center justify-center px-2 py-2"
            onClick={() => {
              if (!isEdit) {
                onClick();
              }
            }}
            disabled={isEdit}
          >
            <div className="pr-1.5 pt-[2px]">
              <WatchListIcon />
            </div>

            <div className="grow">
              <input
                ref={inputRef}
                type="text"
                className="w-full bg-transparent border-none focus:border-none disabled:cursor-pointer"
                value={item.name}
                disabled={!isEdit}
                // eslint-disable-next-line @typescript-eslint/no-misused-promises
                onChange={async (e) => {
                  e.preventDefault();
                  e.stopPropagation();
                  await onChange(e.target.value);
                }}
                onClick={(e) => {
                  if (isEdit) {
                    e.preventDefault();
                    e.stopPropagation();
                  }
                }}
              />
              {/* <span>A</span> */}
            </div>
          </button>

          <div className="flex flex-row items-center justify-center px-2 py-2">
            <button
              type="button"
              className="mx-1 px-0.5 py-0.5 border-none rounded-sm cursor-pointer shadow-md bg-blue-500"
              onClick={(e) => {
                e.preventDefault();
                e.stopPropagation();
                onEdit();
              }}
            >
              <EditIcon />
            </button>
            <button
              type="button"
              className="mx-1 px-0.5 py-0.5 border-none rounded-sm cursor-pointer shadow-md bg-red-500"
              // eslint-disable-next-line @typescript-eslint/no-misused-promises
              onClick={async (e) => {
                e.preventDefault();
                e.stopPropagation();
                await onDelete();
              }}
            >
              <RemoveIcon />
            </button>
          </div>
        </div>
      )}
    </Menu.Item>
  );
};

interface WatchListsProps {
  currentItem: Item | null;
  items: Item[];
  onClick: (item: Item) => void;
  onEdit: (id: number, value: string) => Promise<void>;
  onDelete: (id: number) => Promise<void>;
  onAdd: () => Promise<void>;
}

const WatchListsDropdown: FC<WatchListsProps> = ({
  currentItem,
  items,
  onClick,
  onEdit,
  onDelete,
  onAdd,
}) => {
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
            {currentItem !== null && currentItem.name}
            {currentItem === null && 'Không có danh mục nào'}
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
          enter="transition ease-out duration-100"
          enterFrom="transform opacity-0 scale-95"
          enterTo="transform opacity-100 scale-100"
          leave="transition ease-in duration-75"
          leaveFrom="transform opacity-100 scale-100"
          leaveTo="transform opacity-0 scale-95"
          beforeEnter={() => {
            setItemUpdate(null);
          }}
          afterLeave={() => {
            setItemUpdate(null);
          }}
        >
          <Menu.Items
            className={cx(
              'absolute w-full h-[50vh] right-0 mt-1 px-1.5 origin-top-right',
              'divide-y divide-white divide-opacity-20',
              'rounded-md shadow-lg',
              'bg-slate-700',
              'ring-1 ring-black ring-opacity-5',
              'focus:outline-none',
              'flex flex-col'
            )}
          >
            <div className="w-full grow overflow-y-auto py-1 bg-red">
              {items.map((item) => (
                <WatchListItem
                  key={item.id}
                  item={item}
                  isActive={currentItem !== null && currentItem.id === item.id}
                  isEdit={item.id === itemUpdate}
                  onClick={() => {
                    setItemUpdate(null);
                    onClick(item);
                  }}
                  onEdit={() => {
                    if (itemUpdate === null || itemUpdate !== item.id)
                      setItemUpdate(item.id);
                    else setItemUpdate(null);
                  }}
                  onDelete={async () => {
                    setItemUpdate(null);
                    await onDelete(item.id);
                  }}
                  onChange={async (value) => {
                    await onEdit(item.id, value);
                  }}
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
                    onClick={(e) => {
                      e.preventDefault();
                      e.stopPropagation();
                      onAdd();
                    }}
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

export default memo(WatchListsDropdown);
