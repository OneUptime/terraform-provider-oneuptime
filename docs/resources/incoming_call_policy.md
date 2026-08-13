---
page_title: "oneuptime_incoming_call_policy Resource - oneuptime"
subcategory: "Other"
description: |-
  Manage incoming call routing policies with escalation rules for on-call teams
---

# oneuptime_incoming_call_policy (Resource)

Manage incoming call routing policies with escalation rules for on-call teams

## Example Usage

```terraform
resource "oneuptime_incoming_call_policy" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Any friendly name of this policy..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `greeting_message` (String) Custom TTS greeting message for incoming calls..
- `no_answer_message` (String) Message when escalation is exhausted and no one answers..
- `no_one_available_message` (String) Message when no one is on-call or reachable..
- `is_enabled` (Bool) Enable or disable this incoming call policy..
- `repeat_policy_if_no_one_answers` (Bool) Restart from first rule if all fail..
- `repeat_policy_if_no_one_answers_times` (Number) Maximum repeat attempts if no one answers..
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `project_call_sms_config_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `routing_phone_number` (String) Phone object.
- `call_provider_phone_number_id` (String) The call provider's ID for the phone number (e.g., Twilio SID)..
- `phone_number_country_code` (String) Country code of the phone number (US, GB, etc.)..
- `phone_number_area_code` (String) Area code of the phone number..
- `phone_number_purchased_at` (String) A date time object..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_incoming_call_policy.example <id>
```
