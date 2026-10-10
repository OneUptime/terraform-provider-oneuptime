---
page_title: "oneuptime_video_call_connection Data Source - oneuptime"
subcategory: "Other"
description: |-
  Zoom, Google Meet, Microsoft Teams or a standing meeting link, used to start a dedicated video call for incidents and alerts.
---

# oneuptime_video_call_connection (Data Source)

Zoom, Google Meet, Microsoft Teams or a standing meeting link, used to start a dedicated video call for incidents and alerts.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one video call connection may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_video_call_connection" "example" {
  name = "Example video call connection"
}

# Or by id:
data "oneuptime_video_call_connection" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `auth_method` (String) How this connection signs in to its provider: OAuth when someone connected it by signing in to Zoom, Google or Microsoft (Connect in Project Settings > Video Calls), AppCredentials when it uses the project's own app - a Zoom Server-to-Server OAuth app, a Google service account or a Microsoft Entra app registration. Empty for a meeting link. Fixed once created. A connection made by signing in is created by signing in, never through the API.
- `connected_account` (String) For a connection made by signing in: the Zoom, Google or Microsoft account that signed in, which every meeting is created as. Set by OneUptime when someone connects or reconnects, and cleared when the account removes OneUptime.
- `connected_account_id` (String) For a connection made by signing in: the provider's id of the account that signed in (a Zoom user ID, a Google account ID, a Microsoft Entra object ID). Connections signed in as the same account share one sign-in, because Zoom keeps only one per account.
- `created_by_user_id` (String) ID of the user who created this connection. The ID of a `oneuptime_user` (see the data source).
- `description` (String) What this connection is for.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_error` (String) Why the most recent call could not be started, with credentials redacted. Cleared when a call starts.
- `name` (String) Friendly name for this connection, shown wherever a call is started, e.g. 'Incident Zoom'.
- `provider_value` (String) Which provider this connection starts calls with: Zoom, GoogleMeet, MicrosoftTeams or CustomLink. Fixed once created.

### Read-Only

- `config` (String) Provider-specific, non-secret settings such as the Zoom account and meeting host, the Google Workspace user or the Microsoft Entra tenant and organizer. Keys are defined by the provider catalog. A JSON value: write it with `jsonencode()`.
- `created_at` (String) Date and Time when the object was created.
- `last_call_started_at` (String) When this connection last started a call.
- `last_error_at` (String) When the most recent call could not be started.
- `project_id` (String) ID of the project this connection belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
