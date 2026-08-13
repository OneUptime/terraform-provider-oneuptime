---
page_title: "oneuptime_kubernetes_cluster Resource - oneuptime"
subcategory: "Other"
description: |-
  Kubernetes Clusters that are being monitored in this project. Each cluster is auto-discovered when the OneUptime kubernetes-agent sends metrics, or can be manually registered.
---

# oneuptime_kubernetes_cluster (Resource)

Kubernetes Clusters that are being monitored in this project. Each cluster is auto-discovered when the OneUptime kubernetes-agent sends metrics, or can be manually registered.

## Example Usage

```terraform
resource "oneuptime_kubernetes_cluster" "example" {
  name = "Example short text"
  cluster_identifier = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Friendly name for this Kubernetes cluster..
- `cluster_identifier` (String) Unique identifier for this cluster, sourced from the k8s.cluster.name OTel resource attribute..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description for this Kubernetes cluster..
- `provider` (String) Cloud provider or platform running this cluster (EKS, GKE, AKS, self-managed, unknown)..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_archived` (Bool) Is this Kubernetes cluster archived? Archived Kubernetes clusters are hidden from lists but keep collecting telemetry...
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this Kubernetes cluster. Leave blank to use the project-wide default...
- `telemetry_retention_config` (String) Per-pillar retention overrides for this Kubernetes cluster (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the cluster default, then the project's retention settings...
- `otel_collector_status` (String) Connection status of the OTel Collector agent (connected or disconnected)..
- `agent_version` (String) Version of the OneUptime Kubernetes agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute..
- `last_seen_at` (String) A date time object..
- `node_count` (Number) Cached count of nodes in this cluster..
- `pod_count` (Number) Cached count of pods in this cluster..
- `namespace_count` (Number) Cached count of namespaces in this cluster..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `archived_at` (String) A date time object..
- `archived_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_kubernetes_cluster.example <id>
```
