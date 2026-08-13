---
page_title: "oneuptime_log_drop_filter Data Source - oneuptime"
subcategory: "Other"
description: |-
  Configure rules to drop or sample logs before storage to reduce volume and cost.
---

# oneuptime_log_drop_filter (Data Source)

Configure rules to drop or sample logs before storage to reduce volume and cost. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_log_drop_filter" "by_name" {
  name = "example-log_drop_filter"
}

data "oneuptime_log_drop_filter" "by_id" {
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
- `description` (String) Description of what this drop filter does... Computed.
- `filter_query` (String) Filter expression that identifies which logs to drop or sample... Computed.
- `action` (String) What to do with matching logs: 'drop' to discard entirely, 'sample' to keep a percentage... Computed.
- `sample_percentage` (Number) When action is 'sample', the percentage of matching logs to keep (1-99)... Computed.
- `is_enabled` (Bool) Whether this drop filter is active... Computed.
- `sort_order` (Number) Determines the evaluation order of this filter relative to others... Computed.
- `dropped_count` (Number) Total number of logs this filter has discarded since it was created... Computed.
- `last_dropped_at` (String) A date time object.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
