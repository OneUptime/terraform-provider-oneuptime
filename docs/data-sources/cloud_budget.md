---
page_title: "oneuptime_cloud_budget Data Source - oneuptime"
subcategory: "Cloud"
description: |-
  Monthly budgets for what your connected cloud accounts cost. OneUptime alerts the first time in a month the spend passes each threshold, and when the month is on course to pass the budget.
---

# oneuptime_cloud_budget (Data Source)

Monthly budgets for what your connected cloud accounts cost. OneUptime alerts the first time in a month the spend passes each threshold, and when the month is on course to pass the budget.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one cloud budget may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_cloud_budget" "example" {
  name = "Example cloud budget"
}

# Or by id:
data "oneuptime_cloud_budget" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alert_on_projection` (Boolean) Also alert when the month is on course to pass the budget, at the pace of the last full week: before the money is spent.
- `alert_severity_id` (String) ID of the severity of the alerts the budget raises. Empty: the project's least severe. The ID of a `oneuptime_alert_severity`.
- `alert_thresholds` (String) The shares of the budget, in percent, an alert is raised at, comma separated: the first time in a month the spend passes each.
- `cloud_account_id` (String) ID of the account the budget is for. Empty: every account of the project. The ID of a `oneuptime_cloud_account`.
- `created_by_user_id` (String) ID of the user who created this budget. The ID of a `oneuptime_user` (see the data source).
- `currency` (String) The currency of the amount, as the clouds bill the accounts: an account billed in another is not counted.
- `description` (String) What the budget is for, in your team's words.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether the budget is watched and raises alerts.
- `monthly_amount` (Number) What the accounts may cost in a calendar month (UTC), in the budget's currency.
- `name` (String) What the budget is called.
- `projected_this_month` (Number) What the accounts are on course to cost this month, at the pace of the last full week, when last checked.
- `spent_this_month` (Number) What the accounts cost this month so far, when last checked.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_checked_at` (String) When the budget was last compared with what the accounts cost.
- `on_call_duty_policies` (Set of String) The on-call policies the budget's alerts page. IDs of `oneuptime_on_call_policy` resources.
- `project_id` (String) ID of the project this budget belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
