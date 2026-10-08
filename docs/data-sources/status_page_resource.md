---
page_title: "oneuptime_status_page_resource Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Add resources like monitors to your status page
---

# oneuptime_status_page_resource (Data Source)

Add resources like monitors to your status page

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one status page resource may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_status_page_resource" "example" {
  status_page_id = oneuptime_status_page.example.id
}

# Or by id:
data "oneuptime_status_page_resource" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `column_axis_value` (String) Column this resource belongs to when its status page group is rendered as a grid. Should match one of the column axis values defined on the group.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `display_description` (String) Display description of the monitor on the Status Page. This is in markdown format.
- `display_name` (String) Display name of the monitor on the Status Page.
- `display_tooltip` (String) Tooltip of the monitor on the Status Page.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `monitor_group_id` (String) Relation to Monitor Group ID Resource in which this object belongs. The ID of a `oneuptime_monitor_group`.
- `monitor_id` (String) Relation to Monitor ID Resource in which this object belongs. The ID of a `oneuptime_monitor`.
- `order` (Number) Order / Priority of this resource.
- `row_axis_value` (String) Row this resource belongs to when its status page group is rendered as a grid. Should match one of the row axis values defined on the group.
- `show_current_status` (Boolean) Show current status like offline, operational or degraded.
- `show_status_history_chart` (Boolean) Show a 90 day uptime history of this monitor.
- `show_uptime_percent` (Boolean) Show uptime percent of this monitor for the last 90 days.
- `status_page_group_id` (String) Does this monitor belong to a status page group? The ID of a `oneuptime_status_page_group`.
- `status_page_id` (String) ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.
- `status_page_monitor_rule_id` (String) ID of the rule that added this resource, if it was added by a rule instead of by hand. The ID of a `oneuptime_status_page_monitor_rule`.
- `uptime_percent_precision` (String) Precision of uptime percent of this monitor for the last 90 days.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
