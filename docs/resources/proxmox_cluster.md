---
page_title: "oneuptime_proxmox_cluster Resource - oneuptime"
subcategory: "Other"
description: |-
  Proxmox VE clusters that are being monitored in this project. Each cluster is auto-discovered when the OneUptime Proxmox Agent sends metrics, or can be manually registered.
---

# oneuptime_proxmox_cluster (Resource)

Proxmox VE clusters that are being monitored in this project. Each cluster is auto-discovered when the OneUptime Proxmox Agent sends metrics, or can be manually registered.

## Example Usage

```terraform
resource "oneuptime_proxmox_cluster" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this Proxmox cluster. This is the join key — it must match the proxmox.cluster.name OTel resource attribute stamped by the OneUptime Proxmox Agent...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description for this Proxmox cluster..
- `ceph_cluster_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_archived` (Bool) Is this Proxmox cluster archived? Archived Proxmox clusters are hidden from lists but keep collecting telemetry...
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this Proxmox cluster. Leave blank to use the project-wide default...
- `telemetry_retention_config` (String) Per-pillar retention overrides for this Proxmox cluster (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the Proxmox cluster default, then the project's retention settings...
- `otel_collector_status` (String) Connection status of the OTel Collector agent (connected or disconnected)..
- `agent_version` (String) Version of the OneUptime Proxmox agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute..
- `pve_version` (String) Proxmox VE version reported by this cluster..
- `last_seen_at` (String) A date time object..
- `node_count` (Number) Cached count of nodes in this cluster..
- `online_node_count` (Number) Cached count of nodes currently online (pve_up == 1) in this cluster. Rendered as 'Nodes X/Y online' next to nodeCount...
- `guest_count` (Number) Cached count of guests (VMs and containers) in this cluster..
- `storage_count` (Number) Cached count of storage pools in this cluster..
- `guests_without_backup_count` (Number) Cached count of guests not covered by ANY backup job (pve_not_backed_up_total). NULL until the exporter's cluster-level backup-info collector reports. Coverage by a job is NOT the same as recent/successful backups — freshness needs the PVE task log or PBS API...

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
terraform import oneuptime_proxmox_cluster.example <id>
```
