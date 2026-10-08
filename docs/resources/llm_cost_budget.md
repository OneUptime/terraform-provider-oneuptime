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
  name                = "Example llm cost budget"
  daily_budget_in_usd = 42
  description         = "Managed by Terraform"
}
```

## Schema

### Required

- `daily_budget_in_usd` (Number) Daily LLM spend budget in USD, evaluated over the UTC day. Spend and percent-used are published as metrics for monitors to alert on.
- `name` (String) Friendly name for this budget.

### Optional

- `description` (String) Description of what this budget covers.
- `is_enabled` (Boolean) Whether this budget is evaluated. Defaults to `true`.
- `llm_model` (String) Optional model filter, e.g. gpt-4o (matches the span's requested model exactly). Leave empty to count every model.
- `llm_system` (String) Optional provider filter, e.g. openai, anthropic, aws.bedrock (matches the span's gen_ai provider). Leave empty to count every provider.
- `service_id` (String) Optional telemetry service ID scoping this budget. Null means the budget is project-wide. The ID of a `oneuptime_service`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this budget. The ID of a `oneuptime_user` (see the data source).
- `current_day_spend_in_usd` (Number) LLM spend accrued so far in the current UTC day, in USD. Computed by the worker.
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this LLM cost budget belongs to. The ID of a `oneuptime_project`.
- `spend_last_evaluated_at` (String) The last time the worker evaluated this budget. Computed by the worker.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing llm cost budget by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_llm_cost_budget.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_llm_cost_budget.example <id>
```
