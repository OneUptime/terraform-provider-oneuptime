---
page_title: "oneuptime_slo_monitor_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Configure rules that automatically attach matching monitors to a Service Level Objective, instead of picking every monitor by hand
---

# oneuptime_slo_monitor_rule (Data Source)

Configure rules that automatically attach matching monitors to a Service Level Objective, instead of picking every monitor by hand

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one slo monitor rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_slo_monitor_rule" "example" {
  name = "Example slo monitor rule"
}

# Or by id:
data "oneuptime_slo_monitor_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this SLO monitor rule.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this rule is enabled. A disabled rule matches nothing, so monitors that only this rule attached are detached from the SLO.
- `monitor_description_pattern` (String) Regex (case-insensitive) matched against the monitor description. Leave empty to skip the description filter.
- `monitor_name_pattern` (String) Regex (case-insensitive) matched against the monitor name. Leave empty to skip the name filter. Use .* to match every monitor.
- `monitor_type` (String) Only match monitors of this type. Leave empty to skip the type filter.
- `name` (String) Name of this SLO monitor rule.
- `service_level_objective_id` (String) ID of the Service Level Objective this monitor rule belongs to. The ID of a `oneuptime_service_level_objective`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `monitor_labels` (Set of String) Only match monitors that carry at least one of these labels. Leave empty to skip the label filter. IDs of `oneuptime_label` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
