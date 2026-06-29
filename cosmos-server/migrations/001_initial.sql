-- ============================================================================
-- Cosmos: Full Schema — all 8 domains
-- ============================================================================

-- Domain 1: Tenants & Orgs
-- ----------------------------------------------------------------------------

CREATE TYPE tenant_tier AS ENUM ('free', 'pro', 'enterprise');
CREATE TYPE tenant_status AS ENUM ('active', 'suspended', 'deleted');

CREATE TABLE tenants (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          TEXT NOT NULL,
    slug          TEXT NOT NULL UNIQUE,
    tier          tenant_tier NOT NULL DEFAULT 'free',
    status        tenant_status NOT NULL DEFAULT 'active',
    billing_email TEXT,
    settings      JSONB DEFAULT '{}',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ
);

-- Domain 2: Users & Auth
-- ----------------------------------------------------------------------------

CREATE TYPE user_role AS ENUM ('owner', 'admin', 'user', 'service');
CREATE TYPE user_status AS ENUM ('active', 'invited', 'suspended', 'deleted');
CREATE TYPE mfa_type AS ENUM ('totp', 'webauthn', 'sms');
CREATE TYPE identity_provider AS ENUM ('google', 'github', 'azure', 'oidc', 'saml');

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL REFERENCES tenants(id),
    email         TEXT NOT NULL,
    display_name  TEXT,
    password_hash TEXT,
    mfa_enabled   BOOLEAN NOT NULL DEFAULT false,
    mfa_type      mfa_type,
    role          user_role NOT NULL DEFAULT 'user',
    status        user_status NOT NULL DEFAULT 'active',
    last_login_at TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ,
    UNIQUE (tenant_id, email)
);

CREATE TABLE user_identities (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider      identity_provider NOT NULL,
    provider_id   TEXT NOT NULL,
    provider_data JSONB DEFAULT '{}',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (provider, provider_id)
);

CREATE TABLE mfa_devices (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type       mfa_type NOT NULL,
    secret     TEXT NOT NULL,
    verified   BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE api_tokens (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID NOT NULL REFERENCES tenants(id),
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name         TEXT NOT NULL,
    token_hash   TEXT NOT NULL UNIQUE,
    last_used_at TIMESTAMPTZ,
    expires_at   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_groups (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    name        TEXT NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, name)
);

CREATE TABLE user_group_members (
    user_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES user_groups(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, group_id)
);

-- Domain 3: Machines & Infrastructure
-- ----------------------------------------------------------------------------

CREATE TYPE server_status AS ENUM ('online', 'offline', 'unreachable', 'archived');

CREATE TABLE servers (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL REFERENCES tenants(id),
    name          TEXT NOT NULL,
    hostname      TEXT,
    ip_address    INET,
    os            TEXT,
    os_version    TEXT,
    agent_version TEXT,
    status        server_status NOT NULL DEFAULT 'offline',
    last_seen_at  TIMESTAMPTZ,
    metadata      JSONB DEFAULT '{}',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ
);

CREATE TABLE server_labels (
    server_id UUID NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
    key       TEXT NOT NULL,
    value     TEXT NOT NULL,
    PRIMARY KEY (server_id, key)
);

CREATE TABLE server_access_keys (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    server_id   UUID NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    public_key  TEXT NOT NULL,
    fingerprint TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE server_groups (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  UUID NOT NULL REFERENCES tenants(id),
    name       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, name)
);

CREATE TABLE server_group_servers (
    group_id  UUID NOT NULL REFERENCES server_groups(id) ON DELETE CASCADE,
    server_id UUID NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, server_id)
);

CREATE TYPE service_account_status AS ENUM ('active', 'revoked');

CREATE TABLE service_accounts (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID NOT NULL REFERENCES tenants(id),
    server_id    UUID NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
    name         TEXT NOT NULL,
    description  TEXT,
    public_key   TEXT NOT NULL,
    fingerprint  TEXT NOT NULL,
    status       service_account_status NOT NULL DEFAULT 'active',
    last_used_at TIMESTAMPTZ,
    expires_at   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Domain 4: Connexions & Protocoles
-- ----------------------------------------------------------------------------

CREATE TYPE connection_protocol AS ENUM ('ssh', 'rdp', 'vnc', 'telnet', 'k8s', 'db');
CREATE TYPE connection_status AS ENUM ('pending', 'active', 'disconnected', 'failed');

CREATE TABLE connections (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL REFERENCES tenants(id),
    user_id       UUID NOT NULL REFERENCES users(id),
    server_id     UUID NOT NULL REFERENCES servers(id),
    protocol      connection_protocol NOT NULL,
    port          INTEGER NOT NULL DEFAULT 22,
    guacd_conn_id TEXT,
    status        connection_status NOT NULL DEFAULT 'pending',
    client_ip     INET,
    metadata      JSONB DEFAULT '{}',
    started_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    ended_at      TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Domain 5: Zero Trust Policies
-- ----------------------------------------------------------------------------

CREATE TYPE policy_action AS ENUM ('allow', 'deny');
CREATE TYPE policy_status AS ENUM ('enabled', 'disabled');

CREATE TABLE policies (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    name        TEXT NOT NULL,
    description TEXT,
    priority    INTEGER NOT NULL DEFAULT 0,
    action      policy_action NOT NULL DEFAULT 'allow',
    status      policy_status NOT NULL DEFAULT 'enabled',
    conditions  JSONB NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Domain 6: Audit & Compliance (DORA/RGPD/ISO 27001)
-- ----------------------------------------------------------------------------

CREATE TYPE audit_event_type AS ENUM (
    'auth', 'connection', 'policy', 'user', 'server',
    'billing', 'admin', 'api', 'system'
);
CREATE TYPE audit_severity AS ENUM ('info', 'warning', 'critical');

CREATE TABLE audit_logs (
    id            UUID DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL REFERENCES tenants(id),
    actor_type    TEXT NOT NULL,
    actor_id      UUID,
    event_type    audit_event_type NOT NULL,
    event_action  TEXT NOT NULL,
    resource_type TEXT NOT NULL,
    resource_id   UUID,
    ip_address    INET,
    user_agent    TEXT,
    changes       JSONB,
    metadata      JSONB DEFAULT '{}',
    severity      audit_severity NOT NULL DEFAULT 'info',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

CREATE TYPE recording_status AS ENUM ('recording', 'completed', 'deleted');

CREATE TABLE session_recordings (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(id),
    connection_id   UUID NOT NULL REFERENCES connections(id),
    user_id         UUID NOT NULL REFERENCES users(id),
    server_id       UUID NOT NULL REFERENCES servers(id),
    protocol        connection_protocol NOT NULL,
    storage_path    TEXT NOT NULL,
    size_bytes      BIGINT DEFAULT 0,
    duration_ms     BIGINT DEFAULT 0,
    status          recording_status NOT NULL DEFAULT 'recording',
    retention_until TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TYPE traffic_event_type AS ENUM ('keystroke', 'file_transfer', 'clipboard', 'command', 'screen');

CREATE TABLE traffic_events (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL REFERENCES tenants(id),
    connection_id UUID NOT NULL REFERENCES connections(id),
    event_type    traffic_event_type NOT NULL,
    data          JSONB NOT NULL DEFAULT '{}',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TYPE export_type AS ENUM ('audit', 'compliance', 'rgpd');
CREATE TYPE export_status AS ENUM ('pending', 'processing', 'completed', 'failed');

CREATE TABLE compliance_exports (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID NOT NULL REFERENCES tenants(id),
    type         export_type NOT NULL,
    date_range   DATERANGE NOT NULL,
    status       export_status NOT NULL DEFAULT 'pending',
    file_path    TEXT,
    requested_by UUID NOT NULL REFERENCES users(id),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Domain 7: Billing & Licensing
-- ----------------------------------------------------------------------------

CREATE TYPE plan_tier AS ENUM ('free', 'pro', 'enterprise');
CREATE TYPE subscription_status AS ENUM ('active', 'past_due', 'canceled', 'expired', 'trialing');

CREATE TABLE subscriptions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID NOT NULL REFERENCES tenants(id),
    plan                plan_tier NOT NULL DEFAULT 'free',
    status              subscription_status NOT NULL DEFAULT 'active',
    current_period_start TIMESTAMPTZ NOT NULL,
    current_period_end  TIMESTAMPTZ NOT NULL,
    stripe_sub_id       TEXT,
    stripe_customer_id  TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE licenses (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    license_key TEXT NOT NULL UNIQUE,
    tier        plan_tier NOT NULL,
    seats       INTEGER NOT NULL DEFAULT 1,
    features    JSONB NOT NULL DEFAULT '{}',
    issued_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Domain 8: Feature Flags & Config
-- ----------------------------------------------------------------------------

CREATE TABLE tenant_config (
    tenant_id               UUID PRIMARY KEY REFERENCES tenants(id),
    theme                   TEXT NOT NULL DEFAULT 'dark',
    session_timeout_minutes INTEGER NOT NULL DEFAULT 480,
    max_concurrent_conns    INTEGER NOT NULL DEFAULT 5,
    require_mfa             BOOLEAN NOT NULL DEFAULT false,
    allowed_protocols       TEXT[] DEFAULT '{ssh,rdp,vnc}',
    audit_retention_days    INTEGER NOT NULL DEFAULT 365,
    custom_domain           TEXT,
    logo_url                TEXT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE feature_flags (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id),
    feature_key TEXT NOT NULL,
    enabled     BOOLEAN NOT NULL DEFAULT false,
    overrides   JSONB DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, feature_key)
);

CREATE TABLE webhooks (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  UUID NOT NULL REFERENCES tenants(id),
    url        TEXT NOT NULL,
    events     TEXT[] NOT NULL,
    secret     TEXT NOT NULL,
    enabled    BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
