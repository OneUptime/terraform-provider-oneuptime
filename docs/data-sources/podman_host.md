---
page_title: "oneuptime_podman_host Data Source - oneuptime"
subcategory: "Other"
description: |-
  Podman Hosts that are being monitored in this project. Each host is auto-discovered when the OneUptime Podman Agent sends metrics, or can be manually registered.
---

# oneuptime_podman_host (Data Source)

Podman Hosts that are being monitored in this project. Each host is auto-discovered when the OneUptime Podman Agent sends metrics, or can be manually registered. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_podman_host" "by_name" {
  name = "example-podman_host"
}

data "oneuptime_podman_host" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `description` (String) Friendly description for this Podman host.. Computed.
- `host_identifier` (String) Unique identifier for this Podman host, sourced from the host.name OTel resource attribute.. Computed.
- `otel_collector_status` (String) Connection status of the OTel Collector agent (connected or disconnected).. Computed.
- `agent_version` (String) Version of the OneUptime Podman agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute.. Computed.
- `last_seen_at` (String) A date time object.. Computed.
- `containers_running` (Number) Cached count of running containers on this host.. Computed.
- `containers_stopped` (Number) Cached count of stopped containers on this host.. Computed.
- `containers_paused` (Number) Cached count of paused containers on this host.. Computed.
- `os_type` (String) Operating system type of the Podman host.. Computed.
- `os_version` (String) Operating system version of the Podman host.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_archived` (Bool) Is this Podman host archived? Archived Podman hosts are hidden from lists but keep collecting telemetry... Computed.
- `archived_at` (String) A date time object.. Computed.
- `archived_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this Podman host. Leave blank to use the project-wide default... Computed.
- `telemetry_retention_config` (String) Per-pillar retention overrides for this Podman host (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the Podman host default, then the project's retention settings... Computed.
