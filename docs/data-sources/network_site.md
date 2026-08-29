---
page_title: "oneuptime_network_site Data Source - oneuptime"
subcategory: "Other"
description: |-
  Self-nesting sites (Account Type -> Region / Franchisee -> Market -> Unit) that group Network Devices into a drill-down hierarchy with a persisted health rollup.
---

# oneuptime_network_site (Data Source)

Self-nesting sites (Account Type -> Region / Franchisee -> Market -> Unit) that group Network Devices into a drill-down hierarchy with a persisted health rollup. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_network_site" "by_name" {
  name = "example-network_site"
}

data "oneuptime_network_site" "by_id" {
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
- `description` (String) Friendly description for this network site.. Computed.
- `site_type` (String) Deprecated legacy site type string. Use the Network Site Type relation instead; this column exists only for the backfill migration and will be removed... Computed.
- `network_site_type_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `parent_site_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `materialized_path` (String) Slash-separated ancestor IDs of this site (e.g. '/rootId/childId/'). Managed by the server on parent changes; used for subtree queries and rollups... Computed.
- `depth` (Number) Number of ancestors above this site (0 for root sites). Managed by the server on parent changes... Computed.
- `address` (String) Street address of this site, shown on map views.. Computed.
- `latitude` (Number) Latitude of this site, for US and world map views.. Computed.
- `longitude` (Number) Longitude of this site, for US and world map views.. Computed.
- `current_monitor_status_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `last_rollup_at` (String) A date time object.. Computed.
- `health_rollup_policy` (String) How this site's status is derived from the devices beneath it: WorstStatus (any device offline makes the site offline) or PercentThreshold (the share of devices that are down decides)... Computed.
- `offline_threshold_percent` (Number) With the PercentThreshold rollup policy: the share of reporting devices beneath this site that must be non-operational before the site itself is marked offline. Below it (but above zero) the site is degraded... Computed.
- `should_alert_when_unhealthy` (Bool) When enabled, an alert opens when this site's health rollup turns non-operational and auto-resolves when it recovers... Computed.
- `alert_severity_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `current_active_alert_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
