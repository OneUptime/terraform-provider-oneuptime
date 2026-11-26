---
page_title: "oneuptime_alert_episode_feed Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Log of the entire alert episode activity. This is a log of all the episode state changes, alerts added/removed, notes, etc.
---

# oneuptime_alert_episode_feed (Resource)

Log of the entire alert episode activity. This is a log of all the episode state changes, alerts added/removed, notes, etc.

## Example Usage

```terraform
resource "oneuptime_alert_episode_feed" "example" {
  alert_episode_id              = oneuptime_alert_episode.example.id
  feed_info_in_markdown         = "# Heading\n\nThis is **markdown** content"
  alert_episode_feed_event_type = "EpisodeCreated"
  display_color                 = "#ff0000"
}
```

## Schema

### Required

- `alert_episode_feed_event_type` (String) Alert Episode Feed Event Type. Allowed values: `EpisodeCreated`, `EpisodeStateChanged`, `EpisodeUpdated`, `AlertAdded`, `AlertRemoved`, `OwnerUserAdded`, `OwnerTeamAdded`, `OwnerUserRemoved`, `OwnerTeamRemoved`, `OwnerNotificationSent`, `PrivateNote`, `RootCause`, `OnCallPolicy`, `OnCallNotification`, `SeverityChanged`, `LabelRuleExecuted`, `OwnerRuleExecuted`, `PrivacyRuleExecuted`, `OnCallRuleExecuted`, `ChannelMemberAdded`.
- `alert_episode_id` (String) Relation to Alert Episode ID in which this resource belongs. The ID of a `oneuptime_alert_episode`.
- `display_color` (String) Display color for the alert episode log.
- `feed_info_in_markdown` (String) Log of the entire alert episode activity in Markdown.

### Optional

- `more_information_in_markdown` (String) More information in Markdown.
- `posted_at` (String) Date and time when the feed was posted.
- `user_email` (String) The email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) User who this feed belongs to (if this feed belongs to a User). The ID of a `oneuptime_user` (see the data source). Or name them by email address, in `user_email`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `saved_from_chat_message` (String) For a note saved from a Slack or Microsoft Teams message: the chat it was posted in (workspaceType), who wrote it (authorName, and authorUserId when their chat account is connected to OneUptime) and when it was posted there (postedAt). The feed item's user is whoever saved it. Written by OneUptime. A JSON value: read it with `jsondecode()`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing alert episode feed by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_alert_episode_feed.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_alert_episode_feed.example <id>
```
