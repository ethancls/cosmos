"use client";

import { ScrollArea } from "@components/ScrollArea";
import { cn } from "@utils/helpers";
import { isAgentNetworkOnly } from "@utils/netbird";
import AccessControlIcon from "@/assets/icons/AccessControlIcon";
import ControlCenterIcon from "@/assets/icons/ControlCenterIcon";
import SettingsIcon from "@/assets/icons/SettingsIcon";
import TeamIcon from "@/assets/icons/TeamIcon";
import SidebarItem from "@/components/SidebarItem";
import { useAnnouncement } from "@/contexts/AnnouncementProvider";
import { useApplicationContext } from "@/contexts/ApplicationProvider";
import { usePermissions } from "@/contexts/PermissionsProvider";
import { headerHeight } from "@/layouts/Header";
import * as React from "react";
import ActivityIcon from "@/assets/icons/ActivityIcon";
import { ServerIcon, TerminalSquareIcon } from "lucide-react";

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
  const { permission, isRestricted } = usePermissions();

  return (
    <div
      data-navigation
      className={cn(
        "whitespace-nowrap md:border-r dark:border-zinc-700/40 bg-gray-50 dark:bg-nb-gray relative group/navigation",
        hideOnMobile ? "hidden md:block" : "",
        fullWidth
          ? "w-auto max-w-[22rem]"
          : "w-[15rem] max-w-[15rem] min-w-[15rem] overflow-y-auto",
        isNavigationCollapsed &&
          "md:w-[70px] md:min-w-[70px] md:fixed md:overflow-hidden md:hover:w-[15rem] md:hover:min-w-[15rem] md:hover:overflow-visible md:z-50 md:shadow-none md:hover:shadow-2xl",
        "transition-[width,min-width,box-shadow] duration-300 ease-out",
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
              "flex flex-col pt-3 justify-between w-[15rem] max-w-[15rem] min-w-[15rem]",
              isNavigationCollapsed &&
                "md:w-[70px] md:min-w-[70px] md:group-hover/navigation:w-[15rem] md:group-hover/navigation:min-w-[15rem] md:overflow-hidden md:group-hover/navigation:overflow-visible",
              "transition-[width,min-width] duration-300 ease-out",
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
                  icon={<ControlCenterIcon size={16} />}
                  label="Control Center"
                  href={"/control-center"}
                  visible={permission.policies.read}
                />

                <SidebarItem
                  icon={<ServerIcon size={16} />}
                  label="Resources"
                  href={"/resources"}
                  visible={permission.policies.read && !isRestricted}
                />

                <SidebarItem
                  icon={<TerminalSquareIcon size={16} />}
                  label="Sessions"
                  href={"/sessions"}
                  visible={permission.events.read}
                />

                <SidebarItem
                  icon={<AccessControlIcon />}
                  label="Access"
                  href={"/access-control"}
                  collapsible
                  visible={permission.policies.read}
                >
                  <SidebarItem
                    label="Policies"
                    href={"/access-control"}
                    isChild
                    exactPathMatch={true}
                    visible={permission.policies.read}
                  />
                  <SidebarItem
                    label="Groups"
                    isChild
                    href={"/groups"}
                    visible={permission.policies.read}
                  />
                  <SidebarItem
                    label="Posture Checks"
                    isChild
                    href={"/posture-checks"}
                    exactPathMatch={true}
                    visible={permission.policies.read}
                  />
                </SidebarItem>

                <SidebarItem
                  icon={<TeamIcon />}
                  label="Team"
                  href={"/team"}
                  collapsible
                  visible={permission.users.read}
                >
                  <SidebarItem
                    label="Users"
                    isChild
                    href={"/team/users"}
                    visible={permission.users.read}
                  />
                  <SidebarItem
                    label="Service Accounts"
                    isChild
                    href={"/team/service-users"}
                    visible={permission.users.read}
                  />
                </SidebarItem>
                <ActivityNavigationItem />
              </SidebarItemGroup>

              <SidebarItemGroup>
                <SidebarItem
                  icon={<SettingsIcon />}
                  label="Settings"
                  href={"/settings"}
                  exactPathMatch={true}
                  visible={permission.settings.read}
                />
              </SidebarItemGroup>
            </div>
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

const ActivityNavigationItem = () => {
  const { permission } = usePermissions();

  return (
    <SidebarItem
      icon={<ActivityIcon />}
      label="Activity"
      href={"/events"}
      collapsible
      visible={permission.events.read && !isAgentNetworkOnly()}
    >
      <SidebarItem
        label="Audit Events"
        href={"/events/audit"}
        isChild
        exactPathMatch={true}
        visible={permission.events.read}
      />
      <SidebarItem
        label="Traffic Events"
        isChild
        href={"/events/traffic"}
        exactPathMatch={true}
        visible={permission.events.read}
      />
    </SidebarItem>
  );
};
