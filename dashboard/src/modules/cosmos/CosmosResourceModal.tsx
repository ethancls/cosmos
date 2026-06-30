"use client";

import Button from "@components/Button";
import FancyToggleSwitch from "@components/FancyToggleSwitch";
import { Input } from "@components/Input";
import { Label } from "@components/Label";
import {
  Modal,
  ModalClose,
  ModalContent,
  ModalFooter,
} from "@components/modal/Modal";
import ModalHeader from "@components/modal/ModalHeader";
import { notify } from "@components/Notification";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@components/Select";
import { useApiCall } from "@utils/api";
import React, { useMemo, useState } from "react";
import { ServerIcon } from "lucide-react";
import { PeerGroupSelector } from "@components/PeerGroupSelector";
import { useGroups } from "@/contexts/GroupsProvider";
import {
  CosmosProtocol,
  CosmosResource,
  CosmosResourceRequest,
} from "@/interfaces/Cosmos";
import { Group } from "@/interfaces/Group";

type Props = {
  open: boolean;
  setOpen: (open: boolean) => void;
  resource?: CosmosResource;
  onSaved?: (resource: CosmosResource) => void;
};

const protocolDefaults: Record<CosmosProtocol, number> = {
  ssh: 22,
  rdp: 3389,
  vnc: 5900,
};

export default function CosmosResourceModal({
  open,
  setOpen,
  resource,
  onSaved,
}: Props) {
  return (
    <Modal open={open} onOpenChange={setOpen}>
      <CosmosResourceModalContent
        key={`${open}-${resource?.id ?? "new"}`}
        resource={resource}
        onSaved={(saved) => {
          onSaved?.(saved);
          setOpen(false);
        }}
      />
    </Modal>
  );
}

function CosmosResourceModalContent({
  resource,
  onSaved,
}: {
  resource?: CosmosResource;
  onSaved?: (resource: CosmosResource) => void;
}) {
  const [name, setName] = useState(resource?.name ?? "");
  const [description, setDescription] = useState(resource?.description ?? "");
  const [protocol, setProtocol] = useState<CosmosProtocol>(
    resource?.protocol ?? "ssh",
  );
  const [host, setHost] = useState(resource?.host ?? "");
  const [port, setPort] = useState(String(resource?.port ?? protocolDefaults.ssh));
  const { groups: allGroups } = useGroups();
  const [groupIDs, setGroupIDs] = useState(resource?.group_ids ?? "");
  const [enabled, setEnabled] = useState(resource ? resource.enabled : true);

  const initialGroups = useMemo(() => {
    return groupIDs
      .split(",")
      .map((id) => id.trim())
      .filter(Boolean)
      .map((id) => allGroups?.find((g) => g.id === id))
      .filter(Boolean) as Group[];
  }, [groupIDs, allGroups]);
  const [recordingEnabled, setRecordingEnabled] = useState(
    resource ? resource.recording_enabled : true,
  );

  const create = useApiCall<CosmosResource>("/cosmos/resources").post;
  const update = useApiCall<CosmosResource>(
    `/cosmos/resources/${resource?.id}`,
  ).put;

  const portNumber = Number(port);
  const canSave =
    name.trim().length > 0 &&
    host.trim().length > 0 &&
    Number.isInteger(portNumber) &&
    portNumber > 0 &&
    portNumber <= 65535;

  const payload = useMemo<CosmosResourceRequest>(
    () => ({
      name: name.trim(),
      description: description.trim(),
      protocol,
      host: host.trim(),
      port: portNumber,
      group_ids: groupIDs
        .split(",")
        .map((g) => g.trim())
        .filter(Boolean),
      enabled,
      recording_enabled: recordingEnabled,
    }),
    [
      description,
      enabled,
      host,
      groupIDs,
      name,
      portNumber,
      protocol,
      recordingEnabled,
    ],
  );

  const save = () => {
    const promise = (resource ? update(payload) : create(payload)).then((r) => {
      onSaved?.(r);
      return r;
    });

    notify({
      title: resource ? "Resource Updated" : "Resource Created",
      description: `${name.trim()} is ready for controlled access.`,
      loadingMessage: resource ? "Updating resource..." : "Creating resource...",
      promise,
    });
  };

  const changeProtocol = (value: CosmosProtocol) => {
    setProtocol(value);
    if (!resource) setPort(String(protocolDefaults[value]));
  };

  return (
    <ModalContent maxWidthClass="max-w-2xl">
      <ModalHeader
        icon={<ServerIcon size={19} />}
        title={resource ? "Edit Resource" : "Add Resource"}
        description="Register an SSH, RDP, or VNC target for browser-based bastion access."
        color="cosmos"
      />
      <div className="px-8 py-6 grid gap-5">
        <div>
          <Label>Name</Label>
          <Input
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Production database"
          />
        </div>
        <div>
          <Label>Description</Label>
          <Input
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Optional context for operators"
          />
        </div>
        <div className="grid grid-cols-1 md:grid-cols-[160px_1fr_120px] gap-4">
          <div>
            <Label>Protocol</Label>
            <Select value={protocol} onValueChange={changeProtocol}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="ssh">SSH</SelectItem>
                <SelectItem value="rdp">RDP</SelectItem>
                <SelectItem value="vnc">VNC</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div>
            <Label>Host</Label>
            <Input
              value={host}
              onChange={(e) => setHost(e.target.value)}
              placeholder="10.0.12.15 or server.internal"
            />
          </div>
          <div>
            <Label>Port</Label>
            <Input
              value={port}
              onChange={(e) => setPort(e.target.value)}
              type="number"
              min={1}
              max={65535}
            />
          </div>
        </div>
        <div>
          <Label>Groups</Label>
          <PeerGroupSelector
            onChange={(groups) => {
              const gs = typeof groups === "function" ? groups(initialGroups) : groups;
              setGroupIDs(gs.map((g) => g.id).join(","));
            }}
            values={initialGroups}
          />
        </div>
        <div className="grid gap-3">
          <FancyToggleSwitch
            value={enabled}
            onChange={setEnabled}
            label="Enabled"
            helpText="Allow users with matching policies to start sessions to this resource."
          />
          <FancyToggleSwitch
            value={recordingEnabled}
            onChange={setRecordingEnabled}
            label="Record sessions"
            helpText="Store session metadata now, and enable full recording once the bastion stream is wired."
          />
        </div>
      </div>
      <ModalFooter>
        <ModalClose asChild>
          <Button variant="secondary">Cancel</Button>
        </ModalClose>
        <Button variant="primary" disabled={!canSave} onClick={save}>
          {resource ? "Save Changes" : "Add Resource"}
        </Button>
      </ModalFooter>
    </ModalContent>
  );
}
