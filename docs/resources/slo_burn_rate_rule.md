---
page_title: "oneuptime_slo_burn_rate_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure multi-window burn rate rules that raise alerts when a Service Level Objective consumes its error budget too quickly
---

# oneuptime_slo_burn_rate_rule (Resource)

Configure multi-window burn rate rules that raise alerts when a Service Level Objective consumes its error budget too quickly

## Example Usage

```terraform
resource "oneuptime_slo_burn_rate_rule" "example" {
  service_level_objective_id = "123e4567-e89b-12d3-a456-426614174000"
  name = "Example short text"
  burn_rate_threshold = 42
  long_window_in_minutes = 42
  short_window_in_minutes = 42
}
```

## Schema

### Required

- `service_level_objective_id` (String) A unique identifier for an object, represented as a UUID..
- `name` (String) Name of this burn rate rule..
- `burn_rate_threshold` (Number) Alert when the burn rate in both the long and short windows is at or above this threshold (e.g. 14.4)...
- `long_window_in_minutes` (Number) Length of the long lookback window in minutes (e.g. 60). The alert fires when both windows exceed the threshold and resolves when the long window drops below it...
- `short_window_in_minutes` (Number) Length of the short lookback window in minutes (e.g. 5). Guards against alerting on burn that has already stopped...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `is_enabled` (Bool) Whether this burn rate rule is enabled..
- `minimum_sample_count` (Number) For event-based SLIs only: skip this rule when the long window has fewer than this many total events. Prevents noisy alerts on low traffic...
- `refire_suppression_minutes` (Number) Minimum number of minutes after an alert resolves before this rule can fire again. Defaults to the long window length when not set...
- `alert_severity_id` (String) A unique identifier for an object, represented as a UUID..
- `on_call_duty_policies` (Set) On-call duty policies attached to alerts created by this burn rate rule...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `last_alert_created_at` (String) A date time object..
- `last_alert_resolved_at` (String) A date time object..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_slo_burn_rate_rule.example <id>
```
