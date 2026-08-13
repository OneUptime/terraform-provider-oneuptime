---
page_title: "oneuptime_status_page_resource Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Add resources like monitors to your status page
---

# oneuptime_status_page_resource (Data Source)

Add resources like monitors to your status page Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_status_page_resource" "by_name" {
  name = "example-status_page_resource"
}

data "oneuptime_status_page_resource" "by_id" {
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
- `status_page_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `monitor_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `monitor_group_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `status_page_group_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `status_page_monitor_rule_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `display_name` (String) Display name of the monitor on the Status Page.. Computed.
- `display_description` (String) Display description of the monitor on the Status Page. This is in markdown format... Computed.
- `display_tooltip` (String) Tooltip of the monitor on the Status Page.. Computed.
- `show_current_status` (Bool) Show current status like offline, operational or degraded... Computed.
- `show_uptime_percent` (Bool) Show uptime percent of this monitor for the last 90 days.. Computed.
- `uptime_percent_precision` (String) Precision of uptime percent of this monitor for the last 90 days.. Computed.
- `show_status_history_chart` (Bool) Show a 90 day uptime history of this monitor.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `order` (Number) Order / Priority of this resource.. Computed.
- `row_axis_value` (String) Row this resource belongs to when its status page group is rendered as a grid. Should match one of the row axis values defined on the group... Computed.
- `column_axis_value` (String) Column this resource belongs to when its status page group is rendered as a grid. Should match one of the column axis values defined on the group... Computed.
