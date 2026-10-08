---
page_title: "oneuptime_status_page_group Resource - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage groups on your status page and categorize resources like monitors into these groups.
---

# oneuptime_status_page_group (Resource)

Manage groups on your status page and categorize resources like monitors into these groups.

## Example Usage

```terraform
resource "oneuptime_status_page_group" "example" {
  status_page_id = oneuptime_status_page.example.id
  name           = "Example status page group"
  description    = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of the Group.
- `status_page_id` (String) ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.

### Optional

- `column_axis_label` (String) Label shown above the column axis when the group is rendered as a grid (e.g. 'Region', 'Environment'). Free-form so you can use any dimension you like.
- `column_axis_values` (String) Comma-separated list of column labels for the grid (e.g. 'US-East, EU-West, Asia'). Determines column order in the grid layout.
- `description` (String) Description for this group. This is visible on Status Page. This can be in markdown format.
- `is_expanded_by_default` (Boolean) Is this group expanded by default on Status Page? Defaults to `true`.
- `order` (Number) Order / Priority of this resource.
- `parent_status_page_group_id` (String) ID of the Status Page Group this group is nested under. Empty for top level groups. The ID of a `oneuptime_status_page_group`.
- `row_axis_label` (String) Label shown above the row axis when the group is rendered as a grid (e.g. 'Service', 'Tenant'). Free-form so you can use any dimension you like.
- `row_axis_values` (String) Comma-separated list of row labels for the grid (e.g. 'Auth, API, Database'). Determines row order in the grid layout.
- `show_current_status` (Boolean) Show current status like offline, operational or degraded. Defaults to `true`.
- `show_uptime_percent` (Boolean) Show uptime percent of this group for the last 90 days. Defaults to `false`.
- `uptime_percent_precision` (String) Precision of uptime percent of this group for the last 90 days.
- `view_mode` (String) Layout of this group on the status page. 'List' renders resources stacked vertically (default). 'Grid' renders resources as a matrix using row and column axes. Defaults to `List`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing status page group by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_status_page_group.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_status_page_group.example <id>
```
