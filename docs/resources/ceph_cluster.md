---
page_title: "oneuptime_ceph_cluster Resource - oneuptime"
subcategory: "Other"
description: |-
  Ceph clusters that are being monitored in this project. Each cluster is auto-discovered when the OneUptime Ceph Agent sends metrics, or can be manually registered.
---

# oneuptime_ceph_cluster (Resource)

Ceph clusters that are being monitored in this project. Each cluster is auto-discovered when the OneUptime Ceph Agent sends metrics, or can be manually registered.

## Example Usage

```terraform
resource "oneuptime_ceph_cluster" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this Ceph cluster. This is the join key — it must match the ceph.cluster.name OTel resource attribute stamped by the OneUptime Ceph Agent...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description for this Ceph cluster..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_archived` (Bool) Is this Ceph cluster archived? Archived Ceph clusters are hidden from lists but keep collecting telemetry...
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this Ceph cluster. Leave blank to use the project-wide default...
- `telemetry_retention_config` (String) Per-pillar retention overrides for this Ceph cluster (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the Ceph cluster default, then the project's retention settings...
- `fsid` (String) Ceph cluster fsid, sourced from the ceph.cluster.fsid OTel resource attribute when known..
- `otel_collector_status` (String) Connection status of the OTel Collector agent (connected or disconnected)..
- `agent_version` (String) Version of the OneUptime Ceph agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute..
- `ceph_version` (String) Ceph version reported by this cluster..
- `last_seen_at` (String) A date time object..
- `mon_count` (Number) Cached count of Ceph monitors (mons) in this cluster..
- `osd_count` (Number) Cached count of OSDs in this cluster..
- `osd_up_count` (Number) Cached count of OSDs that are up (ceph_osd_up == 1) in this cluster. Rendered as 'X up / Y in / Z total' next to osdCount...
- `osd_in_count` (Number) Cached count of OSDs that are in the cluster (ceph_osd_in == 1). Rendered as 'X up / Y in / Z total' next to osdCount...
- `pool_count` (Number) Cached count of pools in this cluster..
- `health_status` (Number) Cached latest ceph_health_status value: 0 = HEALTH_OK, 1 = HEALTH_WARN, 2 = HEALTH_ERR. Rendered as the OK/WARN/ERR health pill. Null until the first metric batch arrives...
- `capacity_used_percent` (Number) Cached cluster capacity usage percent (ceph_cluster_total_used_bytes / ceph_cluster_total_bytes * 100). Stored as decimal so sub-percent precision survives the round trip. Null until both series appear in one metric batch...

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
terraform import oneuptime_ceph_cluster.example <id>
```
