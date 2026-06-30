"use client";

import Badge from "@components/Badge";
import Button from "@components/Button";
import InlineLink from "@components/InlineLink";
import SquareIcon from "@components/SquareIcon";
import { DataTable } from "@components/table/DataTable";
import DataTableHeader from "@components/table/DataTableHeader";
import DataTableRefreshButton from "@components/table/DataTableRefreshButton";
import DataTableResetFilterButton from "@components/table/DataTableResetFilterButton";
import {
  CheckboxListPicker,
  CheckboxOption,
  formatCheckboxChip,
} from "@components/table/filters/CheckboxListPicker";
import {
  TableFilterChips,
  TableFilterDef,
  TableFiltersButton,
} from "@components/table/TableFilters";
import { SmallBadge } from "@components/ui/SmallBadge";
import GetStartedTest from "@components/ui/GetStartedTest";
import { notify } from "@components/Notification";
import { ColumnDef, SortingState } from "@tanstack/react-table";
import { useApiCall } from "@utils/api";
import dayjs from "dayjs";
import {
  ExternalLinkIcon,
  MonitorIcon,
  ServerIcon,
  StopCircleIcon,
  TerminalSquareIcon,
} from "lucide-react";
import { useRouter } from "next/navigation";
import React, { useMemo } from "react";
import { useSWRConfig } from "swr";
import { useLocalStorage } from "@/hooks/useLocalStorage";
import { CosmosSession } from "@/interfaces/Cosmos";

type Props = {
  sessions?: CosmosSession[];
  isLoading: boolean;
  headingTarget?: HTMLHeadingElement | null;
};

export default function CosmosSessionsTable({
  sessions,
  isLoading,
  headingTarget,
}: Props) {
  const { mutate } = useSWRConfig();
  const router = useRouter();
  const sessionApi = useApiCall<CosmosSession>("/cosmos/sessions");
  const [sorting, setSorting] = useLocalStorage<SortingState>(
    "cosmos-table-sort/sessions",
    [{ id: "started_at", desc: true }],
  );

  const statusOptions = useMemo<CheckboxOption<string>[]>(
    () => [
      { value: "active", label: "Active" },
      { value: "closed", label: "Closed" },
      { value: "pending", label: "Pending" },
      { value: "denied", label: "Denied" },
    ],
    [],
  );

  const protocolOptions = useMemo<CheckboxOption<string>[]>(
    () => [
      { value: "ssh", label: "SSH" },
      { value: "rdp", label: "RDP" },
      { value: "vnc", label: "VNC" },
    ],
    [],
  );

  const filterDefs = useMemo<TableFilterDef[]>(
    () => [
      {
        id: "status_filter",
        label: "Status",
        renderPicker: (p) => (
          <CheckboxListPicker
            value={p.value as string[] | undefined}
            onChange={p.onChange}
            close={p.close}
            options={statusOptions}
          />
        ),
        formatChip: (v) =>
          formatCheckboxChip(v as string[] | undefined, statusOptions, "statuses"),
      },
      {
        id: "protocol_filter",
        label: "Protocol",
        renderPicker: (p) => (
          <CheckboxListPicker
            value={p.value as string[] | undefined}
            onChange={p.onChange}
            close={p.close}
            options={protocolOptions}
          />
        ),
        formatChip: (v) =>
          formatCheckboxChip(
            v as string[] | undefined,
            protocolOptions,
            "protocols",
          ),
      },
    ],
    [protocolOptions, statusOptions],
  );

  const closeSession = (session: CosmosSession) => {
    const promise = sessionApi
      .post({}, `/${session.id}/close`)
      .then((closed) => {
        mutate("/cosmos/sessions");
        return closed;
      });

    notify({
      title: "Session Closed",
      description: `${session.resource_name} session has been closed.`,
      loadingMessage: "Closing session...",
      promise,
    });
  };

  const columns = useMemo<ColumnDef<CosmosSession>[]>(
    () => [
      {
        id: "resource",
        accessorFn: (session) => session.resource_name,
        header: ({ column }) => (
          <DataTableHeader column={column}>Resource</DataTableHeader>
        ),
        sortingFn: "text",
        cell: ({ row }) => {
          const p = row.original.protocol;
          const Icon =
            p === "ssh" ? TerminalSquareIcon : p === "rdp" ? MonitorIcon : ServerIcon;
          return (
            <div className="flex items-center gap-3 min-w-[220px]">
              <div className="h-9 w-9 rounded-md bg-nb-gray-900 border border-nb-gray-800 flex items-center justify-center text-cosmos">
                <Icon size={17} />
              </div>
              <div className="min-w-0">
                <div className="font-medium text-nb-gray-50 truncate">
                  {row.original.resource_name}
                </div>
                <div className="text-xs text-nb-gray-300 truncate">
                  {row.original.user_name || row.original.user_email || row.original.user_id}
                </div>
              </div>
            </div>
          );
        },
      },
      {
        id: "protocol",
        accessorFn: (session) => session.protocol,
        header: ({ column }) => (
          <DataTableHeader column={column}>Protocol</DataTableHeader>
        ),
        cell: ({ row }) => {
          const p = row.original.protocol;
          const variant =
            p === "ssh" ? "green" : p === "rdp" ? "blue" : "purple";
          const Icon =
            p === "ssh" ? TerminalSquareIcon : p === "rdp" ? MonitorIcon : ServerIcon;
          return (
            <div className={"flex"}>
              <Badge variant={variant} className={"uppercase font-medium"}>
                <Icon size={12} />
                {p.toUpperCase()}
              </Badge>
            </div>
          );
        },
      },
      {
        id: "user",
        accessorFn: (session) => session.user_email || session.user_name || "",
        header: ({ column }) => (
          <DataTableHeader column={column}>User</DataTableHeader>
        ),
        cell: ({ row }) => (
          <div className="text-sm">
            <div className="text-nb-gray-100">
              {row.original.user_name || "Unknown user"}
            </div>
            <div className="text-xs text-nb-gray-400">
              {row.original.user_email || row.original.user_id}
            </div>
          </div>
        ),
      },
      {
        id: "status",
        accessorFn: (session) => session.status,
        header: ({ column }) => (
          <DataTableHeader column={column}>Status</DataTableHeader>
        ),
        cell: ({ row }) => {
          const s = row.original.status;
          const variant =
            s === "active" ? "green" : s === "denied" ? "yellow" : "cosmos";
          return <SmallBadge text={s} variant={variant} />;
        },
      },
      {
        id: "status_filter",
        accessorFn: (session) => [session.status],
        filterFn: "arrIncludesSome",
      },
      {
        id: "protocol_filter",
        accessorFn: (session) => [session.protocol],
        filterFn: "arrIncludesSome",
      },
      {
        id: "client_ip",
        accessorFn: (session) => session.client_ip || "",
        header: ({ column }) => (
          <DataTableHeader column={column}>Client IP</DataTableHeader>
        ),
        cell: ({ row }) => (
          <span className="font-mono text-sm text-nb-gray-300">
            {row.original.client_ip || "-"}
          </span>
        ),
      },
      {
        id: "started_at",
        accessorFn: (session) => session.started_at,
        header: ({ column }) => (
          <DataTableHeader column={column}>Started</DataTableHeader>
        ),
        sortingFn: "datetime",
        cell: ({ row }) => (
          <span className="text-sm text-nb-gray-300">
            {dayjs(row.original.started_at).fromNow()}
          </span>
        ),
      },
      {
        id: "duration",
        accessorFn: (session) => session.ended_at || session.started_at,
        header: ({ column }) => (
          <DataTableHeader column={column}>Duration</DataTableHeader>
        ),
        cell: ({ row }) => {
          const end = row.original.ended_at
            ? dayjs(row.original.ended_at)
            : dayjs();
          const minutes = Math.max(
            0,
            end.diff(dayjs(row.original.started_at), "minute"),
          );
          return (
            <span className="text-sm text-nb-gray-300">
              {minutes < 1 ? "< 1 min" : `${minutes} min`}
            </span>
          );
        },
      },
      {
        id: "actions",
        header: "",
        cell: ({ row }) => (
          <div className="flex justify-end">
            <Button
              size="xs"
              variant="secondary"
              disabled={row.original.status !== "active"}
              onClick={() => closeSession(row.original)}
            >
              <StopCircleIcon size={14} />
              Close
            </Button>
          </div>
        ),
      },
    ],
    [sessionApi],
  );

  return (
    <DataTable
      headingTarget={headingTarget}
      text="Sessions"
      sorting={sorting}
      setSorting={setSorting}
      initialPageSize={25}
      showResetFilterButton={false}
      columns={columns}
      data={sessions}
      isLoading={isLoading}
      searchPlaceholder="Search by resource, user, status, client IP..."
      tableClassName="px-8 pt-4"
      aboveTable={(table) => (
        <TableFilterChips table={table} filters={filterDefs} />
      )}
      columnVisibility={{
        status_filter: false,
        protocol_filter: false,
      }}
      getStartedCard={
        <GetStartedTest
          icon={
            <SquareIcon
              icon={<TerminalSquareIcon size={23} />}
              color="gray"
              size="large"
            />
          }
          title="Start Session"
          description="Sessions appear here when an operator opens an SSH, RDP, or VNC resource through the Cosmos bastion."
          button={
            <Button variant="primary" onClick={() => router.push("/resources")}>
              Go to Resources
            </Button>
          }
          learnMore={
            <>
              Learn more about
              <InlineLink href="/resources">
                Resources
                <ExternalLinkIcon size={12} />
              </InlineLink>
            </>
          }
        />
      }
    >
      {(table) => (
        <>
          <TableFiltersButton
            table={table}
            filters={filterDefs}
            disabled={sessions?.length == 0}
          />
          <DataTableResetFilterButton
            table={table}
            onClick={() => {
              table.setPageIndex(0);
              table.setColumnFilters([]);
              table.setGlobalFilter("");
            }}
          />
          <DataTableRefreshButton
            isDisabled={sessions?.length == 0}
            onClick={() => mutate("/cosmos/sessions")}
          />
        </>
      )}
    </DataTable>
  );
}
