---
page_title: "oneuptime_user_override Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  While someone is away, a user override sends the alerts that would page them to the person who covers, for a set time. An override on an on-call policy applies to that policy only; one without a policy applies to every on-call policy.
---

# oneuptime_user_override (Data Source)

While someone is away, a user override sends the alerts that would page them to the person who covers, for a set time. An override on an on-call policy applies to that policy only; one without a policy applies to every on-call policy.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one user override may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_user_override" "example" {
  on_call_duty_policy_id = oneuptime_on_call_policy.example.id
}

# Or by id:
data "oneuptime_user_override" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `on_call_duty_policy_id` (String) ID of the on-call policy this override applies to. Leave it empty for a global override, which applies to every on-call policy. The ID of a `oneuptime_on_call_policy`.
- `override_user_id` (String) ID of the user who is away. While the override is in force, alerts that would page this user go to the user in Route Alerts To User ID instead. The ID of a `oneuptime_user` (see the data source).
- `route_alerts_to_user_id` (String) ID of the user who covers. While the override is in force, this user gets the alerts that would page the user in Override User ID. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `ends_at` (String) When the override ends, which has to be after it starts. From then on, alerts page the Override User again.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `starts_at` (String) When the override starts sending the Override User's alerts to the Route Alerts To User.
- `updated_at` (String) Date and Time when the object was updated.
