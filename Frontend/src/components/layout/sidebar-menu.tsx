"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Collapse } from "rizzui";
import { Icon } from "@iconify/react";
import { PiCaretDownBold } from "react-icons/pi";
import { cn } from "@/lib/utils";
import { useAuth } from "@/hooks/use-auth";
import { MenuItem } from "@/services/auth/user-schema";

function hasActiveDescendant(item: MenuItem, pathname: string): boolean {
  if (item.href === pathname) return true;
  return (item.children || []).some((child) => hasActiveDescendant(child, pathname));
}

interface MenuNodeProps {
  item: MenuItem;
  pathname: string;
  depth: number;
}

function MenuNode({ item, pathname, depth }: MenuNodeProps) {
  const hasChildren = !!item.children && item.children.length > 0;

  if (hasChildren) {
    const isOpenByDefault = hasActiveDescendant(item, pathname);
    return (
      <Collapse
        defaultOpen={isOpenByDefault}
        header={({ open, toggle }) => (
          <div
            onClick={toggle}
            style={depth > 0 ? { paddingLeft: 12 + depth * 14 } : undefined}
            className={cn(
              "group relative mx-3 flex cursor-pointer items-center justify-between rounded-md px-3 py-2.5 font-medium lg:my-1 2xl:mx-5 2xl:my-2 transition-all duration-200",
              isOpenByDefault
                ? "text-primary bg-primary/10 before:absolute before:left-0 before:top-1/2 before:-translate-y-1/2 before:h-3/5 before:w-1 before:rounded-r-md before:bg-primary"
                : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
            )}
          >
            <span className="flex items-center truncate">
              {item.icon && (
                <span
                  className={cn(
                    "me-2.5 inline-flex h-5 w-5 shrink-0 items-center justify-center rounded-md transition-colors duration-200",
                    isOpenByDefault
                      ? "text-primary"
                      : "text-muted-foreground group-hover:text-accent-foreground"
                  )}
                >
                  <Icon icon={item.icon} width={20} height={20} />
                </span>
              )}
              <span className="truncate">{item.name}</span>
            </span>

            <PiCaretDownBold
              strokeWidth={3}
              className={cn(
                "h-3.5 w-3.5 shrink-0 -rotate-90 text-muted-foreground transition-transform duration-200 rtl:rotate-90",
                open && "rotate-0 rtl:rotate-0"
              )}
            />
          </div>
        )}
      >
        {item.children!.map((child) => (
          <MenuNode key={child.id} item={child} pathname={pathname} depth={depth + 1} />
        ))}
      </Collapse>
    );
  }

  const isActive = item.href === pathname;

  if (depth === 0) {
    return (
      <Link
        href={item.href || "#"}
        className={cn(
          "group relative mx-3 my-0.5 flex items-center justify-between rounded-md px-3 py-2.5 font-medium lg:my-1 2xl:mx-5 2xl:my-2 transition-all duration-200",
          isActive
            ? "text-primary bg-primary/10 font-semibold before:absolute before:left-0 before:top-1/2 before:-translate-y-1/2 before:h-3/5 before:w-1 before:rounded-r-md before:bg-primary"
            : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
        )}
      >
        <div className="flex items-center truncate">
          {item.icon && (
            <span
              className={cn(
                "me-2.5 inline-flex h-5 w-5 shrink-0 items-center justify-center rounded-md transition-colors duration-200",
                isActive ? "text-primary" : "text-muted-foreground group-hover:text-accent-foreground"
              )}
            >
              <Icon icon={item.icon} width={20} height={20} />
            </span>
          )}
          <span className="truncate">{item.name}</span>
        </div>
      </Link>
    );
  }

  return (
    <Link
      href={item.href || "#"}
      style={{ paddingLeft: 12 + depth * 14 }}
      className={cn(
        "group mx-3.5 mb-0.5 flex items-center justify-between rounded-md py-2 pr-2.5 font-medium last-of-type:mb-1 lg:last-of-type:mb-2 2xl:mx-5 transition-all duration-200",
        isActive
          ? "text-primary bg-primary/10 font-semibold"
          : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
      )}
    >
      <div className="flex items-center truncate">
        <span
          className={cn(
            "me-[18px] inline-flex h-1.5 w-1.5 shrink-0 rounded-full transition-all duration-200",
            isActive
              ? "bg-primary ring-[1px] ring-primary"
              : "bg-muted-foreground/40 group-hover:bg-muted-foreground/80"
          )}
        />
        <span className="truncate">{item.name}</span>
      </div>
    </Link>
  );
}

export function SidebarMenu() {
  const pathname = usePathname();
  const { user } = useAuth();
  const items = user?.menuItems || [];

  return (
    <div className="mt-4 pb-3 3xl:mt-6">
      {items.map((item) => (
        <MenuNode key={item.id} item={item} pathname={pathname} depth={0} />
      ))}
    </div>
  );
}
