---
page_title: "oneuptime_incident_reminder_rule Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Configure reminder rules to periodically notify incident owners while an incident is still open
---

# oneuptime_incident_reminder_rule (Resource)

Configure reminder rules to periodically notify incident owners while an incident is still open

## Example Usage

```terraform
resource "oneuptime_incident_reminder_rule" "example" {
  name = "Example short text"
  reminder_interval_in_minutes = 42
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this reminder rule..
- `reminder_interval_in_minutes` (Number) How often (in minutes) to remind incident owners while the incident is still open. For example, set to 30 to remind owners every 30 minutes...

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource...
- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this reminder rule..
- `order` (Number) Where this rule sits in the list. Rules are evaluated from the top of the list down, lowest number first, and the first one that matches wins. A new rule is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them...
- `is_enabled` (Bool) Whether this reminder rule is enabled..
- `stop_reminders_on_state` (String) Stop sending reminders once the incident reaches this state. Select Acknowledged to stop reminders when the incident is acknowledged, or Resolved to keep reminding until the incident is resolved...
- `incident_severities` (Set) Only apply this reminder rule to incidents with these severities. Leave empty to match incidents of any severity...
- `labels` (Set) Only apply this reminder rule to incidents with these labels. Leave empty to match incidents with any labels...
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
terraform import oneuptime_incident_reminder_rule.example <id>
```
