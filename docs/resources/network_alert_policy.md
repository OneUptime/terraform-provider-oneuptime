---
page_title: "oneuptime_network_alert_policy Resource - oneuptime"
subcategory: "Other"
description: |-
  Alert on a set of network devices at once: every device matching the policy's sites, roles and labels gets a Network Device monitor provisioned from the policy's monitor template, and kept as devices come and go.
---

# oneuptime_network_alert_policy (Resource)

Alert on a set of network devices at once: every device matching the policy's sites, roles and labels gets a Network Device monitor provisioned from the policy's monitor template, and kept as devices come and go.

## Example Usage

```terraform
resource "oneuptime_network_alert_policy" "example" {
  name = "Example short text"
  monitor_template_id = "123e4567-e89b-12d3-a456-426614174000"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object..
- `monitor_template_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `is_enabled` (Bool) Whether this policy is active. Disable it to stop provisioning monitors for matching devices without deleting the policy...
- `scope` (String) Which devices this policy covers: site ids, device role ids and label ids. A device must match every kind that is listed (AND) and any id within a kind (OR); a kind left empty matches every device. Empty altogether means every device in the project...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `last_sync_at` (String) A date time object..
- `last_sync_error` (String) Why the engine's last reconciliation of this policy failed, if it did. Cleared by the next successful pass. Managed by the engine...
- `covered_device_count` (Number) How many devices matched this policy's scope at the engine's last reconciliation. Managed by the engine...
- `template_synced_at` (String) A date time object..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_network_alert_policy.example <id>
```
