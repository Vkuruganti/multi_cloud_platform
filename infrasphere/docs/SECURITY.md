# Security

Security is a product surface, not an afterthought.

## Authentication

The MVP has JWT sessions and seeded login. Production should use bcrypt or Argon2 password hashes, refresh tokens, OIDC, SAML where required, MFA, and session revocation.

## Authorization

RBAC roles include Platform Admin, Cloud Admin, SRE, Developer, FinOps Analyst, Security Analyst, and Read-only Auditor. Every API must check organization membership and role permissions.

## Tenant Isolation

Every resource row is organization-scoped. Queries must always include organization ID from trusted auth context, never from a user-controlled body alone.

## Credential Encryption

Provider credentials belong in `provider_credentials.encrypted_payload`, encrypted with envelope encryption and KMS-backed key IDs. APIs must never return raw credentials.

## Audit Logging

Every login, provider connection, sync, AI action, deployment, approval, and destructive request should write an audit event.

## Approval Workflows

Require explicit approval for deleting resources, scaling down production, changing security groups, deleting storage, modifying IAM, destroying Terraform-managed resources, or restarting production clusters.

## Data Privacy

Redact secrets in logs, prompts, traces, and model context. Do not send raw credentials, tokens, passwords, or private keys to AI providers.

