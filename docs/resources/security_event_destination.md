---
page_title: "oneuptime_security_event_destination Resource - oneuptime"
subcategory: "Other"
description: |-
  SIEMs such as Google SecOps that receive the project's detections, or all its security events, from OneUptime.
---

# oneuptime_security_event_destination (Resource)

SIEMs such as Google SecOps that receive the project's detections, or all its security events, from OneUptime.

## Example Usage

```terraform
resource "oneuptime_security_event_destination" "example" {
  name           = "Example security event destination"
  provider_value = "Example short text"
  secrets        = "This is an example of very long text content that might be stored in this field. It can contain a lot of information, such as detailed descriptions, comments, or any other lengthy text data that needs to be stored in the database."
  description    = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Friendly name for this destination.
- `provider_value` (String) Which SIEM this destination sends to, e.g. 'google-secops'. Fixed once created.
- `secrets` (String, Sensitive) Provider-specific secrets, such as a service account key, as a JSON object. Encrypted at rest and never returned by the API.

### Optional

- `config` (String) Provider-specific settings that are not secret, such as the region and the instance. Keys are defined by the provider catalog. A JSON value: write it with `jsonencode()`.
- `description` (String) What this destination is for.
- `is_enabled` (Boolean) Whether events are sent to this destination. Defaults to `true`.
- `send_scope` (String) 'detections' sends what OneUptime found (its detection rules' and threat intel's matches and imported alerts); 'all-events' sends every security event. Defaults to `detections`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this destination. The ID of a `oneuptime_user` (see the data source).
- `cursor` (String) The time up to which events have been sent, as an ISO string: the end of the last window delivered.
- `id` (String) Unique identifier for the resource.
- `last_attempt_at` (String) When OneUptime last tried to send to this destination.
- `last_error` (String) Why the last attempt failed, with credentials redacted, if it did. Cleared by the next attempt that succeeds.
- `last_send_result` (String) What the last attempt sent: events delivered, events the SIEM refused, and the window they came from. A JSON value: write it with `jsonencode()`.
- `last_sent_at` (String) When events were last delivered to this destination.
- `project_id` (String) ID of the project whose events are sent. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing security event destination by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_security_event_destination.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_security_event_destination.example <id>
```
