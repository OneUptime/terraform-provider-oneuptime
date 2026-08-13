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
  status_page_id = "123e4567-e89b-12d3-a456-426614174000"
  name = "Example short text"
  description = "# Heading

This is **markdown** content"
}
```

## Schema

### Required

- `status_page_id` (String) A unique identifier for an object, represented as a UUID..
- `name` (String) Name of the Group..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `parent_status_page_group_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description for this group. This is visible on Status Page. This can be in markdown format...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `order` (Number) Order / Priority of this resource..
- `is_expanded_by_default` (Bool) Is this group expanded by default on Status Page?..
- `show_current_status` (Bool) Show current status like offline, operational or degraded...
- `show_uptime_percent` (Bool) Show uptime percent of this group for the last 90 days..
- `uptime_percent_precision` (String) Precision of uptime percent of this group for the last 90 days..
- `view_mode` (String) Layout of this group on the status page. 'List' renders resources stacked vertically (default). 'Grid' renders resources as a matrix using row and column axes...
- `row_axis_label` (String) Label shown above the row axis when the group is rendered as a grid (e.g. 'Service', 'Tenant'). Free-form so you can use any dimension you like...
- `column_axis_label` (String) Label shown above the column axis when the group is rendered as a grid (e.g. 'Region', 'Environment'). Free-form so you can use any dimension you like...
- `row_axis_values` (String) Comma-separated list of row labels for the grid (e.g. 'Auth, API, Database'). Determines row order in the grid layout...
- `column_axis_values` (String) Comma-separated list of column labels for the grid (e.g. 'US-East, EU-West, Asia'). Determines column order in the grid layout...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_status_page_group.example <id>
```
