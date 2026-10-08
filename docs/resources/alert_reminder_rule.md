---
page_title: "oneuptime_alert_reminder_rule Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Configure reminder rules to periodically notify alert owners while an alert is still open
---

# oneuptime_alert_reminder_rule (Resource)

Configure reminder rules to periodically notify alert owners while an alert is still open

## Example Usage

```terraform
resource "oneuptime_alert_reminder_rule" "example" {
  name                         = "Example alert reminder rule"
  reminder_interval_in_minutes = 42
  description                  = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this reminder rule.
- `reminder_interval_in_minutes` (Number) How often (in minutes) to remind alert owners while the alert is still open. For example, set to 30 to remind owners every 30 minutes.

### Optional

- `alert_severities` (Set of String) Only apply this reminder rule to alerts with these severities. Leave empty to match alerts of any severity. IDs of `oneuptime_alert_severity` resources.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this reminder rule.
- `is_enabled` (Boolean) Whether this reminder rule is enabled. Defaults to `true`.
- `labels` (Set of String) Only apply this reminder rule to alerts with these labels. Leave empty to match alerts with any labels. IDs of `oneuptime_label` resources.
- `order` (Number) Where this rule sits in the list. Rules are evaluated from the top of the list down, lowest number first, and the first one that matches wins. A new rule is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.
- `stop_reminders_on_state` (String) Stop sending reminders once the alert reaches this state. Select Acknowledged to stop reminders when the alert is acknowledged, or Resolved to keep reminding until the alert is resolved. Defaults to `Resolved`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing alert reminder rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_alert_reminder_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_alert_reminder_rule.example <id>
```
