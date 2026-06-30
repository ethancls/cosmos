"use client";

import Breadcrumbs from "@components/Breadcrumbs";
import Paragraph from "@components/Paragraph";
import { RestrictedAccess } from "@components/ui/RestrictedAccess";
import { usePortalElement } from "@hooks/usePortalElement";
import useFetchApi from "@utils/api";
import { ServerIcon } from "lucide-react";
import React from "react";
import { usePermissions } from "@/contexts/PermissionsProvider";
import { CosmosResource } from "@/interfaces/Cosmos";
import PageContainer from "@/layouts/PageContainer";
import CosmosResourcesTable from "@/modules/cosmos/CosmosResourcesTable";

export default function ResourcesPage() {
  const { permission } = usePermissions();
  const { data: resources, isLoading } =
    useFetchApi<CosmosResource[]>("/cosmos/resources");
  const { ref: headingRef, portalTarget } =
    usePortalElement<HTMLHeadingElement>();

  return (
    <PageContainer>
      <div className="p-default py-6">
        <Breadcrumbs>
          <Breadcrumbs.Item
            href="/resources"
            label="Resources"
            icon={<ServerIcon size={14} />}
            active
          />
        </Breadcrumbs>
        <h1 ref={headingRef}>Resources</h1>
        <Paragraph>
          SSH, RDP, and VNC targets available through the Cosmos bastion.
        </Paragraph>
      </div>
      <RestrictedAccess page="Resources" hasAccess={permission.policies.read}>
        <CosmosResourcesTable
          resources={resources}
          isLoading={isLoading}
          headingTarget={portalTarget}
        />
      </RestrictedAccess>
    </PageContainer>
  );
}
