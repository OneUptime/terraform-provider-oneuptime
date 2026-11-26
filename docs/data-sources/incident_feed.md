---
page_title: "oneuptime_incident_feed Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Log of the entire incident state change. This is a log of all the incident state changes, public notes, more etc.
---

# oneuptime_incident_feed (Data Source)

Log of the entire incident state change. This is a log of all the incident state changes, public notes, more etc.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident feed may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_feed" "example" {
  last_modified_by_user_id = data.oneuptime_user.example.id
  user_email               = "alice@example.com"
}

# Or by id:
data "oneuptime_incident_feed" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `ai_run_id` (String) AI investigation run that produced this feed item when the event is a root cause.
- `api_key_id` (String) The project API key that made this change, when an API key did rather than a person. Kept when the key is revoked or deleted.
- `api_key_name` (String) The name the API key that made this change had at the time, shown as the item's author.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `feed_info_in_markdown` (String) Log of the entire incident state change in Markdown.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_feed_event_type` (String) Incident Feed Event.
- `incident_id` (String) Relation to Incident ID in which this resource belongs. The ID of a `oneuptime_incident`.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `more_information_in_markdown` (String) More information in Markdown.
- `user_email` (String) Look it up by the email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) User who this feed belongs to (if this feed belongs to a User). The ID of a `oneuptime_user` (see the data source). Or look them up by email address, in `user_email`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `display_color` (String) Display color for the incident log.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `posted_at` (String) Date and time when the feed was posted.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `saved_from_chat_message` (String) For a note saved from a Slack or Microsoft Teams message: the chat it was posted in (workspaceType), who wrote it (authorName, and authorUserId when their chat account is connected to OneUptime) and when it was posted there (postedAt). The feed item's user is whoever saved it. Written by OneUptime. A JSON value: read it with `jsondecode()`.
- `updated_at` (String) Date and Time when the object was updated.
