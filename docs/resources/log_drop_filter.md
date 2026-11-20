---
page_title: "oneuptime_log_drop_filter Resource - oneuptime"
subcategory: "Logs, Traces & Metrics"
description: |-
  Configure rules to drop or sample logs before storage to reduce volume and cost.
---

# oneuptime_log_drop_filter (Resource)

Configure rules to drop or sample logs before storage to reduce volume and cost.

## Example Usage

```terraform
resource "oneuptime_log_drop_filter" "example" {
  name         = "Example log drop filter"
  filter_query = "This is an example of longer text content that might be stored in this field."
  action       = "drop"
  description  = "Managed by Terraform"
}
```

<a id="order"></a>
## Order

`sort_order` is this log drop filter's place among the log drop filters of its project, lowest first. Setting it places this log drop filter and moves no other, so a configuration that numbers every log drop filter of a project is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.log_drop_filters`:

```terraform
sort_order = index(local.log_drop_filters, each.key) + 1
```

- Give each log drop filter a number of its own: two log drop filters at one number are shown the older first.
- A log drop filter created without `sort_order` goes after the last one of its project.
- Deleting a log drop filter leaves a gap in the numbers, which changes nothing: log drop filters go by their numbers, not by counting.
- Moving a log drop filter in the OneUptime Dashboard numbers the project's log drop filters 1..n again; the next plan shows that as changes to `sort_order`.


## Schema

### Required

- `action` (String) What to do with matching logs: 'drop' to discard entirely, 'sample' to keep a percentage. Allowed values: `drop`, `sample`.
- `filter_query` (String) Filter expression that identifies which logs to drop or sample.
- `name` (String) Friendly name for this drop filter.

### Optional

- `description` (String) Description of what this drop filter does.
- `is_enabled` (Boolean) Whether this drop filter is active. Defaults to `true`.
- `sample_percentage` (Number) When action is 'sample', the percentage of matching logs to keep (1-99).
- `sort_order` (Number) Where this filter is evaluated among the project's log drop filters, lowest number first. A new filter is added to the end of the list. Setting a number places it there and moves no other one: two at one number are read the older first. In the dashboard, drag the rows to reorder them. See [Order](#order).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this log drop filter. The ID of a `oneuptime_user` (see the data source).
- `dropped_count` (Number) Total number of logs this filter has discarded since it was created.
- `id` (String) Unique identifier for the resource.
- `last_dropped_at` (String) When this filter most recently discarded a log. Null means it has never matched anything.
- `project_id` (String) ID of the project this log drop filter belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing log drop filter by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_log_drop_filter.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_log_drop_filter.example <id>
```
