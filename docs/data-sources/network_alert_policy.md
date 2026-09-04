---
page_title: "oneuptime_network_alert_policy Data Source - oneuptime"
subcategory: "Other"
description: |-
  Alert on a set of network devices at once: every device matching the policy's sites, roles and labels gets a Network Device monitor provisioned from the policy's monitor template, and kept as devices come and go.
---

# oneuptime_network_alert_policy (Data Source)

Alert on a set of network devices at once: every device matching the policy's sites, roles and labels gets a Network Device monitor provisioned from the policy's monitor template, and kept as devices come and go. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_network_alert_policy" "by_name" {
  name = "example-network_alert_policy"
}

data "oneuptime_network_alert_policy" "by_id" {
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
- `description` (String) Friendly description that will help you remember.. Computed.
- `is_enabled` (Bool) Whether this policy is active. Disable it to stop provisioning monitors for matching devices without deleting the policy... Computed.
- `monitor_template_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `scope` (String) Which devices this policy covers: site ids, device role ids and label ids. A device must match every kind that is listed (AND) and any id within a kind (OR); a kind left empty matches every device. Empty altogether means every device in the project... Computed.
- `last_sync_at` (String) A date time object.. Computed.
- `last_sync_error` (String) Why the engine's last reconciliation of this policy failed, if it did. Cleared by the next successful pass. Managed by the engine... Computed.
- `covered_device_count` (Number) How many devices matched this policy's scope at the engine's last reconciliation. Managed by the engine... Computed.
- `template_synced_at` (String) A date time object.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
