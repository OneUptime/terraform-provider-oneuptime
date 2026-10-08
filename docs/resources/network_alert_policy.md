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
  name                = "Example network alert policy"
  monitor_template_id = oneuptime_monitor_template.example.id
  description         = "Managed by Terraform"
}
```

## Schema

### Required

- `monitor_template_id` (String) ID of the Network Device monitor template every matching device gets a monitor cloned from. Null only after the template was deleted, which disables the policy. The ID of a `oneuptime_monitor_template`.
- `name` (String) Any friendly name of this object.

### Optional

- `description` (String) Friendly description that will help you remember.
- `is_enabled` (Boolean) Whether this policy is active. Disable it to stop provisioning monitors for matching devices without deleting the policy. Defaults to `true`.
- `scope` (String) Which devices this policy covers: site ids, device role ids and label ids. A device must match every kind that is listed (AND) and any id within a kind (OR); a kind left empty matches every device. Empty altogether means every device in the project. A JSON value: write it with `jsonencode()`. Defaults to `[object Object]`.

### Read-Only

- `covered_device_count` (Number) How many devices matched this policy's scope at the engine's last reconciliation. Managed by the engine.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_sync_at` (String) When the engine last reconciled this policy's monitors against its matching devices. Managed by the engine.
- `last_sync_error` (String) Why the engine's last reconciliation of this policy failed, if it did. Cleared by the next successful pass. Managed by the engine.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `template_synced_at` (String) When this policy's provisioned monitors were last re-synced from the monitor template. Managed by the engine.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing network alert policy by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_network_alert_policy.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_network_alert_policy.example <id>
```
