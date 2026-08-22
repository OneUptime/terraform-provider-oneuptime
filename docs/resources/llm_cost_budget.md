---
page_title: "oneuptime_llm_cost_budget Resource - oneuptime"
subcategory: "Other"
description: |-
  Daily USD spend budgets for LLM / GenAI telemetry. A worker sums the day's LLM span cost and publishes it as the oneuptime.llm.budget.* metrics, so Metrics monitors, dashboards and anomaly detection can act on spend.
---

# oneuptime_llm_cost_budget (Resource)

Daily USD spend budgets for LLM / GenAI telemetry. A worker sums the day's LLM span cost and publishes it as the oneuptime.llm.budget.* metrics, so Metrics monitors, dashboards and anomaly detection can act on spend.

## Example Usage

```terraform
resource "oneuptime_llm_cost_budget" "example" {
  name = jsonencode({
    "_type": "Name",
    "value": "John Doe"
  })
  daily_budget_in_usd = 42
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name object.
- `daily_budget_in_usd` (Number) Daily LLM spend budget in USD, evaluated over the UTC day. Spend and percent-used are published as metrics for monitors to alert on...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of what this budget covers...
- `is_enabled` (Bool) Whether this budget is evaluated...
- `service_id` (String) A unique identifier for an object, represented as a UUID..
- `llm_system` (String) Optional provider filter, e.g. openai, anthropic, aws.bedrock (matches the span's gen_ai provider). Leave empty to count every provider...
- `llm_model` (String) Optional model filter, e.g. gpt-4o (matches the span's requested model exactly). Leave empty to count every model...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `current_day_spend_in_usd` (Number) LLM spend accrued so far in the current UTC day, in USD. Computed by the worker...
- `spend_last_evaluated_at` (String) A date time object..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_llm_cost_budget.example <id>
```
