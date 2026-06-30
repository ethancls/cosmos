"use client";

import Breadcrumbs from "@components/Breadcrumbs";
import Button from "@components/Button";
import Badge from "@components/Badge";
import Card from "@components/Card";
import { Label } from "@components/Label";
import Paragraph from "@components/Paragraph";
import FullScreenLoading from "@components/ui/FullScreenLoading";
import PageContainer from "@/layouts/PageContainer";
import useFetchApi, { useApiCall } from "@utils/api";
import {
  CircleDotIcon,
  MonitorIcon,
  ServerIcon,
  StopCircleIcon,
  TerminalSquareIcon,
} from "lucide-react";
import { useParams, useRouter } from "next/navigation";
import React, { useMemo } from "react";
import dayjs from "dayjs";
import { useSWRConfig } from "swr";
import { notify } from "@components/Notification";
import { cn } from "@utils/helpers";
import { CosmosSession, CosmosResource } from "@/interfaces/Cosmos";

type SessionAccessResponse = {
  session: CosmosSession;
  resource: CosmosResource;
  connection: {
    protocol: string;
    host: string;
    port: number;
  };
};

export default function SessionPage() {
  const params = useParams<{ id: string }>();
  const router = useRouter();
  const { mutate } = useSWRConfig();
  const { data, isLoading, error } = useFetchApi<SessionAccessResponse>(
    `/cosmos/sessions/${params.id}/access`,
  );

  const closeApi = useApiCall<CosmosSession>(
    `/cosmos/sessions/${params.id}/close`,
  ).post;

  const session = data?.session;
  const resource = data?.resource;
  const connection = data?.connection;

  const protocolIcon = useMemo(() => {
    const p = session?.protocol;
    if (p === "ssh") return <TerminalSquareIcon size={18} />;
    if (p === "rdp") return <MonitorIcon size={18} />;
    return <ServerIcon size={18} />;
  }, [session?.protocol]);

  const protocolBadgeVariant =
    session?.protocol === "ssh"
      ? "green"
      : session?.protocol === "rdp"
        ? "blue"
        : "purple";

  const handleClose = async () => {
    await closeApi({});
    mutate("/cosmos/sessions");
    notify({
      title: "Session Closed",
      description: `${session?.resource_name} session has been closed.`,
      loadingMessage: "Closing session...",
    });
  };

  if (isLoading) return <FullScreenLoading />;
  if (error || !session) {
    return (
      <PageContainer>
        <div className="p-default py-6">
          <h1>Session not found</h1>
          <Button variant="secondary" onClick={() => router.push("/sessions")}>
            Back to Sessions
          </Button>
        </div>
      </PageContainer>
    );
  }

  return (
    <PageContainer>
      <div className="p-default py-6">
        <Breadcrumbs>
          <Breadcrumbs.Item href="/sessions" label="Sessions" active />
        </Breadcrumbs>
        <div className="flex items-center gap-4">
          <h1>{resource?.name || session.resource_name}</h1>
          <Badge variant={protocolBadgeVariant} className="uppercase font-medium">
            {protocolIcon}
            <span className="ml-1">{session.protocol.toUpperCase()}</span>
          </Badge>
          <Badge
            variant={session.status === "active" ? "green" : "yellow"}
            className="uppercase font-medium"
          >
            {session.status}
          </Badge>
        </div>
        <Paragraph>
          {session.status === "active"
            ? "Session is active. Connect using the details below."
            : `Session ${session.status}.`}
        </Paragraph>
      </div>

      <div className="p-default grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card className="p-6">
          <h2 className="text-lg font-medium mb-4">Connection Details</h2>
          <div className="grid gap-3">
            <div>
              <Label>Host</Label>
              <p className="font-mono text-lg text-nb-gray-100">
                {connection?.host}
              </p>
            </div>
            <div>
              <Label>Port</Label>
              <p className="font-mono text-lg text-nb-gray-100">
                {connection?.port}
              </p>
            </div>
            <div>
              <Label>Protocol</Label>
              <p className="text-lg text-nb-gray-100 uppercase">
                {connection?.protocol}
              </p>
            </div>
          </div>
        </Card>

        <Card className="p-6">
          <h2 className="text-lg font-medium mb-4">Session Info</h2>
          <div className="grid gap-3">
            <div>
              <Label>User</Label>
              <p className="text-nb-gray-100">
                {session.user_name} ({session.user_email || session.user_id})
              </p>
            </div>
            <div>
              <Label>Client IP</Label>
              <p className="font-mono text-nb-gray-100">
                {session.client_ip || "-"}
              </p>
            </div>
            <div>
              <Label>Started</Label>
              <p className="text-nb-gray-100">
                {dayjs(session.started_at).format("YYYY-MM-DD HH:mm:ss")}
              </p>
            </div>
            <div>
              <Label>Duration</Label>
              <p className="text-nb-gray-100">
                {session.ended_at
                  ? `${dayjs(session.ended_at).diff(dayjs(session.started_at), "minute")} min`
                  : `${dayjs().diff(dayjs(session.started_at), "minute")} min (active)`}
              </p>
            </div>
            <div>
              <Label>Recording</Label>
              <div className="flex items-center gap-2">
                {session.recording_enabled ? (
                  <Badge variant="green" className="uppercase font-medium">
                    <CircleDotIcon size={12} />
                    Recording
                  </Badge>
                ) : (
                  <span className="text-nb-gray-400 text-sm">
                    Not recording
                  </span>
                )}
              </div>
            </div>
          </div>
        </Card>
      </div>

      {session.status === "active" && (
        <div className="p-default pt-6">
          <Card className="p-6 text-center">
            <h2 className="text-lg font-medium mb-2">Bastion Connection</h2>
            <Paragraph className="mb-4">
              Browser-based {session.protocol.toUpperCase()} access via Apache
              Guacamole. The connection is proxied through the Cosmos bastion.
            </Paragraph>
            <div
              className={cn(
                "inline-flex items-center justify-center rounded-lg p-8 bg-nb-gray-900 border border-nb-gray-800",
                "w-full max-w-2xl min-h-[400px]",
              )}
            >
              <div className="text-center">
                <div className="mb-4">{protocolIcon}</div>
                <p className="text-nb-gray-300 text-sm">
                  Guacamole client will load here.
                  <br />
                  Connection to {connection?.host}:{connection?.port} via{" "}
                  {connection?.protocol?.toUpperCase()}.
                </p>
              </div>
            </div>
            <div className="mt-4">
              <Button
                variant="danger-outline"
                size="sm"
                onClick={handleClose}
              >
                <StopCircleIcon size={14} />
                End Session
              </Button>
            </div>
          </Card>
        </div>
      )}
    </PageContainer>
  );
}
