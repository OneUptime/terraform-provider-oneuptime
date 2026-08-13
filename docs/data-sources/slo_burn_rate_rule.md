---
page_title: "oneuptime_slo_burn_rate_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Configure multi-window burn rate rules that raise alerts when a Service Level Objective consumes its error budget too quickly
---

# oneuptime_slo_burn_rate_rule (Data Source)

Configure multi-window burn rate rules that raise alerts when a Service Level Objective consumes its error budget too quickly Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_slo_burn_rate_rule" "by_name" {
  name = "example-slo_burn_rate_rule"
}

data "oneuptime_slo_burn_rate_rule" "by_id" {
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
- `service_level_objective_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_enabled` (Bool) Whether this burn rate rule is enabled.. Computed.
- `burn_rate_threshold` (Number) Alert when the burn rate in both the long and short windows is at or above this threshold (e.g. 14.4)... Computed.
- `long_window_in_minutes` (Number) Length of the long lookback window in minutes (e.g. 60). The alert fires when both windows exceed the threshold and resolves when the long window drops below it... Computed.
- `short_window_in_minutes` (Number) Length of the short lookback window in minutes (e.g. 5). Guards against alerting on burn that has already stopped... Computed.
- `minimum_sample_count` (Number) For event-based SLIs only: skip this rule when the long window has fewer than this many total events. Prevents noisy alerts on low traffic... Computed.
- `refire_suppression_minutes` (Number) Minimum number of minutes after an alert resolves before this rule can fire again. Defaults to the long window length when not set... Computed.
- `alert_severity_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_duty_policies` (Set) On-call duty policies attached to alerts created by this burn rate rule... Computed.
- `last_alert_created_at` (String) A date time object.. Computed.
- `last_alert_resolved_at` (String) A date time object.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
