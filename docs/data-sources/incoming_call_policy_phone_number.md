---
page_title: "oneuptime_incoming_call_policy_phone_number Data Source - oneuptime"
subcategory: "Other"
description: |-
  Phone numbers that route incoming calls to an incoming call policy.
---

# oneuptime_incoming_call_policy_phone_number (Data Source)

Phone numbers that route incoming calls to an incoming call policy. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_incoming_call_policy_phone_number" "by_name" {
  name = "example-incoming_call_policy_phone_number"
}

data "oneuptime_incoming_call_policy_phone_number" "by_id" {
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
- `incoming_call_policy_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `project_call_sms_config_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `phone_number` (String) Phone object. Computed.
- `call_provider_phone_number_id` (String) The call provider identifier for this phone number... Computed.
- `country_code` (String) Country code associated with this phone number... Computed.
- `area_code` (String) Area code associated with this phone number... Computed.
- `phone_number_purchased_at` (String) A date time object.. Computed.
