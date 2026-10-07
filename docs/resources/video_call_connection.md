---
page_title: "oneuptime_video_call_connection Resource - oneuptime"
subcategory: "Other"
description: |-
  Zoom, Google Meet, Microsoft Teams or a standing meeting link, used to start a dedicated video call for incidents and alerts.
---

# oneuptime_video_call_connection (Resource)

Zoom, Google Meet, Microsoft Teams or a standing meeting link, used to start a dedicated video call for incidents and alerts.

## Example Usage

```terraform
resource "oneuptime_video_call_connection" "example" {
  name = jsonencode({
    "_type": "Name",
    "value": "John Doe"
  })
  provider = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name object.
- `provider` (String) Which provider this connection starts calls with: Zoom, GoogleMeet, MicrosoftTeams or CustomLink. Fixed once created...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) What this connection is for...
- `config` (String) Provider-specific, non-secret settings such as the Zoom account and meeting host, the Google Workspace user or the Microsoft Entra tenant and organizer. Keys are defined by the provider catalog...
- `secrets` (String) Provider-specific secrets (a client secret or a service account key) as a JSON object. Encrypted at rest and never returned by the API...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `last_call_started_at` (String) A date time object..
- `last_error` (String) Why the most recent call could not be started, with credentials redacted. Cleared when a call starts...
- `last_error_at` (String) A date time object..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_video_call_connection.example <id>
```
