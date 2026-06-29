"use client";

import { ScrollArea } from "@components/ScrollArea";
import { cn } from "@utils/helpers";
import {
  LayoutDashboard,
  ScrollText,
  Server,
  Settings,
  ShieldCheck,
  Terminal,
  Users,
} from "lucide-react";
import DocsIcon from "@/assets/icons/DocsIcon";
import IntegrationIcon from "@/assets/icons/IntegrationIcon";
import SidebarItem from "@/components/SidebarItem";
import { NavigationVersionInfo } from "@/components/VersionInfo";
import { useAnnouncement } from "@/contexts/AnnouncementProvider";
import { useApplicationContext } from "@/contexts/ApplicationProvider";
import { headerHeight } from "@/layouts/Header";
import * as React from "react";

type Props = {
  fullWidth?: boolean;
  hideOnMobile?: boolean;
};

export default function Navigation({
  fullWidth = false,
  hideOnMobile = false,
}: Readonly<Props>) {
  const { bannerHeight } = useAnnouncement();
  const { isNavigationCollapsed } = useApplicationContext();

  return (
    <div
      data-navigation
      className={cn(
        "whitespace-nowrap md:border-r dark:border-zinc-700/40 bg-gray-50 dark:bg-nb-gray relative group/navigation transition-all",
        hideOnMobile ? "hidden md:block" : "",
        fullWidth
          ? "w-auto max-w-[22rem]"
          : "w-[15rem] max-w-[15rem] min-w-[15rem] overflow-y-auto",
        isNavigationCollapsed &&
          "md:w-[70px] md:min-w-[70px] md:fixed md:overflow-hidden md:hover:w-[15rem] md:hover:max-w-[15rem] md:hover:min-w-[15rem] md:z-50",
      )}
      style={{
        height: `calc(100vh - ${headerHeight + bannerHeight}px)`,
      }}
    >
      <div className={cn(fullWidth ? "w-10/12" : "fixed z-0")}>
        <ScrollArea
          style={{
            height: !fullWidth
              ? `calc(100vh - ${headerHeight + bannerHeight}px)`
              : "100%",
          }}
        >
          <div
            className={cn(
              "flex flex-col pt-3 justify-between w-[15rem] max-w-[15rem] min-w-[15rem] transition-all",
              isNavigationCollapsed &&
                "md:w-[70px] md:min-w-[70px] md:group-hover/navigation:w-[15rem] md:group-hover/navigation:max-w-[15rem] md:group-hover/navigation:min-w-[15rem] md:overflow-x-clip",
            )}
            style={{
              height: !fullWidth
                ? `calc(100vh - ${headerHeight + bannerHeight}px)`
                : "100%",
            }}
          >
            <div>
              <SidebarItemGroup>
                <SidebarItem
                  icon={<LayoutDashboard size={16} />}
                  label="Dashboard"
                  href={"/"}
                  visible={true}
                />

                <SidebarItem
                  icon={<Server size={16} />}
                  label="Servers"
                  href={"/servers"}
                  visible={true}
                />

                <SidebarItem
                  icon={<Terminal size={16} />}
                  label="Connections"
                  href={"/connections"}
                  visible={true}
                />

                <SidebarItem
                  icon={<Users size={16} />}
                  label="Users"
                  href={"/users"}
                  collapsible
                  visible={true}
                >
                  <SidebarItem
                    label="Users"
                    isChild
                    href={"/users"}
                    exactPathMatch={true}
                    visible={true}
                  />
                  <SidebarItem
                    label="Service Accounts"
                    isChild
                    href={"/users/service-accounts"}
                    visible={true}
                  />
                </SidebarItem>

                <SidebarItem
                  icon={<ShieldCheck size={16} />}
                  label="Policies"
                  href={"/policies"}
                  visible={true}
                />

                <AuditNavigationItem />
              </SidebarItemGroup>

              <SidebarItemGroup>
                <SidebarItem
                  icon={<Settings size={16} />}
                  label="Settings"
                  href={"/settings"}
                  exactPathMatch={true}
                  visible={true}
                />
                <SidebarItem
                  icon={<IntegrationIcon />}
                  label="Integrations"
                  href={"/integrations"}
                  exactPathMatch={true}
                  visible={true}
                />
                <SidebarItem
                  icon={<DocsIcon />}
                  href={"/docs"}
                  label="Documentation"
                  visible={true}
                />
              </SidebarItemGroup>
            </div>
            <NavigationVersionInfo />
          </div>
        </ScrollArea>
      </div>
    </div>
  );
}

type SidebarItemGroupProps = {
  children: React.ReactNode;
};

export function SidebarItemGroup({ children }: SidebarItemGroupProps) {
  return (
    <div
      className={
        "mt-4 border-t border-gray-200 pt-4 first:mt-0 first:border-t-0 first:pt-0 dark:border-zinc-700/40 space-y-[3px]"
      }
    >
      {children}
    </div>
  );
}

const AuditNavigationItem = () => {
  return (
    <SidebarItem
      icon={<ScrollText size={16} />}
      label="Audit Logs"
      href={"/audit"}
      collapsible
      visible={true}
    >
      <SidebarItem
        label="Connection Logs"
        href={"/audit/connections"}
        isChild
        exactPathMatch={true}
        visible={true}
      />
      <SidebarItem
        label="Traffic Events"
        isChild
        href={"/audit/traffic"}
        exactPathMatch={true}
        visible={true}
      />
    </SidebarItem>
  );
};
