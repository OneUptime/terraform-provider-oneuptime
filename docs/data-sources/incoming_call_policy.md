---
page_title: "oneuptime_incoming_call_policy Data Source - oneuptime"
subcategory: "On-Call & Escalation"
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
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `greeting_message` (String) Custom TTS greeting message for incoming calls.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Enable or disable this incoming call policy.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Any friendly name of this policy.
- `no_answer_message` (String) Message when escalation is exhausted and no one answers.
- `no_one_available_message` (String) Message when no one is on-call or reachable.
- `phone_number_area_code` (String) Area code of the phone number.
- `phone_number_country_code` (String) Country code of the phone number (US, GB, etc.).
- `project_call_sms_config_id` (String) ID of the project-level Twilio configuration. If set, uses this config instead of global config and billing does not apply. The ID of a `oneuptime_call_and_sms_config`.
- `repeat_policy_if_no_one_answers` (Boolean) Restart from first rule if all fail.
- `repeat_policy_if_no_one_answers_times` (Number) How many more times the escalation rules ring again from the first when nobody answers: a whole number of 1 or more while Repeat Policy If No One Answers is on.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `labels` (Set of String) Labels to group and filter by. IDs of `oneuptime_label` resources.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `phone_number_purchased_at` (String) When the phone number was purchased.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `routing_phone_number` (String) The phone number for incoming calls to this policy.
- `updated_at` (String) Date and Time when the object was updated.
