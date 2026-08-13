---
page_title: "oneuptime_incoming_call_policy Data Source - oneuptime"
subcategory: "Other"
description: |-
  Manage incoming call routing policies with escalation rules for on-call teams
---

# oneuptime_incoming_call_policy (Data Source)

Manage incoming call routing policies with escalation rules for on-call teams Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_incoming_call_policy" "by_name" {
  name = "example-incoming_call_policy"
}

data "oneuptime_incoming_call_policy" "by_id" {
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
- `description` (String) Friendly description that will help you remember.. Computed.
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `routing_phone_number` (String) Phone object. Computed.
- `call_provider_phone_number_id` (String) The call provider's ID for the phone number (e.g., Twilio SID).. Computed.
- `phone_number_country_code` (String) Country code of the phone number (US, GB, etc.).. Computed.
- `phone_number_area_code` (String) Area code of the phone number.. Computed.
- `phone_number_purchased_at` (String) A date time object.. Computed.
- `greeting_message` (String) Custom TTS greeting message for incoming calls.. Computed.
- `no_answer_message` (String) Message when escalation is exhausted and no one answers.. Computed.
- `no_one_available_message` (String) Message when no one is on-call or reachable.. Computed.
- `is_enabled` (Bool) Enable or disable this incoming call policy.. Computed.
- `repeat_policy_if_no_one_answers` (Bool) Restart from first rule if all fail.. Computed.
- `repeat_policy_if_no_one_answers_times` (Number) Maximum repeat attempts if no one answers.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
- `project_call_sms_config_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
