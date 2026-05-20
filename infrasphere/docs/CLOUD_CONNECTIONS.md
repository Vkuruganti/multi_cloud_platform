# Cloud Connections

Cloud accounts are organization-scoped and credentials must be write-only. Raw secrets should never be returned by an API response.

## AWS

Use an IAM role with external ID, or access keys for development only. Minimum discovery permissions include EC2 describe, ELB describe, EBS describe, S3 list/get bucket metadata, EKS list/describe, RDS describe, IAM get/list, CloudWatch read, and Cost Explorer read.

Recommended production model: cross-account role with least privilege and CloudTrail audit enabled.

## GCP

Use a service account with project-level read roles for Compute, GKE, Cloud SQL, IAM policy viewer, Monitoring viewer, and Billing Account Costs Manager or billing export access. Prefer workload identity federation over static JSON keys.

## Azure

Use an Entra ID app registration with subscription-scoped Reader, Cost Management Reader, Monitoring Reader, and resource-specific roles for AKS, SQL, Storage, and Network discovery. Store client secrets encrypted, or use workload identity where available.

## VCF / vCenter / NSX / SDDC Manager

Use service accounts with read-only access to vCenter inventory, ESXi hosts, clusters, datastores, distributed switches, NSX segments/gateways, Tanzu Kubernetes clusters, SDDC Manager inventory, and Aria Operations metrics.

## Least Privilege

Start read-only for inventory, cost, security, and observability. Deployment permissions should be separated and protected by approval workflows.

