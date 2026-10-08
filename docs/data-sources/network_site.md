---
page_title: "oneuptime_network_site Data Source - oneuptime"
subcategory: "Other"
description: |-
  Self-nesting sites (Account Type -> Region / Franchisee -> Market -> Unit) that group Network Devices into a drill-down hierarchy with a persisted health rollup.
---

# oneuptime_network_site (Data Source)

Self-nesting sites (Account Type -> Region / Franchisee -> Market -> Unit) that group Network Devices into a drill-down hierarchy with a persisted health rollup.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one network site may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_network_site" "example" {
  name = "Example network site"
}

# Or by id:
data "oneuptime_network_site" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `address` (String) Street address of this site, shown on map views.
- `alert_severity_id` (String) ID of the severity used for site-unhealthy alerts. The ID of a `oneuptime_alert_severity`.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `current_active_alert_id` (String) ID of the currently open site-unhealthy alert, if any. Managed by the rollup engine.
- `current_monitor_status_id` (String) Whats the current rolled-up status ID of this site? Computed from the devices and child sites below it. The ID of a `oneuptime_monitor_status`.
- `depth` (Number) Number of ancestors above this site (0 for root sites). Managed by the server on parent changes.
- `description` (String) Friendly description for this network site.
- `health_rollup_policy` (String) How this site's status is derived from the devices beneath it: WorstStatus (any device offline makes the site offline) or PercentThreshold (the share of devices that are down decides).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `latitude` (Number) Latitude of this site, for US and world map views.
- `longitude` (Number) Longitude of this site, for US and world map views.
- `materialized_path` (String) Slash-separated ancestor IDs of this site (e.g. '/rootId/childId/'). Managed by the server on parent changes; used for subtree queries and rollups.
- `name` (String) Friendly name for this network site.
- `network_site_type_id` (String) ID of the Network Site Type this site belongs to. The ID of a `oneuptime_network_site_type`.
- `offline_threshold_percent` (Number) With the PercentThreshold rollup policy: the share of reporting devices beneath this site that must be non-operational before the site itself is marked offline. Below it (but above zero) the site is degraded.
- `parent_site_id` (String) ID of the parent Network Site this site is nested under (empty for root sites). The ID of a `oneuptime_network_site`.
- `probe_id` (String) ID of the probe that polls devices in this site by default. The ID of a `oneuptime_probe`.
- `should_alert_when_unhealthy` (Boolean) When enabled, an alert opens when this site's health rollup turns non-operational and auto-resolves when it recovers.
- `site_type` (String) Deprecated legacy site type string. Use the Network Site Type relation instead; this column exists only for the backfill migration and will be removed.
- `slug` (String) Friendly globally unique name for your object.
- `snmp_credential_profile_id` (String) ID of the SNMP Credential Profile devices in this site inherit. The ID of a `oneuptime_snmp_credential_profile`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_rollup_at` (String) When the health rollup for this site was last computed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
