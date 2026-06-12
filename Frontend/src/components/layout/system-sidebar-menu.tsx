'use client';

import Link from 'next/link';
import { Fragment } from 'react';
import { usePathname } from 'next/navigation';
import { Title, Collapse } from 'rizzui';
import { cn } from '@/lib/utils';
import { PiCaretDownBold } from 'react-icons/pi';
import { menuItems } from './system-menu-items';

export function SystemSidebarMenu() {
  const pathname = usePathname();

  return (
    <div className="mt-4 pb-3 3xl:mt-6">
      {menuItems.map((item, index) => {
        const isActive = pathname === (item?.href as string);
        const pathnameExistInDropdowns = item?.dropdownItems?.some(
          (dropdownItem) => dropdownItem.href === pathname
        );
        const isDropdownOpen = Boolean(pathnameExistInDropdowns);

        return (
          <Fragment key={item.name + '-' + index}>
            {item?.href ? (
              <>
                {item?.dropdownItems ? (
                  <Collapse
                    defaultOpen={isDropdownOpen}
                    header={({ open, toggle }) => (
                      <div
                        onClick={toggle}
                        className={cn(
                          'group relative mx-3 flex cursor-pointer items-center justify-between rounded-md px-3 py-2.5 font-medium lg:my-1 2xl:mx-5 2xl:my-2 transition-all duration-200',
                          isDropdownOpen
                            ? 'text-primary bg-primary/10 before:absolute before:left-0 before:top-1/2 before:-translate-y-1/2 before:h-3/5 before:w-1 before:rounded-r-md before:bg-primary'
                            : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
                        )}
                      >
                        <span className="flex items-center">
                          {item?.icon && (
                            <span
                              className={cn(
                                'me-2.5 inline-flex h-5 w-5 items-center justify-center rounded-md [&>svg]:h-[20px] [&>svg]:w-[20px] transition-colors duration-200',
                                isDropdownOpen
                                  ? 'text-primary'
                                  : 'text-muted-foreground group-hover:text-accent-foreground'
                              )}
                            >
                              {item?.icon}
                            </span>
                          )}
                          {item.name}
                        </span>

                        <PiCaretDownBold
                          strokeWidth={3}
                          className={cn(
                            'h-3.5 w-3.5 -rotate-90 text-muted-foreground transition-transform duration-200 rtl:rotate-90',
                            open && 'rotate-0 rtl:rotate-0'
                          )}
                        />
                      </div>
                    )}
                  >
                    {item?.dropdownItems?.map((dropdownItem, dropdownIndex) => {
                      const isChildActive =
                        pathname === (dropdownItem?.href as string);

                      return (
                        <Link
                          href={dropdownItem?.href || '#'}
                          key={dropdownItem?.name + dropdownIndex}
                          className={cn(
                            'group mx-3.5 mb-0.5 flex items-center justify-between rounded-md px-2.5 py-2 font-medium last-of-type:mb-1 lg:last-of-type:mb-2 2xl:mx-5 2xl:px-3.5 transition-all duration-200',
                            isChildActive
                              ? 'text-primary bg-primary/10 font-semibold'
                              : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
                          )}
                        >
                          <div className="flex items-center truncate">
                            <span
                              className={cn(
                                'me-[18px] ms-1 inline-flex h-1.5 w-1.5 rounded-full transition-all duration-200',
                                isChildActive
                                  ? 'bg-primary ring-[1px] ring-primary'
                                  : 'bg-muted-foreground/40 group-hover:bg-muted-foreground/80'
                              )}
                            />{' '}
                            <span className="truncate">
                              {dropdownItem?.name}
                            </span>
                          </div>
                        </Link>
                      );
                    })}
                  </Collapse>
                ) : (
                  <Link
                    href={item?.href}
                    className={cn(
                      'group relative mx-3 my-0.5 flex items-center justify-between rounded-md px-3 py-2.5 font-medium lg:my-1 2xl:mx-5 2xl:my-2 transition-all duration-200',
                      isActive
                        ? 'text-primary bg-primary/10 font-semibold before:absolute before:left-0 before:top-1/2 before:-translate-y-1/2 before:h-3/5 before:w-1 before:rounded-r-md before:bg-primary'
                        : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
                    )}
                  >
                    <div className="flex items-center truncate">
                      {item?.icon && (
                        <span
                          className={cn(
                            'me-2.5 inline-flex h-5 w-5 items-center justify-center rounded-md transition-colors duration-200 [&>svg]:h-[20px] [&>svg]:w-[20px]',
                            isActive
                              ? 'text-primary'
                              : 'text-muted-foreground group-hover:text-accent-foreground'
                          )}
                        >
                          {item?.icon}
                        </span>
                      )}
                      <span className="truncate">{item.name}</span>
                    </div>
                  </Link>
                )}
              </>
            ) : (
              <Title
                as="h6"
                className={cn(
                  'mb-2 truncate px-6 text-xs font-semibold uppercase tracking-widest text-muted-foreground/60 2xl:px-8',
                  index !== 0 && 'mt-6 3xl:mt-7'
                )}
              >
                {item.name}
              </Title>
            )}
          </Fragment>
        );
      })}
    </div>
  );
}
