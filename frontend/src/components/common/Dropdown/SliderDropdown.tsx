/* eslint-disable react/jsx-props-no-spreading */
import type { FC } from 'react';
import { Fragment, memo } from 'react';
import cx from 'classnames';
import { Menu, Transition } from '@headlessui/react';
import { ChevronDownIcon } from '@heroicons/react/20/solid';

interface SliderDropdownProps {
  currentItem: string;
  items: string[];
  onChange: (item: string) => void;
}

const SliderDropdown: FC<SliderDropdownProps> = ({
  currentItem,
  items,
  onChange,
}) => {
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
            {currentItem}
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
        >
          <Menu.Items
            className={cx(
              'absolute right-0 mt-1 w-full origin-top-right',
              'rounded-md shadow-lg',
              'bg-slate-700',
              'ring-1 ring-black ring-opacity-5',
              'focus:outline-none'
            )}
          >
            <div className="px-2 py-1 ">
              {items.map((item) => (
                <Menu.Item key={item}>
                  {({ active }) => (
                    <button
                      type="button"
                      className={cx(
                        'group flex w-full items-center rounded-sm px-2 py-2 text-sm font-sans font-normal text-gray-100',
                        {
                          'bg-slate-600': active || currentItem === item,
                        }
                      )}
                      onClick={() => onChange(item)}
                    >
                      {item}
                    </button>
                  )}
                </Menu.Item>
              ))}
            </div>
          </Menu.Items>
        </Transition>
      </Menu>
    </div>
  );
};

export default memo(SliderDropdown);
