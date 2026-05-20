# API Reference

Base URL: `http://localhost:8080`

## Auth

- `POST /api/auth/login` returns JWT token.
- `POST /api/auth/logout` returns stateless logout acknowledgement.
- `POST /api/auth/refresh` returns a refreshed JWT token.
- `GET /api/me` returns current user, roles, and organizations.

## Organizations

- `GET /api/orgs`
- `POST /api/orgs` planned.
- `GET /api/orgs/:id` planned.

## Providers and Cloud Accounts

- `GET /api/providers`
- `POST /api/cloud-accounts`
- `GET /api/cloud-accounts`
- `GET /api/cloud-accounts/:id` planned.
- `POST /api/cloud-accounts/:id/validate`
- `POST /api/cloud-accounts/:id/sync`

## Inventory

- `GET /api/inventory/resources`
- `GET /api/inventory/resources/:id`
- `GET /api/inventory/relationships`
- `GET /api/inventory/topology`

## Deployments

- `POST /api/deployments`
- `GET /api/deployments`
- `GET /api/deployments/:id` planned.
- `POST /api/deployments/:id/plan`
- `POST /api/deployments/:id/approve`
- `POST /api/deployments/:id/apply`
- `POST /api/deployments/:id/rollback`

## Observability

- `GET /api/observability/metrics`
- `GET /api/observability/logs`
- `GET /api/observability/traces`
- `GET /api/observability/alerts`

## AI

- `POST /api/ai/chat`
- `POST /api/ai/recommend`
- `POST /api/ai/explain-resource`
- `POST /api/ai/triage-incident`

## Cost

- `GET /api/cost/summary`
- `GET /api/cost/resources`
- `GET /api/cost/anomalies`
- `GET /api/cost/recommendations`

## Security and Audit

- `GET /api/security/findings`
- `GET /api/audit-logs`
