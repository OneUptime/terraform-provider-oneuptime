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
  name = "Example short text"
  reminder_interval_in_minutes = 42
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this reminder rule..
- `reminder_interval_in_minutes` (Number) How often (in minutes) to remind alert owners while the alert is still open. For example, set to 30 to remind owners every 30 minutes...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this reminder rule..
- `order` (Number) Order/priority of this rule. Rules are evaluated in order (lowest first). First matching rule wins...
- `is_enabled` (Bool) Whether this reminder rule is enabled..
- `stop_reminders_on_state` (String) Stop sending reminders once the alert reaches this state. Select Acknowledged to stop reminders when the alert is acknowledged, or Resolved to keep reminding until the alert is resolved...
- `alert_severities` (Set) Only apply this reminder rule to alerts with these severities. Leave empty to match alerts of any severity...
- `labels` (Set) Only apply this reminder rule to alerts with these labels. Leave empty to match alerts with any labels...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_alert_reminder_rule.example <id>
```
