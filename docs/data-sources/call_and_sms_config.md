---
page_title: "oneuptime_call_and_sms_config Data Source - oneuptime"
subcategory: "Notification Channels"
description: |-
  A project's own Twilio account. The project's default accounts send the SMS and make the calls to its team members instead of OneUptime's own. A status page can name one for the SMS to its subscribers, and an incoming call policy one for its calls.
---

# oneuptime_call_and_sms_config (Data Source)

A project's own Twilio account. The project's default accounts send the SMS and make the calls to its team members instead of OneUptime's own. A status page can name one for the SMS to its subscribers, and an incoming call policy one for its calls.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one call and sms config may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_call_and_sms_config" "example" {
  name = "Example call and sms config"
}

# Or by id:
data "oneuptime_call_and_sms_config" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_project_default` (Boolean) Whether the SMS and calls sent to project team members (on-call notifications, alerts, verification codes, etc.) go through this Twilio config instead of the global config - the ones it is used for. A project has one default for SMS and one for calls, so a config used for SMS only and one used for calls only can both be the default. Setting this to true makes the config the project default for everything it is used for, taking that from the project's other configs; false makes it the default for nothing. It reads true while the config is the default for SMS, calls or both: isProjectDefaultForSms and isProjectDefaultForCalls say which. A project's first Twilio config becomes the project default when this is not set on create; an explicit false is kept. This does not affect status pages — they use the config explicitly assigned to each status page.
- `is_project_default_for_calls` (Boolean) Whether the calls made to project team members go through this Twilio config: it is the project default for calls. Set by OneUptime from isProjectDefault and what the config is used for; at most one config of a project is the default for calls.
- `is_project_default_for_sms` (Boolean) Whether the SMS sent to project team members go through this Twilio config: it is the project default for SMS. Set by OneUptime from isProjectDefault and what the config is used for; at most one config of a project is the default for SMS.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Permissions - Create: [Project Owner, Project Admin, Create Call and SMS], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Call and SMS], Update: [Project Owner, Project Admin, Edit Call and SMS]
- `slug` (String) Friendly globally unique name for your object.
- `twilio_account_sid` (String) Account SID for your Twilio Account.
- `twilio_secondary_phone_numbers` (String) Secondary Phone Number for your Twilio account.
- `twilio_sms_phone_numbers` (String) Comma-separated numbers to send SMS from, for phone numbers that can only make calls. Once set, every SMS goes out from the one in the recipient's country, or else the first. Leave it empty to send SMS from the same numbers as calls.
- `used_for` (String) What this Twilio config is used for: "Calls and SMS" (the default), "SMS only" or "Calls only". Calls and SMS and Calls only need a Twilio Primary Phone Number; SMS only needs Twilio SMS Phone Numbers. A config used for SMS only makes no calls, and one used for calls only sends no SMS: those go through OneUptime's own Twilio account, as for a project without a config.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `twilio_auth_token` (String, Sensitive) Auth Token for your Twilio Account.
- `twilio_primary_phone_number` (String) Primary Phone Number for your Twilio account.
- `updated_at` (String) Date and Time when the object was updated.
