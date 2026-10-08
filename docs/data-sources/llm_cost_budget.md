---
page_title: "oneuptime_llm_cost_budget Data Source - oneuptime"
subcategory: "Other"
description: |-
  Daily USD spend budgets for LLM / GenAI telemetry. A worker sums the day's LLM span cost and publishes it as the oneuptime.llm.budget.* metrics, so Metrics monitors, dashboards and anomaly detection can act on spend.
---

# oneuptime_llm_cost_budget (Data Source)

Daily USD spend budgets for LLM / GenAI telemetry. A worker sums the day's LLM span cost and publishes it as the oneuptime.llm.budget.* metrics, so Metrics monitors, dashboards and anomaly detection can act on spend.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one llm cost budget may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_llm_cost_budget" "example" {
  name = "Example llm cost budget"
}

# Or by id:
data "oneuptime_llm_cost_budget" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the user who created this budget. The ID of a `oneuptime_user` (see the data source).
- `current_day_spend_in_usd` (Number) LLM spend accrued so far in the current UTC day, in USD. Computed by the worker.
- `daily_budget_in_usd` (Number) Daily LLM spend budget in USD, evaluated over the UTC day. Spend and percent-used are published as metrics for monitors to alert on.
- `description` (String) Description of what this budget covers.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this budget is evaluated.
- `llm_model` (String) Optional model filter, e.g. gpt-4o (matches the span's requested model exactly). Leave empty to count every model.
- `llm_system` (String) Optional provider filter, e.g. openai, anthropic, aws.bedrock (matches the span's gen_ai provider). Leave empty to count every provider.
- `name` (String) Friendly name for this budget.
- `service_id` (String) Optional telemetry service ID scoping this budget. Null means the budget is project-wide. The ID of a `oneuptime_service`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this LLM cost budget belongs to. The ID of a `oneuptime_project`.
- `spend_last_evaluated_at` (String) The last time the worker evaluated this budget. Computed by the worker.
- `updated_at` (String) Date and Time when the object was updated.
