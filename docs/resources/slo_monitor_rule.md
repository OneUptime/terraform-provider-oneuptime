---
page_title: "oneuptime_slo_monitor_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure rules that automatically attach matching monitors to a Service Level Objective, instead of picking every monitor by hand
---

# oneuptime_slo_monitor_rule (Resource)

Configure rules that automatically attach matching monitors to a Service Level Objective, instead of picking every monitor by hand

## Example Usage

```terraform
resource "oneuptime_slo_monitor_rule" "example" {
  service_level_objective_id = "123e4567-e89b-12d3-a456-426614174000"
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `service_level_objective_id` (String) A unique identifier for an object, represented as a UUID..
- `name` (String) Name of this SLO monitor rule..

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource...
- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this SLO monitor rule..
- `is_enabled` (Bool) Whether this rule is enabled. A disabled rule matches nothing, so monitors that only this rule attached are detached from the SLO...
- `monitor_labels` (Set) Only match monitors that carry at least one of these labels. Leave empty to skip the label filter...
- `monitor_name_pattern` (String) Regex (case-insensitive) matched against the monitor name. Leave empty to skip the name filter. Use .* to match every monitor...
- `monitor_description_pattern` (String) Regex (case-insensitive) matched against the monitor description. Leave empty to skip the description filter...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_slo_monitor_rule.example <id>
```
