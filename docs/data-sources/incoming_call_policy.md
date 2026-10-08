---
page_title: "oneuptime_incoming_call_policy Data Source - oneuptime"
subcategory: "Other"
description: |-
  Manage incoming call routing policies with escalation rules for on-call teams
---

# oneuptime_incoming_call_policy (Data Source)

Manage incoming call routing policies with escalation rules for on-call teams

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incoming call policy may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incoming_call_policy" "example" {
  name = "Example incoming call policy"
}

# Or by id:
data "oneuptime_incoming_call_policy" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `call_provider_phone_number_id` (String) The call provider's ID for the phone number (e.g., Twilio SID).
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `greeting_message` (String) Custom TTS greeting message for incoming calls.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Enable or disable this incoming call policy.
- `name` (String) Any friendly name of this policy.
- `no_answer_message` (String) Message when escalation is exhausted and no one answers.
- `no_one_available_message` (String) Message when no one is on-call or reachable.
- `phone_number_area_code` (String) Area code of the phone number.
- `phone_number_country_code` (String) Country code of the phone number (US, GB, etc.).
- `project_call_sms_config_id` (String) ID of the project-level Twilio configuration. If set, uses this config instead of global config and billing does not apply.
- `repeat_policy_if_no_one_answers` (Boolean) Restart from first rule if all fail.
- `repeat_policy_if_no_one_answers_times` (Number) Maximum repeat attempts if no one answers.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `phone_number_purchased_at` (String) When the phone number was purchased.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `routing_phone_number` (String) The phone number for incoming calls to this policy.
- `updated_at` (String) Date and Time when the object was updated.
