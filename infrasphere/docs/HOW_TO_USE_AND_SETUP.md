# How To Use And Setup InfraSphere

This guide is the operational quick path for running, connecting, and extending InfraSphere.

## Run Locally

```bash
cd infrasphere
docker compose up --build
```

Open `http://localhost:5173`.

## Create An Admin User

For the MVP, use the seeded admin:

```text
email: admin@infrasphere.local
password: ChangeMe123!
```

For production, create a bootstrap command that inserts a user with a bcrypt or Argon2 password hash, then assigns the Platform Admin role in `organization_members`.

## Connect AWS

1. Create a cross-account IAM role with external ID.
2. Grant read-only discovery for EC2, VPC, ELB, EBS, S3, EKS, RDS, IAM, CloudWatch, and Cost Explorer.
3. In InfraSphere, open Cloud Connections.
4. Add provider `aws`, account ID, default region, role ARN, and external ID.
5. Validate credentials.
6. Run Sync Inventory.

## Connect GCP

1. Create a service account or workload identity federation binding.
2. Grant Compute Viewer, Kubernetes Engine Viewer, Cloud SQL Viewer, Monitoring Viewer, IAM Security Reviewer, and billing export access.
3. Add provider `gcp`, project ID, default region, and credential reference.
4. Validate and sync.

## Connect Azure

1. Create an Entra ID app registration or workload identity.
2. Assign Reader, Monitoring Reader, Cost Management Reader, and relevant resource readers.
3. Add provider `azure`, tenant ID, subscription ID, client ID, secret reference, and default region.
4. Validate and sync.

## Connect VCF

1. Create read-only service accounts for vCenter, NSX, SDDC Manager, and Aria Operations.
2. Add provider `vcf`, vCenter URL, NSX URL, SDDC Manager URL, and credential references.
3. Validate TLS and credentials.
4. Sync inventory.

## Configure AI Providers

Set one or more:

```bash
OPENAI_API_KEY=
ANTHROPIC_API_KEY=
AZURE_OPENAI_ENDPOINT=
AZURE_OPENAI_API_KEY=
GOOGLE_VERTEX_PROJECT=
AWS_REGION=
OLLAMA_BASE_URL=http://localhost:11434
```

AI is read-only by default. Any mutating action must create an approval request and show a diff.

## Configure Observability Tools

Add integrations for Prometheus, Grafana, Datadog, New Relic, Splunk, Elastic, CloudWatch, Azure Monitor, GCP Monitoring, or Aria Operations. Store tokens as encrypted config. Map workloads to service names and resource IDs.

## Sync Inventory

Use Cloud Connections and click sync, or call:

```bash
curl -X POST http://localhost:8080/api/cloud-accounts/acct-mock/sync \
  -H "Authorization: Bearer <token>"
```

Sync populates resources, relationships, cost snapshots, health state, and security findings.

## Deploy A Workload

1. Open Deploy.
2. Select provider, account, region, workload type, networking, storage, scaling, observability, and budget.
3. Review estimated cost.
4. Create a draft deployment.
5. Plan.
6. Approve.
7. Apply.
8. Watch events, logs, status, cost, and observability links.

See `examples-workload.yaml` for a provider-neutral workload definition.

## Use The AI Assistant

Ask questions such as:

- Show all underutilized VMs.
- Which workloads are exposed to the internet?
- Why did the deployment fail?
- Which resources are driving cost this month?
- Recommend a cheaper placement for this workload.

The assistant should cite inventory, logs, metrics, cost, alerts, or deployment records.

## Troubleshooting

- Frontend cannot reach API: confirm `backend` is healthy and `VITE_API_URL=http://localhost:8080`.
- Login fails: use the seeded admin exactly as shown.
- PostgreSQL migration fails: ensure Docker volume is fresh or remove the old volume after backing up.
- Port conflict: change Compose ports for `5173`, `8080`, `5432`, or `6379`.
- Cloud validation fails: check IAM/service account permissions, region, external ID, network access, and clock skew.

## Deploy To Kubernetes

1. Build and push backend/frontend images.
2. Create PostgreSQL and Redis services or use managed services.
3. Create Kubernetes Secrets for database URL, JWT secret, encryption key, AI keys, and provider secrets.
4. Install with Helm:

```bash
helm install infrasphere infra/helm/infrasphere
```

5. Add Ingress, TLS, HPA, NetworkPolicy, and backup automation.

## Extend With A New Provider

1. Add `backend/internal/providers/<provider>`.
2. Implement `CloudProvider`.
3. Normalize inventory into `models.Resource`.
4. Add resource relationships.
5. Add credential validation.
6. Add least-privilege docs.
7. Register provider in API/worker factories.
8. Add frontend provider badge and connection form fields.
9. Add tests with mock API responses.

