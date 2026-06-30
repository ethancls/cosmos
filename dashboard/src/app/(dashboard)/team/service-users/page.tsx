"use client";

import Breadcrumbs from "@components/Breadcrumbs";
import Paragraph from "@components/Paragraph";
import SkeletonTable from "@components/skeletons/SkeletonTable";
import { RestrictedAccess } from "@components/ui/RestrictedAccess";
import { usePortalElement } from "@hooks/usePortalElement";
import { IconSettings2 } from "@tabler/icons-react";
import useFetchApi from "@utils/api";
import React, { lazy, Suspense } from "react";
import TeamIcon from "@/assets/icons/TeamIcon";
import { usePermissions } from "@/contexts/PermissionsProvider";
import { User } from "@/interfaces/User";
import PageContainer from "@/layouts/PageContainer";

const ServiceUsersTable = lazy(
  () => import("@/modules/users/ServiceUsersTable"),
);

export default function ServiceUsers() {
  const { permission } = usePermissions();
  const { data: users, isLoading } = useFetchApi<User[]>(
    "/users?service_user=true",
  );

  const { ref: headingRef, portalTarget } =
    usePortalElement<HTMLHeadingElement>();

  return (
    <PageContainer>
      <div className={"p-default py-6"}>
        <Breadcrumbs>
          <Breadcrumbs.Item
            href={"/team"}
            label={"Team"}
            icon={<TeamIcon size={13} />}
          />
          <Breadcrumbs.Item
            href={"/team/service-users"}
            label={"Service Accounts"}
            active
            icon={<IconSettings2 size={17} />}
          />
        </Breadcrumbs>
        <h1 ref={headingRef}>Service Accounts</h1>
        <Paragraph>
          Use service accounts to create API tokens for automation and gateway
          integrations.
        </Paragraph>
      </div>
      <RestrictedAccess
        page={"Service Users"}
        hasAccess={permission.users.read}
      >
        <Suspense fallback={<SkeletonTable />}>
          <ServiceUsersTable
            users={users}
            isLoading={isLoading}
            headingTarget={portalTarget}
          />
        </Suspense>
      </RestrictedAccess>
    </PageContainer>
  );
}
