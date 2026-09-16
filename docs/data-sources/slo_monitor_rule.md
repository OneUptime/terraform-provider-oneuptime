---
page_title: "oneuptime_slo_monitor_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Configure rules that automatically attach matching monitors to a Service Level Objective, instead of picking every monitor by hand
---

# oneuptime_slo_monitor_rule (Data Source)

Configure rules that automatically attach matching monitors to a Service Level Objective, instead of picking every monitor by hand Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_slo_monitor_rule" "by_name" {
  name = "example-slo_monitor_rule"
}

data "oneuptime_slo_monitor_rule" "by_id" {
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
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource... Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `service_level_objective_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `description` (String) Description of this SLO monitor rule.. Computed.
- `is_enabled` (Bool) Whether this rule is enabled. A disabled rule matches nothing, so monitors that only this rule attached are detached from the SLO... Computed.
- `monitor_labels` (Set) Only match monitors that carry at least one of these labels. Leave empty to skip the label filter... Computed.
- `monitor_name_pattern` (String) Regex (case-insensitive) matched against the monitor name. Leave empty to skip the name filter. Use .* to match every monitor... Computed.
- `monitor_description_pattern` (String) Regex (case-insensitive) matched against the monitor description. Leave empty to skip the description filter... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
