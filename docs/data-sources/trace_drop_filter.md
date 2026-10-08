---
page_title: "oneuptime_trace_drop_filter Data Source - oneuptime"
subcategory: "Other"
description: |-
  Configure rules to drop or sample spans before storage to reduce volume and cost.
---

# oneuptime_trace_drop_filter (Data Source)

Configure rules to drop or sample spans before storage to reduce volume and cost.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one trace drop filter may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_trace_drop_filter" "example" {
  name = "Example trace drop filter"
}

# Or by id:
data "oneuptime_trace_drop_filter" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `action` (String) What to do with matching spans: 'drop' to discard entirely, 'sample' to keep a percentage.
- `created_by_user_id` (String) ID of the user who created this trace drop filter. The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of what this drop filter does.
- `dropped_count` (Number) Total number of spans this filter has discarded since it was created.
- `filter_query` (String) Filter expression that identifies which spans to drop or sample.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this drop filter is active.
- `name` (String) Friendly name for this drop filter.
- `sample_percentage` (Number) When action is 'sample', the percentage of matching spans to keep (1-99).
- `sort_order` (Number) Where this filter is evaluated among the project's span drop filters, lowest number first. A new filter is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_dropped_at` (String) When this filter most recently discarded a span. Null means it has never matched anything.
- `project_id` (String) ID of the project this trace drop filter belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
