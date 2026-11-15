---
page_title: "oneuptime_call_and_sms_config Resource - oneuptime"
subcategory: "Notification Channels"
description: |-
  A project's own Twilio account. The project's default accounts send the SMS and make the calls to its team members instead of OneUptime's own. A status page can name one for the SMS to its subscribers, and an incoming call policy one for its calls.
---

# oneuptime_call_and_sms_config (Resource)

A project's own Twilio account. The project's default accounts send the SMS and make the calls to its team members instead of OneUptime's own. A status page can name one for the SMS to its subscribers, and an incoming call policy one for its calls.

## Example Usage

```terraform
resource "oneuptime_call_and_sms_config" "example" {
  name        = "Example call and sms config"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String)

### Optional

- `description` (String) Friendly description that will help you remember.
- `is_project_default` (Boolean) Whether the SMS and calls sent to project team members (on-call notifications, alerts, verification codes, etc.) go through this Twilio config instead of the global config - the ones it is used for. A project has one default for SMS and one for calls, so a config used for SMS only and one used for calls only can both be the default. Setting this to true makes the config the project default for everything it is used for, taking that from the project's other configs; false makes it the default for nothing. It reads true while the config is the default for SMS, calls or both: isProjectDefaultForSms and isProjectDefaultForCalls say which. A project's first Twilio config becomes the project default when this is not set on create; an explicit false is kept. This does not affect status pages — they use the config explicitly assigned to each status page. Defaults to `false`.
- `twilio_account_sid` (String) Account SID for your Twilio Account.
- `twilio_auth_token` (String, Sensitive) Auth Token for your Twilio Account.
- `twilio_primary_phone_number` (String) Primary Phone Number for your Twilio account.
- `twilio_secondary_phone_numbers` (String) Secondary Phone Number for your Twilio account.
- `twilio_sms_phone_numbers` (String) Comma-separated numbers to send SMS from, for phone numbers that can only make calls. Once set, every SMS goes out from the one in the recipient's country, or else the first. Leave it empty to send SMS from the same numbers as calls.
- `used_for` (String) What this Twilio config is used for: "Calls and SMS" (the default), "SMS only" or "Calls only". Calls and SMS and Calls only need a Twilio Primary Phone Number; SMS only needs Twilio SMS Phone Numbers. A config used for SMS only makes no calls, and one used for calls only sends no SMS: those go through OneUptime's own Twilio account, as for a project without a config. Allowed values: `Calls and SMS`, `SMS only`, `Calls only`. Defaults to `Calls and SMS`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_project_default_for_calls` (Boolean) Whether the calls made to project team members go through this Twilio config: it is the project default for calls. Set by OneUptime from isProjectDefault and what the config is used for; at most one config of a project is the default for calls.
- `is_project_default_for_sms` (Boolean) Whether the SMS sent to project team members go through this Twilio config: it is the project default for SMS. Set by OneUptime from isProjectDefault and what the config is used for; at most one config of a project is the default for SMS.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing call and sms config by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_call_and_sms_config.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_call_and_sms_config.example <id>
```
