---
page_title: "oneuptime_alert_state_timeline Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Change state of the alerts (Created to Acknowledged for example)
---

# oneuptime_alert_state_timeline (Resource)

Change state of the alerts (Created to Acknowledged for example)

## Example Usage

```terraform
resource "oneuptime_alert_state_timeline" "example" {
  alert_id       = oneuptime_alert.example.id
  alert_state_id = oneuptime_alert_state.example.id
}
```

## Schema

### Required

- `alert_id` (String) Relation to Alert ID in which this resource belongs. The ID of a `oneuptime_alert`.
- `alert_state_id` (String) Alert State ID Relation. Which alert state does this alert change to? The ID of a `oneuptime_alert_state`.

### Optional

- `ends_at` (String) When did this status change end?
- `root_cause` (String) What is the root cause of this status change?
- `starts_at` (String) When did this status change? Correct this when the recorded time is wrong - every measurement derived from this timeline is recomputed from the corrected value.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of state change?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `state_change_log` (String) A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing alert state timeline by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_alert_state_timeline.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_alert_state_timeline.example <id>
```
