---
page_title: "oneuptime_scheduled_maintenance_reminder_rule Resource - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Configure reminder rules to periodically notify maintenance event owners while an event is still not complete
---

# oneuptime_scheduled_maintenance_reminder_rule (Resource)

Configure reminder rules to periodically notify maintenance event owners while an event is still not complete

## Example Usage

```terraform
resource "oneuptime_scheduled_maintenance_reminder_rule" "example" {
  name                         = "Example scheduled maintenance reminder rule"
  reminder_interval_in_minutes = 42
  description                  = "Managed by Terraform"
}
```

<a id="order"></a>
## Order

`order` is this scheduled maintenance reminder rule's place among the scheduled maintenance reminder rules of its project, lowest first. Setting it places this scheduled maintenance reminder rule and moves no other, so a configuration that numbers every scheduled maintenance reminder rule of a project is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.scheduled_maintenance_reminder_rules`:

```terraform
order = index(local.scheduled_maintenance_reminder_rules, each.key) + 1
```

- Give each scheduled maintenance reminder rule a number of its own: two scheduled maintenance reminder rules at one number are shown the older first.
- A scheduled maintenance reminder rule created without `order` goes after the last one of its project.
- Deleting a scheduled maintenance reminder rule leaves a gap in the numbers, which changes nothing: scheduled maintenance reminder rules go by their numbers, not by counting.
- Moving a scheduled maintenance reminder rule in the OneUptime Dashboard numbers the project's scheduled maintenance reminder rules 1..n again; the next plan shows that as changes to `order`.


## Schema

### Required

- `name` (String) Name of this reminder rule.
- `reminder_interval_in_minutes` (Number) How often (in minutes) to remind maintenance event owners while the event is still not complete. For example, set to 30 to remind owners every 30 minutes.

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this reminder rule.
- `is_enabled` (Boolean) Whether this reminder rule is enabled. Defaults to `true`.
- `labels` (Set of String) Only apply this reminder rule to maintenance events with these labels. Leave empty to match all events. IDs of `oneuptime_label` resources.
- `order` (Number) Where this rule sits in the list. Rules are evaluated from the top of the list down, lowest number first, and the first one that matches wins. A new rule is added to the end of the list. Setting a number places it there and moves no other one: two at one number are read the older first. In the dashboard, drag the rows to reorder them. See [Order](#order).
- `remind_while_scheduled` (Boolean) Send reminders while the event is still scheduled (before it starts). When disabled, reminders only begin once the event has started. Defaults to `false`.
- `stop_reminders_on_state` (String) Stop sending reminders once the maintenance event reaches this state. Select Ongoing to stop reminders when the event starts, or Completed to keep reminding until the event is completed. Ongoing needs Remind While Event is Scheduled on: a rule that waits for the event to start and stops when it starts would never remind, so it is refused. Allowed values: `Ongoing`, `Completed`. Defaults to `Completed`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `emptied_limits` (String) The lists this rule was limited to whose every record was deleted since: each list, the deleted records' names, how many and when. Deleting the last record a list held would make the rule match every record instead, so OneUptime switched the rule off and wrote this. It is cleared when the rule is switched on again; until a list named here holds a record again, switching the rule on must say what the list holds in the same update. Written by OneUptime; any value sent for it is refused. A JSON value: read it with `jsondecode()`.
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing scheduled maintenance reminder rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_scheduled_maintenance_reminder_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_scheduled_maintenance_reminder_rule.example <id>
```
