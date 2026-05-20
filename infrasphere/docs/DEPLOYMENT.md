# Deployment

## Docker Compose

```bash
cd infrasphere
docker compose up --build
```

Compose runs PostgreSQL, Redis, backend API, worker, and frontend.

## Kubernetes

Use `infra/helm/infrasphere` as the production packaging base. Add Deployments for API, worker, and frontend; Services; Ingress; Secrets; ConfigMaps; HPA; PodDisruptionBudgets; and NetworkPolicies.

## Helm Chart

```bash
helm install infrasphere infra/helm/infrasphere
```

The chart is intentionally minimal and ready for templates to be added as deployment targets are finalized.

## Production Hardening

Use managed PostgreSQL, managed Redis, KMS, TLS, OIDC, external secret management, private networking, image scanning, non-root containers, rate limiting, WAF, backups, and audit log retention.

## Backup Strategy

Back up PostgreSQL continuously with point-in-time recovery. Store audit logs in immutable storage. Export provider credentials only through encrypted KMS workflows.

## Scaling Workers

Scale workers horizontally by queue partitions: inventory sync, deployment execution, observability ingestion, cost ingestion, AI jobs, and security scans.

## TLS

Terminate TLS at ingress or load balancer. Enforce HTTPS, secure cookies, strict CORS, and HSTS in production.

