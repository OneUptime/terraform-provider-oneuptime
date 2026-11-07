---
page_title: "oneuptime_cloud_budget Resource - oneuptime"
subcategory: "Cloud"
description: |-
  Monthly budgets for what your connected cloud accounts cost. OneUptime alerts the first time in a month the spend passes each threshold, and when the month is on course to pass the budget.
---

# oneuptime_cloud_budget (Resource)

Monthly budgets for what your connected cloud accounts cost. OneUptime alerts the first time in a month the spend passes each threshold, and when the month is on course to pass the budget.

## Example Usage

```terraform
resource "oneuptime_cloud_budget" "example" {
  name           = "Example cloud budget"
  monthly_amount = 42
  description    = "Managed by Terraform"
}
```

## Schema

### Required

- `monthly_amount` (Number) What the accounts may cost in a calendar month (UTC), in the budget's currency.
- `name` (String) What the budget is called.

### Optional

- `alert_on_anomaly` (Boolean) Also alert when a day costs far more than usual: a service at twice its usual day, or the whole account at one and a half times. Defaults to `true`.
- `alert_on_projection` (Boolean) Also alert when the month is on course to pass the budget, at the pace of the last full week: before the money is spent. Defaults to `true`.
- `alert_severity_id` (String) ID of the severity of the alerts the budget raises. Empty: the project's least severe. The ID of a `oneuptime_alert_severity`.
- `alert_thresholds` (String) The shares of the budget, in percent, an alert is raised at, comma separated: the first time in a month the spend passes each. Defaults to `80,100`.
- `cloud_account_id` (String) ID of the account the budget is for. Empty: every account of the project. The ID of a `oneuptime_cloud_account`.
- `currency` (String) The currency of the amount, as the clouds bill the accounts: an account billed in another is not counted. Defaults to `USD`.
- `description` (String) What the budget is for, in your team's words.
- `is_enabled` (Boolean) Whether the budget is watched and raises alerts. Defaults to `true`.
- `on_call_duty_policies` (Set of String) The on-call policies the budget's alerts page. IDs of `oneuptime_on_call_policy` resources.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this budget. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_checked_at` (String) When the budget was last compared with what the accounts cost.
- `project_id` (String) ID of the project this budget belongs to. The ID of a `oneuptime_project`.
- `projected_this_month` (Number) What the accounts are on course to cost this month, at the pace of the last full week, when last checked.
- `spent_this_month` (Number) What the accounts cost this month so far, when last checked.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing cloud budget by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_cloud_budget.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_cloud_budget.example <id>
```
