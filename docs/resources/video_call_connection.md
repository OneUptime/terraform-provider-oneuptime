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
  name           = "Example video call connection"
  provider_value = "Example short text"
  description    = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Friendly name for this connection, shown wherever a call is started, e.g. 'Incident Zoom'.
- `provider_value` (String) Which provider this connection starts calls with: Zoom, GoogleMeet, MicrosoftTeams or CustomLink. Fixed once created.

### Optional

- `auth_method` (String) How this connection signs in to its provider: OAuth when someone connected it by signing in to Zoom, Google or Microsoft (Connect in Project Settings > Video Calls), AppCredentials when it uses the project's own app - a Zoom Server-to-Server OAuth app, a Google service account or a Microsoft Entra app registration. Empty for a meeting link. Fixed once created. A connection made by signing in is created by signing in, never through the API.
- `config` (String) Provider-specific, non-secret settings such as the Zoom account and meeting host, the Google Workspace user or the Microsoft Entra tenant and organizer. Keys are defined by the provider catalog. A JSON value: write it with `jsonencode()`.
- `description` (String) What this connection is for.
- `secrets` (String) Provider-specific secrets (a client secret or a service account key) as a JSON object. Encrypted at rest and never returned by the API. A connection made by signing in keeps its sign-in's tokens here, which only OneUptime writes.

### Read-Only

- `connected_account` (String) For a connection made by signing in: the Zoom, Google or Microsoft account that signed in, which every meeting is created as. Set by OneUptime when someone connects or reconnects, and cleared when the account removes OneUptime.
- `connected_account_id` (String) For a connection made by signing in: the provider's id of the account that signed in (a Zoom user ID, a Google account ID, a Microsoft Entra object ID). Connections signed in as the same account share one sign-in, because Zoom keeps only one per account.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this connection. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_call_started_at` (String) When this connection last started a call.
- `last_error` (String) Why the most recent call could not be started, with credentials redacted. Cleared when a call starts.
- `last_error_at` (String) When the most recent call could not be started.
- `project_id` (String) ID of the project this connection belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing video call connection by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_video_call_connection.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_video_call_connection.example <id>
```
