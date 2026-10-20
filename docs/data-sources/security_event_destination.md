---
page_title: "oneuptime_security_event_destination Data Source - oneuptime"
subcategory: "Other"
description: |-
  SIEMs such as Google SecOps that receive the project's detections, or all its security events, from OneUptime.
---

# oneuptime_security_event_destination (Data Source)

SIEMs such as Google SecOps that receive the project's detections, or all its security events, from OneUptime.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one security event destination may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_security_event_destination" "example" {
  name = "Example security event destination"
}

# Or by id:
data "oneuptime_security_event_destination" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the user who created this destination. The ID of a `oneuptime_user` (see the data source).
- `cursor` (String) The time up to which events have been sent, as an ISO string: the end of the last window delivered.
- `description` (String) What this destination is for.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether events are sent to this destination.
- `last_error` (String) Why the last attempt failed, with credentials redacted, if it did. Cleared by the next attempt that succeeds.
- `name` (String) Friendly name for this destination.
- `provider_value` (String) Which SIEM this destination sends to, e.g. 'google-secops'. Fixed once created.
- `send_scope` (String) 'detections' sends what OneUptime found (its detection rules' and threat intel's matches and imported alerts); 'all-events' sends every security event.

### Read-Only

- `config` (String) Provider-specific settings that are not secret, such as the region and the instance. Keys are defined by the provider catalog. A JSON value: write it with `jsonencode()`.
- `created_at` (String) Date and Time when the object was created.
- `last_attempt_at` (String) When OneUptime last tried to send to this destination.
- `last_send_result` (String) What the last attempt sent: events delivered, events the SIEM refused, and the window they came from. A JSON value: write it with `jsonencode()`.
- `last_sent_at` (String) When events were last delivered to this destination.
- `project_id` (String) ID of the project whose events are sent. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
