---
page_title: "oneuptime_status_page_group Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage groups on your status page and categorize resources like monitors into these groups.
---

# oneuptime_status_page_group (Data Source)

Manage groups on your status page and categorize resources like monitors into these groups. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_status_page_group" "by_name" {
  name = "example-status_page_group"
}

data "oneuptime_status_page_group" "by_id" {
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
- `parent_status_page_group_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `description` (String) Description for this group. This is visible on Status Page. This can be in markdown format... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `order` (Number) Order / Priority of this resource.. Computed.
- `is_expanded_by_default` (Bool) Is this group expanded by default on Status Page?.. Computed.
- `show_current_status` (Bool) Show current status like offline, operational or degraded... Computed.
- `show_uptime_percent` (Bool) Show uptime percent of this group for the last 90 days.. Computed.
- `uptime_percent_precision` (String) Precision of uptime percent of this group for the last 90 days.. Computed.
- `view_mode` (String) Layout of this group on the status page. 'List' renders resources stacked vertically (default). 'Grid' renders resources as a matrix using row and column axes... Computed.
- `row_axis_label` (String) Label shown above the row axis when the group is rendered as a grid (e.g. 'Service', 'Tenant'). Free-form so you can use any dimension you like... Computed.
- `column_axis_label` (String) Label shown above the column axis when the group is rendered as a grid (e.g. 'Region', 'Environment'). Free-form so you can use any dimension you like... Computed.
- `row_axis_values` (String) Comma-separated list of row labels for the grid (e.g. 'Auth, API, Database'). Determines row order in the grid layout... Computed.
- `column_axis_values` (String) Comma-separated list of column labels for the grid (e.g. 'US-East, EU-West, Asia'). Determines column order in the grid layout... Computed.
