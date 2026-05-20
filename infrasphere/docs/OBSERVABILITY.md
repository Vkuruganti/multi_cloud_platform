# Observability

Supported targets include Prometheus, Grafana, OpenTelemetry Collector, Datadog, New Relic, Splunk, Elastic, CloudWatch, Azure Monitor, Google Cloud Monitoring, and VMware Aria Operations.

## Metrics Model

Metrics queries include resource ID, metric name, start time, end time, dimensions, and aggregation. Results are normalized into timestamp/value series.

## Logs Model

Logs queries include resource, workload, namespace, severity, time range, and text filter. Logs are redacted before AI use.

## Traces Model

Trace queries include service, operation, trace ID, duration, status, and time range. Traces power deployment-to-incident correlation.

## Alerts

Alerts contain severity, source, title, status, start time, resolved time, and provider metadata.

## SLO Dashboard

SLOs should track availability, latency, error budget burn, saturation, and deployment correlation by workload and environment.

## OpenTelemetry

OTel Collector is the preferred ingestion point for app metrics, traces, and logs. InfraSphere should link workloads to OTel service names and resource attributes.

