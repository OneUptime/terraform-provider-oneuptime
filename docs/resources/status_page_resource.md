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
  status_page_id = oneuptime_status_page.example.id
  display_name   = "Example short text"
}
```

## Schema

### Required

- `display_name` (String) Display name of the monitor on the Status Page.
- `status_page_id` (String) ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.

### Optional

- `column_axis_value` (String) Column this resource belongs to when its status page group is rendered as a grid. Should match one of the column axis values defined on the group.
- `display_description` (String) Display description of the monitor on the Status Page. This is in markdown format.
- `display_tooltip` (String) Tooltip of the monitor on the Status Page.
- `monitor_group_id` (String) Relation to Monitor Group ID Resource in which this object belongs. The ID of a `oneuptime_monitor_group`.
- `monitor_id` (String) Relation to Monitor ID Resource in which this object belongs. The ID of a `oneuptime_monitor`.
- `order` (Number) Order / Priority of this resource.
- `row_axis_value` (String) Row this resource belongs to when its status page group is rendered as a grid. Should match one of the row axis values defined on the group.
- `show_current_status` (Boolean) Show current status like offline, operational or degraded. Defaults to `true`.
- `show_status_history_chart` (Boolean) Show a 90 day uptime history of this monitor. Defaults to `true`.
- `show_uptime_percent` (Boolean) Show uptime percent of this monitor for the last 90 days. Defaults to `false`.
- `status_page_group_id` (String) Does this monitor belong to a status page group? The ID of a `oneuptime_status_page_group`.
- `uptime_percent_precision` (String) Precision of uptime percent of this monitor for the last 90 days.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `status_page_monitor_rule_id` (String) ID of the rule that added this resource, if it was added by a rule instead of by hand. The ID of a `oneuptime_status_page_monitor_rule`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing status page resource by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_status_page_resource.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_status_page_resource.example <id>
```
