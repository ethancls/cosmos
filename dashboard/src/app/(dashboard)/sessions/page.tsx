"use client";

import Breadcrumbs from "@components/Breadcrumbs";
import Paragraph from "@components/Paragraph";
import { RestrictedAccess } from "@components/ui/RestrictedAccess";
import { usePortalElement } from "@hooks/usePortalElement";
import useFetchApi from "@utils/api";
import { TerminalSquareIcon } from "lucide-react";
import React from "react";
import { usePermissions } from "@/contexts/PermissionsProvider";
import { CosmosSession } from "@/interfaces/Cosmos";
import PageContainer from "@/layouts/PageContainer";
import CosmosSessionsTable from "@/modules/cosmos/CosmosSessionsTable";

export default function SessionsPage() {
  const { permission } = usePermissions();
  const { data: sessions, isLoading } =
    useFetchApi<CosmosSession[]>("/cosmos/sessions");
  const { ref: headingRef, portalTarget } =
    usePortalElement<HTMLHeadingElement>();

  return (
    <PageContainer>
      <div className="p-default py-6">
        <Breadcrumbs>
          <Breadcrumbs.Item
            href="/sessions"
            label="Sessions"
            icon={<TerminalSquareIcon size={14} />}
            active
          />
        </Breadcrumbs>
        <h1 ref={headingRef}>Sessions</h1>
        <Paragraph>
          Active and historical bastion access sessions across SSH, RDP, and
          VNC.
        </Paragraph>
      </div>
      <RestrictedAccess page="Sessions" hasAccess={permission.events.read}>
        <CosmosSessionsTable
          sessions={sessions}
          isLoading={isLoading}
          headingTarget={portalTarget}
        />
      </RestrictedAccess>
    </PageContainer>
  );
}
