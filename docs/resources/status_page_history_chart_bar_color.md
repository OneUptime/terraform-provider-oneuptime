---
page_title: "oneuptime_status_page_history_chart_bar_color Resource - oneuptime"
subcategory: "Status Pages"
description: |-
  Modify the colors of the history chart bars on Status Page
---

# oneuptime_status_page_history_chart_bar_color (Resource)

Modify the colors of the history chart bars on Status Page

## Example Usage

```terraform
resource "oneuptime_status_page_history_chart_bar_color" "example" {
  status_page_id                          = oneuptime_status_page.example.id
  uptime_percent_greater_than_or_equal_to = 42
  bar_color                               = "#ff0000"
}
```

## Schema

### Required

- `bar_color` (String) Color of the bar chart when this rule matches (#32a852 for example).
- `status_page_id` (String) ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.
- `uptime_percent_greater_than_or_equal_to` (Number) Uptime percent greater than or equal to this value.

### Optional

- `order` (Number) Where this rule sits in the list. Rules are evaluated from the top of the list down, lowest number first. A new rule is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing status page history chart bar color by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_status_page_history_chart_bar_color.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_status_page_history_chart_bar_color.example <id>
```
