---
page_title: "oneuptime_incoming_call_policy_phone_number Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Phone numbers that route incoming calls to an incoming call policy.
---

# oneuptime_incoming_call_policy_phone_number (Data Source)

Phone numbers that route incoming calls to an incoming call policy.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incoming call policy phone number may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incoming_call_policy_phone_number" "example" {
  created_by_actor_type = "User"
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
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incoming_call_policy_id` (String) ID of the policy that receives calls to this number. The ID of a `oneuptime_incoming_call_policy`.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_call_sms_config_id` (String) ID of the call provider configuration that owns this number. The ID of a `oneuptime_call_and_sms_config`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `phone_number` (String) Phone number that routes calls to the policy.
- `phone_number_purchased_at` (String) When this phone number was attached to the policy.
- `project_id` (String) ID of the project that owns this phone number. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
