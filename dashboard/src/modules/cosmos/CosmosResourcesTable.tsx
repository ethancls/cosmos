"use client";

import Button from "@components/Button";
import Badge from "@components/Badge";
import CopyToClipboardText from "@components/CopyToClipboardText";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@components/DropdownMenu";
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
  formatRadioChip,
  RadioOption,
  RadioPicker,
} from "@components/table/filters/RadioPicker";
import {
  TableFilterChips,
  TableFilterDef,
  TableFiltersButton,
} from "@components/table/TableFilters";
import GetStartedTest from "@components/ui/GetStartedTest";
import DescriptionWithTooltip from "@components/ui/DescriptionWithTooltip";
import TextWithTooltip from "@components/ui/TextWithTooltip";
import { notify } from "@components/Notification";
import { ToggleSwitch } from "@components/ToggleSwitch";
import { ColumnDef, SortingState } from "@tanstack/react-table";
import { useApiCall } from "@utils/api";
import { cn } from "@utils/helpers";
import dayjs from "dayjs";
import {
  CircleDotIcon,
  ExternalLinkIcon,
  MonitorIcon,
  MoreVertical,
  PowerIcon,
  PlayIcon,
  PlusCircleIcon,
  ServerIcon,
  SquarePenIcon,
  TerminalSquareIcon,
  Trash2,
} from "lucide-react";
import React, { useCallback, useMemo, useState } from "react";
import { useSWRConfig } from "swr";
import { useDialog } from "@/contexts/DialogProvider";
import { useGroups } from "@/contexts/GroupsProvider";
import { useLocalStorage } from "@/hooks/useLocalStorage";
import { CosmosResource, CosmosSession } from "@/interfaces/Cosmos";
import CosmosResourceModal from "@/modules/cosmos/CosmosResourceModal";

type Props = {
  resources?: CosmosResource[];
  isLoading: boolean;
  headingTarget?: HTMLHeadingElement | null;
};

const protocolLabel = (protocol: CosmosResource["protocol"]) =>
  protocol.toUpperCase();

export default function CosmosResourcesTable({
  resources,
  isLoading,
  headingTarget,
}: Props) {
  const { mutate } = useSWRConfig();
  const { confirm } = useDialog();
  const [modalOpen, setModalOpen] = useState(false);
  const [resourceToEdit, setResourceToEdit] = useState<CosmosResource>();
  const startSession = useApiCall<CosmosSession>("/cosmos/sessions").post;
  const resourceApi = useApiCall<CosmosResource>("/cosmos/resources");

  const [sorting, setSorting] = useLocalStorage<SortingState>(
    "cosmos-table-sort/resources",
    [{ id: "name", desc: false }],
  );

  const openCreate = useCallback(() => {
    setResourceToEdit(undefined);
    setModalOpen(true);
  }, []);

  const openEdit = useCallback((resource: CosmosResource) => {
    setResourceToEdit(resource);
    setModalOpen(true);
  }, []);

  const statusOptions = useMemo<RadioOption<boolean | undefined>[]>(
    () => [
      { value: undefined, label: "All", dotClass: "bg-nb-gray-500" },
      { value: true, label: "Enabled", dotClass: "bg-green-500" },
      { value: false, label: "Disabled", dotClass: "bg-nb-gray-700" },
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

  const recordingOptions = useMemo<RadioOption<boolean | undefined>[]>(
    () => [
      { value: undefined, label: "All", dotClass: "bg-nb-gray-500" },
      { value: true, label: "Recorded", dotClass: "bg-green-500" },
      { value: false, label: "Not recorded", dotClass: "bg-nb-gray-700" },
    ],
    [],
  );

  const filterDefs = useMemo<TableFilterDef[]>(
    () => [
      {
        id: "enabled_filter",
        label: "Status",
        renderPicker: (p) => (
          <RadioPicker
            value={p.value as boolean | undefined}
            onChange={p.onChange}
            close={p.close}
            options={statusOptions}
          />
        ),
        formatChip: (v) =>
          formatRadioChip(v as boolean | undefined, statusOptions),
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
      {
        id: "recording_filter",
        label: "Recording",
        renderPicker: (p) => (
          <RadioPicker
            value={p.value as boolean | undefined}
            onChange={p.onChange}
            close={p.close}
            options={recordingOptions}
          />
        ),
        formatChip: (v) =>
          formatRadioChip(v as boolean | undefined, recordingOptions),
      },
    ],
    [protocolOptions, recordingOptions, statusOptions],
  );

  const start = useCallback(
    (resource: CosmosResource) => {
      const promise = startSession({ resource_id: resource.id }).then(
        (session) => {
          mutate("/cosmos/sessions");
          return session;
        },
      );

      notify({
        title: "Session Started",
        description: `${resource.name} is now tracked as an active bastion session.`,
        loadingMessage: "Starting session...",
        promise,
      });
    },
    [mutate, startSession],
  );

  const toggleEnabled = useCallback(
    (resource: CosmosResource) => {
      const nextEnabled = !resource.enabled;
      const promise = resourceApi
        .put(
          {
            name: resource.name,
            description: resource.description,
            protocol: resource.protocol,
            host: resource.host,
            port: resource.port,
            group_ids: resource.group_ids
              ?.split(",")
              .map((g) => g.trim())
              .filter(Boolean),
            enabled: nextEnabled,
            recording_enabled: resource.recording_enabled,
          },
          `/${resource.id}`,
        )
        .then((updated) => {
          mutate("/cosmos/resources");
          return updated;
        });

      notify({
        title: "Update Resource",
        description: `'${resource.name}' is now ${
          nextEnabled ? "enabled" : "disabled"
        }.`,
        loadingMessage: "Updating resource...",
        duration: 1200,
        promise,
      });
    },
    [mutate, resourceApi],
  );

  const deleteResource = useCallback(
    async (resource: CosmosResource) => {
      const choice = await confirm({
        title: `Delete '${resource.name}'?`,
        description:
          "Are you sure you want to delete this resource? This action cannot be undone.",
        confirmText: "Delete",
        cancelText: "Cancel",
        type: "danger",
      });
      if (!choice) return;

      const promise = resourceApi.del(undefined, `/${resource.id}`).then((r) => {
        mutate("/cosmos/resources");
        return r;
      });

      notify({
        title: "Delete Resource",
        description: `'${resource.name}' has been deleted.`,
        loadingMessage: "Deleting resource...",
        promise,
      });
    },
    [confirm, mutate, resourceApi],
  );

  const columns = useMemo<ColumnDef<CosmosResource>[]>(
    () => [
      {
        id: "name",
        accessorFn: (resource) => resource.name,
        header: ({ column }) => (
          <DataTableHeader column={column}>Resource</DataTableHeader>
        ),
        sortingFn: "text",
        cell: ({ row }) => (
          <button
            className={"flex gap-4 items-center group min-w-[220px]"}
            onClick={() => openEdit(row.original)}
          >
            <div
              className={cn(
                "flex items-center justify-center rounded-md h-9 w-9 shrink-0 bg-nb-gray-900 transition-all",
                "group-hover:bg-nb-gray-800",
              )}
            >
              {row.original.protocol === "ssh" ? (
                <TerminalSquareIcon size={15} />
              ) : row.original.protocol === "rdp" ? (
                <MonitorIcon size={15} />
              ) : (
                <ServerIcon size={15} />
              )}
            </div>
            <div
              className={cn(
                "flex flex-col gap-0 text-neutral-300 font-light truncate",
                "group-hover:text-neutral-100 text-left",
              )}
            >
              <TextWithTooltip
                text={row.original.name}
                maxChars={28}
                className={"font-normal"}
              />
              <DescriptionWithTooltip
                maxChars={28}
                className={"font-normal"}
                text={row.original.description}
              />
            </div>
          </button>
        ),
      },
      {
        id: "protocol",
        accessorFn: (resource) => resource.protocol,
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
                {protocolLabel(p)}
              </Badge>
            </div>
          );
        },
      },
      {
        id: "protocol_filter",
        accessorFn: (resource) => [resource.protocol],
        filterFn: "arrIncludesSome",
      },
      {
        id: "enabled_filter",
        accessorFn: (resource) => resource.enabled,
      },
      {
        id: "enabled",
        accessorFn: (resource) => resource.enabled,
        header: ({ column }) => (
          <DataTableHeader column={column}>Enabled</DataTableHeader>
        ),
        cell: ({ row }) => (
          <div className={"flex"}>
            <ToggleSwitch
              checked={row.original.enabled}
              size={"small"}
              onClick={() => toggleEnabled(row.original)}
            />
          </div>
        ),
      },
      {
        id: "target",
        accessorFn: (resource) => `${resource.host}:${resource.port}`,
        header: ({ column }) => (
          <DataTableHeader column={column}>Target</DataTableHeader>
        ),
        cell: ({ row }) => (
          <CopyToClipboardText
            message={`${row.original.host}:${row.original.port} has been copied to your clipboard`}
          >
            <div
              className={
                "font-mono dark:text-nb-gray-300 pt-1 flex gap-2 items-center text-[.82rem]"
              }
            >
              {row.original.host}:{row.original.port}
            </div>
          </CopyToClipboardText>
        ),
      },
      {
        id: "groups",
        accessorFn: (resource) => resource.group_ids || "",
        header: ({ column }) => (
          <DataTableHeader column={column}>Groups</DataTableHeader>
        ),
        cell: ({ row }) => (
          <ResourceGroupCell groupIds={row.original.group_ids} />
        ),
      },
      {
        id: "recording",
        accessorFn: (resource) => resource.recording_enabled,
        header: ({ column }) => (
          <DataTableHeader column={column}>Recording</DataTableHeader>
        ),
        cell: ({ row }) =>
          row.original.recording_enabled ? (
            <Badge
              variant={"green"}
              className={"uppercase font-medium"}
            >
              <CircleDotIcon size={12} />
              Recorded
            </Badge>
          ) : (
            <span className="text-sm text-nb-gray-400">-</span>
          ),
      },
      {
        id: "recording_filter",
        accessorFn: (resource) => resource.recording_enabled,
      },
      {
        id: "updated_at",
        accessorFn: (resource) => resource.updated_at,
        header: ({ column }) => (
          <DataTableHeader column={column}>Updated</DataTableHeader>
        ),
        sortingFn: "datetime",
        cell: ({ row }) => (
          <span className="text-sm text-nb-gray-300">
            {dayjs(row.original.updated_at).fromNow()}
          </span>
        ),
      },
      {
        id: "actions",
        header: "",
        cell: ({ row }) => (
          <div className="flex justify-end gap-2">
            <Button
              size="xs"
              variant="secondary"
              disabled={!row.original.enabled}
              onClick={() => start(row.original)}
            >
              <PlayIcon size={14} />
              Connect
            </Button>
            <CosmosResourceActionMenu
              resource={row.original}
              onEdit={openEdit}
              onToggle={toggleEnabled}
              onDelete={deleteResource}
            />
          </div>
        ),
      },
    ],
    [deleteResource, openEdit, start, toggleEnabled],
  );

  return (
    <>
      <DataTable
        headingTarget={headingTarget}
        text="Resources"
        sorting={sorting}
        setSorting={setSorting}
        initialPageSize={25}
        showResetFilterButton={false}
        columns={columns}
        data={resources}
        isLoading={isLoading}
        searchPlaceholder="Search by name, target, protocol, labels..."
        tableClassName="px-8 pt-4"
        rowClassName={(row) => (row.original.enabled ? "" : "opacity-50")}
        aboveTable={(table) => (
          <TableFilterChips table={table} filters={filterDefs} />
        )}
        columnVisibility={{
          protocol_filter: false,
          enabled_filter: false,
          recording_filter: false,
        }}
        getStartedCard={
          <GetStartedTest
            icon={
              <SquareIcon
                icon={<ServerIcon size={23} />}
                color="gray"
                size="large"
              />
            }
            title="Create Resource"
            description="Add SSH, RDP, or VNC resources to make them available through the Cosmos bastion."
            button={
              <Button variant="primary" onClick={openCreate}>
                <PlusCircleIcon size={16} />
                Create Resource
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
        rightSide={() => (
          <>
            {resources && resources.length > 0 && (
              <div className={"flex ml-auto gap-4 items-center"}>
                <Button variant="primary" size="sm" onClick={openCreate}>
                  <PlusCircleIcon size={16} />
                  Add Resource
                </Button>
              </div>
            )}
          </>
        )}
      >
        {(table) => (
          <>
            <TableFiltersButton
              table={table}
              filters={filterDefs}
              disabled={resources?.length == 0}
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
              isDisabled={resources?.length == 0}
              onClick={() => mutate("/cosmos/resources")}
            />
          </>
        )}
      </DataTable>
      <CosmosResourceModal
        open={modalOpen}
        setOpen={setModalOpen}
        resource={resourceToEdit}
        onSaved={() => mutate("/cosmos/resources")}
      />
    </>
  );
}

type ActionMenuProps = {
  resource: CosmosResource;
  onEdit: (resource: CosmosResource) => void;
  onToggle: (resource: CosmosResource) => void;
  onDelete: (resource: CosmosResource) => void;
};

function CosmosResourceActionMenu({
  resource,
  onEdit,
  onToggle,
  onDelete,
}: ActionMenuProps) {
  return (
    <DropdownMenu modal={false}>
      <DropdownMenuTrigger
        asChild={true}
        onClick={(e) => {
          e.stopPropagation();
          e.preventDefault();
        }}
      >
        <Button
          variant={"secondary"}
          className={"!px-3"}
          aria-label={"Resource actions"}
        >
          <MoreVertical size={16} className={"shrink-0"} />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-auto" align="end">
        <DropdownMenuItem onClick={() => onEdit(resource)}>
          <div className={"flex gap-3 items-center"}>
            <SquarePenIcon size={14} className={"shrink-0"} />
            Edit
          </div>
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => onToggle(resource)}>
          <div className={"flex gap-3 items-center"}>
            <PowerIcon size={14} className={"shrink-0"} />
            {resource.enabled ? "Disable" : "Enable"}
          </div>
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={() => onDelete(resource)} variant={"danger"}>
          <div className={"flex gap-3 items-center"}>
            <Trash2 size={14} className={"shrink-0"} />
            Delete
          </div>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

function ResourceGroupCell({ groupIds }: { groupIds?: string }) {
  const { groups: allGroups } = useGroups();

  const ids = groupIds
    ?.split(",")
    .map((g) => g.trim())
    .filter(Boolean);

  if (!ids || ids.length === 0)
    return <span className="text-sm text-nb-gray-400">-</span>;

  const groupNames = ids
    .map((id) => allGroups?.find((g) => g.id === id)?.name)
    .filter(Boolean) as string[];

  return (
    <div className="flex flex-wrap gap-1">
      {groupNames.map((name) => (
        <Badge key={name} variant={"gray"}>
          {name}
        </Badge>
      ))}
    </div>
  );
}
