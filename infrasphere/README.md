# InfraSphere Control Plane

InfraSphere is an enterprise infrastructure control plane for discovering, managing, observing, and deploying workloads across AWS, Google Cloud, Azure, VMware Cloud Foundation, and future private cloud or edge environments.

It is designed as the single operational pane for platform engineering, SRE, cloud, FinOps, and security teams: unified inventory, provider-neutral deployments, observability correlation, cost intelligence, security posture, and AI-assisted infrastructure reasoning.

## Key Features

- Secure login, JWT sessions, RBAC roles, and multi-tenant organization model.
- Cloud account registry for AWS, GCP, Azure, VCF, and mock provider environments.
- Unified inventory model for compute, network, storage, Kubernetes, database, identity, security, cost, and observability resources.
- Deployment wizard for containers, Helm, Terraform/OpenTofu, Kubernetes manifests, VMs, serverless, batch, and AI workloads.
- AI assistant UI with read-only analysis, cited data sources, and approval-oriented guardrails.
- Observability, cost, and security dashboards with extensible provider abstractions.
- PostgreSQL schema, Docker Compose, Helm, Terraform starter modules, and production-oriented documentation.

## Screenshots

Screenshots live here once captured from the local app:

- `docs/images/dashboard.png`
- `docs/images/inventory.png`
- `docs/images/ai-assistant.png`

## Quick Start

```bash
cd infrasphere
docker compose up --build
```

Open:

```text
http://localhost:5173
```

Default admin:

```text
email: admin@infrasphere.local
password: ChangeMe123!
```

## Architecture Summary

- Backend: Go REST API, provider abstractions, in-memory MVP service layer, PostgreSQL migrations, worker entrypoint.
- Frontend: React, TypeScript, Vite, console-style UI, dashboard, inventory, deployment, AI, cost, security, and observability pages.
- Infra: Docker Compose for local development, Helm chart starter, Terraform cloud deployment placeholders.

## Supported Providers

- Mock provider: working local demo data.
- AWS: starter adapter that shares the provider interface and is ready for SDK-backed discovery.
- GCP, Azure, VCF: starter adapters and documentation for credential models and inventory mapping.

## Roadmap

MVP focuses on auth, organizations, cloud connections, mock inventory, AWS starter adapter, deployment planning UI, AI UI, observability abstraction, cost/security skeletons, and docs. Enterprise phases add persistent repositories, real SDK integrations, OIDC, encrypted credentials, policy engine, graph queries, Terraform/OpenTofu execution, live Kubernetes deployment, and governed AI automation.

