---
page_title: "oneuptime_incident_internal_note Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Manage internal notes for your incident
---

# oneuptime_incident_internal_note (Resource)

Manage internal notes for your incident

## Example Usage

```terraform
resource "oneuptime_incident_internal_note" "example" {
  incident_id = oneuptime_incident.example.id
}
```

## Schema

### Required

- `incident_id` (String) Relation to Incident ID in which this resource belongs. The ID of a `oneuptime_incident`.

### Optional

- `attachments` (Set of String) Files attached to this note. IDs of `oneuptime_file` resources.
- `note` (String) Notes in markdown.
- `posted_at` (String) Date and time the note belongs to in the timeline: when it was written, or - for a note saved from a Slack or Microsoft Teams message - when that message was posted.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `posted_from_slack_message_id` (String) Unique identifier for the Slack message this note was created from (channel_id:message_ts). Used to prevent duplicate notes when multiple users react to the same message.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `saved_from_chat_message` (String) When this note was saved from a Slack or Microsoft Teams message - pinned, or given a note emoji - the chat it was posted in (workspaceType), who wrote it (authorName, and authorUserId when their chat account is connected to OneUptime) and when it was posted there (postedAt). Written by OneUptime; the note's creator is whoever saved it. A JSON value: read it with `jsondecode()`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incident internal note by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incident_internal_note.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incident_internal_note.example <id>
```
