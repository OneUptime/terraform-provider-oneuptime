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
  status_page_id = "123e4567-e89b-12d3-a456-426614174000"
  uptime_percent_greater_than_or_equal_to = 42
  bar_color = jsonencode({
    "_type": "Color",
    "value": "#ff0000"
  })
}
```

## Schema

### Required

- `status_page_id` (String) A unique identifier for an object, represented as a UUID..
- `uptime_percent_greater_than_or_equal_to` (Number) Uptime percent greater than or equal to this value..
- `bar_color` (String) Color object.

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `order` (Number) Order / Priority of this resource..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_status_page_history_chart_bar_color.example <id>
```
