---
page_title: "oneuptime_incident_reminder_rule Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Configure reminder rules to periodically notify incident owners while an incident is still open
---

# oneuptime_incident_reminder_rule (Data Source)

Configure reminder rules to periodically notify incident owners while an incident is still open

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident reminder rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_reminder_rule" "example" {
  name = "Example incident reminder rule"
}

# Or by id:
data "oneuptime_incident_reminder_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this reminder rule.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this reminder rule is enabled.
- `name` (String) Name of this reminder rule.
- `order` (Number) Where this rule sits in the list. Rules are evaluated from the top of the list down, lowest number first, and the first one that matches wins. A new rule is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.
- `reminder_interval_in_minutes` (Number) How often (in minutes) to remind incident owners while the incident is still open. For example, set to 30 to remind owners every 30 minutes.
- `stop_reminders_on_state` (String) Stop sending reminders once the incident reaches this state. Select Acknowledged to stop reminders when the incident is acknowledged, or Resolved to keep reminding until the incident is resolved.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `incident_severities` (Set of String) Only apply this reminder rule to incidents with these severities. Leave empty to match incidents of any severity. IDs of `oneuptime_incident_severity` resources.
- `labels` (Set of String) Only apply this reminder rule to incidents with these labels. Leave empty to match incidents with any labels. IDs of `oneuptime_label` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
