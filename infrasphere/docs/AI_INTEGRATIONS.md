# AI Integrations

InfraSphere supports OpenAI, Anthropic Claude, Azure OpenAI, Google Vertex AI, AWS Bedrock, local Ollama, and internal enterprise LLM gateways.

## Architecture

The assistant retrieves inventory, topology, metrics, logs, traces, costs, alerts, deployments, and audit records. It builds a constrained context package, calls the selected model, and returns an answer with citations and risk classification.

## Prompting Strategy

Prompts should include tenant, selected context, available data sources, guardrails, and the exact user question. The assistant must state uncertainty when data is absent.

## Tool Calling

Tools are read-only by default: search inventory, fetch resource details, query cost, query alerts, fetch deployment history, and draft a plan. Mutating tools require approval records and an auditable diff.

## Guardrails

AI can recommend, explain, summarize, generate Terraform/Kubernetes YAML, and draft deployment plans. AI cannot directly delete, scale down production, modify IAM, change firewall rules, destroy Terraform resources, restart production clusters, or alter storage without explicit approval.

## RAG Over Inventory

Inventory resources and relationships should be embedded or indexed by provider, account, region, resource type, tags, owner, workload, and environment. Source IDs must be cited in answers.

## Incident Triage

Triage combines alerts, logs, metrics, deployment events, topology, recent changes, and resource health to propose likely cause, blast radius, remediation, and postmortem draft.

## Example Prompts

- Show all underutilized VMs and estimate monthly savings.
- Which workloads are exposed to the internet?
- Explain this Terraform plan before I apply it.
- Find Kubernetes clusters running old versions.
- Summarize production incidents from the last 24 hours.

