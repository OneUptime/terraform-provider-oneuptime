---
page_title: "oneuptime_network_site Resource - oneuptime"
subcategory: "Other"
description: |-
  Self-nesting sites (Account Type -> Region / Franchisee -> Market -> Unit) that group Network Devices into a drill-down hierarchy with a persisted health rollup.
---

# oneuptime_network_site (Resource)

Self-nesting sites (Account Type -> Region / Franchisee -> Market -> Unit) that group Network Devices into a drill-down hierarchy with a persisted health rollup.

## Example Usage

```terraform
resource "oneuptime_network_site" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Friendly name for this network site..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description for this network site..
- `site_type` (String) Deprecated legacy site type string. Use the Network Site Type relation instead; this column exists only for the backfill migration and will be removed...
- `network_site_type_id` (String) A unique identifier for an object, represented as a UUID..
- `parent_site_id` (String) A unique identifier for an object, represented as a UUID..
- `address` (String) Street address of this site, shown on map views..
- `latitude` (Number) Latitude of this site, for US and world map views..
- `longitude` (Number) Longitude of this site, for US and world map views..
- `health_rollup_policy` (String) How this site's status is derived from the devices beneath it: WorstStatus (any device offline makes the site offline) or PercentThreshold (the share of devices that are down decides)...
- `offline_threshold_percent` (Number) With the PercentThreshold rollup policy: the share of reporting devices beneath this site that must be non-operational before the site itself is marked offline. Below it (but above zero) the site is degraded...
- `should_alert_when_unhealthy` (Bool) When enabled, an alert opens when this site's health rollup turns non-operational and auto-resolves when it recovers...
- `alert_severity_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `materialized_path` (String) Slash-separated ancestor IDs of this site (e.g. '/rootId/childId/'). Managed by the server on parent changes; used for subtree queries and rollups...
- `depth` (Number) Number of ancestors above this site (0 for root sites). Managed by the server on parent changes...
- `current_monitor_status_id` (String) A unique identifier for an object, represented as a UUID..
- `last_rollup_at` (String) A date time object..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `current_active_alert_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_network_site.example <id>
```
