export type CosmosProtocol = "ssh" | "rdp" | "vnc";

export type CosmosSessionStatus = "pending" | "active" | "closed" | "denied";

export type CosmosResource = {
  id: string;
  account_id: string;
  name: string;
  description?: string;
  protocol: CosmosProtocol;
  host: string;
  port: number;
  group_ids?: string;
  enabled: boolean;
  recording_enabled: boolean;
  created_at: string;
  updated_at: string;
};

export type CosmosResourceRequest = {
  name: string;
  description?: string;
  protocol: CosmosProtocol;
  host: string;
  port: number;
  group_ids?: string[];
  enabled: boolean;
  recording_enabled: boolean;
};

export type CosmosSession = {
  id: string;
  account_id: string;
  resource_id: string;
  resource_name: string;
  user_id: string;
  user_name?: string;
  user_email?: string;
  protocol: CosmosProtocol;
  status: CosmosSessionStatus;
  client_ip?: string;
  guacd_connection_id?: string;
  started_at: string;
  ended_at?: string;
  created_at: string;
  updated_at: string;
};

export type CosmosAuditEvent = {
  id: string;
  account_id: string;
  user_id?: string;
  user_name?: string;
  user_email?: string;
  action: string;
  target_type: string;
  target_id?: string;
  target_name?: string;
  timestamp: string;
  description?: string;
  meta?: string;
  created_at: string;
};
