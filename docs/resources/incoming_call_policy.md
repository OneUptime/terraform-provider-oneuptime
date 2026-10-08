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
  name        = "Example incoming call policy"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Any friendly name of this policy.

### Optional

- `description` (String) Friendly description that will help you remember.
- `greeting_message` (String) Custom TTS greeting message for incoming calls.
- `is_enabled` (Boolean) Enable or disable this incoming call policy. Defaults to `true`.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `no_answer_message` (String) Message when escalation is exhausted and no one answers.
- `no_one_available_message` (String) Message when no one is on-call or reachable.
- `project_call_sms_config_id` (String) ID of the project-level Twilio configuration. If set, uses this config instead of global config and billing does not apply.
- `repeat_policy_if_no_one_answers` (Boolean) Restart from first rule if all fail. Defaults to `false`.
- `repeat_policy_if_no_one_answers_times` (Number) Maximum repeat attempts if no one answers. Defaults to `1`.

### Read-Only

- `call_provider_phone_number_id` (String) The call provider's ID for the phone number (e.g., Twilio SID).
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `phone_number_area_code` (String) Area code of the phone number.
- `phone_number_country_code` (String) Country code of the phone number (US, GB, etc.).
- `phone_number_purchased_at` (String) When the phone number was purchased.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `routing_phone_number` (String) The phone number for incoming calls to this policy.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incoming call policy by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incoming_call_policy.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incoming_call_policy.example <id>
```
