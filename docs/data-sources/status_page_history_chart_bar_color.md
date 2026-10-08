---
page_title: "oneuptime_status_page_history_chart_bar_color Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Modify the colors of the history chart bars on Status Page
---

# oneuptime_status_page_history_chart_bar_color (Data Source)

Modify the colors of the history chart bars on Status Page

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one status page history chart bar color may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_status_page_history_chart_bar_color" "example" {
  status_page_id = oneuptime_status_page.example.id
}

# Or by id:
data "oneuptime_status_page_history_chart_bar_color" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `order` (Number) Where this rule sits in the list. Rules are evaluated from the top of the list down, lowest number first. A new rule is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.
- `status_page_id` (String) ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.
- `uptime_percent_greater_than_or_equal_to` (Number) Uptime percent greater than or equal to this value.

### Read-Only

- `bar_color` (String) Color of the bar chart when this rule matches (#32a852 for example).
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
