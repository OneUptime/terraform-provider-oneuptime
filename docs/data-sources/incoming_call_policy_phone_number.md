---
page_title: "oneuptime_incoming_call_policy_phone_number Data Source - oneuptime"
subcategory: "Other"
description: |-
  Phone numbers that route incoming calls to an incoming call policy.
---

# oneuptime_incoming_call_policy_phone_number (Data Source)

Phone numbers that route incoming calls to an incoming call policy.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incoming call policy phone number may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incoming_call_policy_phone_number" "example" {
  incoming_call_policy_id = oneuptime_incoming_call_policy.example.id
}

# Or by id:
data "oneuptime_incoming_call_policy_phone_number" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `area_code` (String) Area code associated with this phone number.
- `call_provider_phone_number_id` (String) The call provider identifier for this phone number.
- `country_code` (String) Country code associated with this phone number.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incoming_call_policy_id` (String) ID of the policy that receives calls to this number. The ID of a `oneuptime_incoming_call_policy`.
- `project_call_sms_config_id` (String) ID of the call provider configuration that owns this number.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `phone_number` (String) Phone number that routes calls to the policy.
- `phone_number_purchased_at` (String) When this phone number was attached to the policy.
- `project_id` (String) ID of the project that owns this phone number. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
