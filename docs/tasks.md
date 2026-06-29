# Cosmos — Tasks

## Palier 0 — Setup & Fondations
- [ ] Init repo cosmos-dashboard (fork netbird/dashboard)
- [ ] Init repo cosmos-server (Go)
- [ ] docker-compose.yml (cosmos-server + dashboard + guacd + postgres)
- [ ] CI/CD basique (build, lint, test)

## Palier 1 — Dashboard Core
- [ ] Remplacer les couleurs Netbird (orange) → Cosmos (bleu)
- [ ] Ajouter le light theme
- [ ] Remplacer le logo par le dragon Cosmos
- [ ] Page Dashboard (overview)
- [ ] Sidebar + navigation (servers, connections, users, audit, policies, settings)
- [ ] Auth pages (login, MFA, SSO)

## Palier 2 — Backend Core
- [ ] API Go: health, auth (JWT), tenant CRUD
- [ ] Migrations PostgreSQL (toutes les tables)
- [ ] Module audit_logs (capture événements API)
- [ ] Module servers (CRUD)
- [ ] Module users (CRUD + invitations)

## Palier 3 — Connexions
- [ ] Intégration guacd (client Go → guacd)
- [ ] Module connections (SSH, RDP, VNC)
- [ ] WebSocket proxy (navigateur → cosmos-server → guacd)
- [ ] Terminal SSH dans le navigateur
- [ ] RDP/VNC dans le navigateur

## Palier 4 — Zero Trust
- [ ] Module policies (CRUD + évaluation)
- [ ] Moteur de règles (conditions → allow/deny)
- [ ] MFA enforcement par policy
- [ ] Session recording

## Palier 5 — SaaS & Enterprise
- [ ] Module license (self-hosted)
- [ ] Module billing (Stripe pour SaaS)
- [ ] Feature flags (free vs pro vs enterprise)
- [ ] Compliance exports (RGPD/DORA)
- [ ] Webhooks
