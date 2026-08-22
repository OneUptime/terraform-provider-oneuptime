---
page_title: "oneuptime_llm_cost_budget Data Source - oneuptime"
subcategory: "Other"
description: |-
  Daily USD spend budgets for LLM / GenAI telemetry. A worker sums the day's LLM span cost and raises alerts when spend crosses the warning and breach thresholds.
---

# oneuptime_llm_cost_budget (Data Source)

Daily USD spend budgets for LLM / GenAI telemetry. A worker sums the day's LLM span cost and raises alerts when spend crosses the warning and breach thresholds. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_llm_cost_budget" "by_name" {
  name = "example-llm_cost_budget"
}

data "oneuptime_llm_cost_budget" "by_id" {
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
- `description` (String) Description of what this budget covers... Computed.
- `is_enabled` (Bool) Whether this budget is evaluated... Computed.
- `daily_budget_in_usd` (Number) Daily LLM spend budget in USD, evaluated over the UTC day. Alerts fire at the warning threshold and at 100%... Computed.
- `warning_threshold_percent` (Number) Percentage of the daily budget at which a warning alert is raised (1-99). A breach alert always fires at 100%... Computed.
- `service_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `llm_system` (String) Optional provider filter, e.g. openai, anthropic, aws.bedrock (matches the span's gen_ai provider). Leave empty to count every provider... Computed.
- `llm_model` (String) Optional model filter, e.g. gpt-4o (matches the span's requested model exactly). Leave empty to count every model... Computed.
- `alert_severity_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_duty_policies` (Set) On-call duty policies to execute when this budget raises an alert... Computed.
- `current_day_spend_in_usd` (Number) LLM spend accrued so far in the current UTC day, in USD. Computed by the worker... Computed.
- `spend_last_evaluated_at` (String) A date time object.. Computed.
- `last_warning_alert_created_at` (String) A date time object.. Computed.
- `last_breach_alert_created_at` (String) A date time object.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
