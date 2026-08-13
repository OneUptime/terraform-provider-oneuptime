---
page_title: "oneuptime_status_page_resource Resource - oneuptime"
subcategory: "Status Pages"
description: |-
  Add resources like monitors to your status page
---

# oneuptime_status_page_resource (Resource)

Add resources like monitors to your status page

## Example Usage

```terraform
resource "oneuptime_status_page_resource" "example" {
  status_page_id = "123e4567-e89b-12d3-a456-426614174000"
  display_name = "Example short text"
}
```

## Schema

### Required

- `status_page_id` (String) A unique identifier for an object, represented as a UUID..
- `display_name` (String) Display name of the monitor on the Status Page..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `monitor_id` (String) A unique identifier for an object, represented as a UUID..
- `monitor_group_id` (String) A unique identifier for an object, represented as a UUID..
- `status_page_group_id` (String) A unique identifier for an object, represented as a UUID..
- `display_description` (String) Display description of the monitor on the Status Page. This is in markdown format...
- `display_tooltip` (String) Tooltip of the monitor on the Status Page..
- `show_current_status` (Bool) Show current status like offline, operational or degraded...
- `show_uptime_percent` (Bool) Show uptime percent of this monitor for the last 90 days..
- `uptime_percent_precision` (String) Precision of uptime percent of this monitor for the last 90 days..
- `show_status_history_chart` (Bool) Show a 90 day uptime history of this monitor..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `order` (Number) Order / Priority of this resource..
- `row_axis_value` (String) Row this resource belongs to when its status page group is rendered as a grid. Should match one of the row axis values defined on the group...
- `column_axis_value` (String) Column this resource belongs to when its status page group is rendered as a grid. Should match one of the column axis values defined on the group...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `status_page_monitor_rule_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_status_page_resource.example <id>
```
