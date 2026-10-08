---
page_title: "oneuptime_network_alert_policy Data Source - oneuptime"
subcategory: "Other"
description: |-
  Alert on a set of network devices at once: every device matching the policy's sites, roles and labels gets a Network Device monitor provisioned from the policy's monitor template, and kept as devices come and go.
---

# oneuptime_network_alert_policy (Data Source)

Alert on a set of network devices at once: every device matching the policy's sites, roles and labels gets a Network Device monitor provisioned from the policy's monitor template, and kept as devices come and go.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one network alert policy may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_network_alert_policy" "example" {
  name = "Example network alert policy"
}

# Or by id:
data "oneuptime_network_alert_policy" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `covered_device_count` (Number) How many devices matched this policy's scope at the engine's last reconciliation. Managed by the engine.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this policy is active. Disable it to stop provisioning monitors for matching devices without deleting the policy.
- `last_sync_error` (String) Why the engine's last reconciliation of this policy failed, if it did. Cleared by the next successful pass. Managed by the engine.
- `monitor_template_id` (String) ID of the Network Device monitor template every matching device gets a monitor cloned from. Null only after the template was deleted, which disables the policy. The ID of a `oneuptime_monitor_template`.
- `name` (String) Any friendly name of this object.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_sync_at` (String) When the engine last reconciled this policy's monitors against its matching devices. Managed by the engine.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `scope` (String) Which devices this policy covers: site ids, device role ids and label ids. A device must match every kind that is listed (AND) and any id within a kind (OR); a kind left empty matches every device. Empty altogether means every device in the project. A JSON value: write it with `jsonencode()`.
- `template_synced_at` (String) When this policy's provisioned monitors were last re-synced from the monitor template. Managed by the engine.
- `updated_at` (String) Date and Time when the object was updated.
