---
page_title: "oneuptime_workspace_notification_rule Data Source - oneuptime"
subcategory: "Slack, Teams & Video Calls"
description: |-
  Notification Rule for Third Party Workspaces
---

# oneuptime_workspace_notification_rule (Data Source)

Notification Rule for Third Party Workspaces

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one workspace notification rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_workspace_notification_rule" "example" {
  name = "Example workspace notification rule"
}

# Or by id:
data "oneuptime_workspace_notification_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of the Notification Rule.
- `event_type` (String) Event Type for the Workspace like Incident Created, Monitor Status Updated, etc.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `name` (String) Name of the Notification Rule.
- `workspace_type` (String) Type of Workspace - slack, microsoft teams etc.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `notification_rule` (String) The rule, as JSON: the events it is for (filters, filterCondition), where it posts, and for incident, alert, incident episode, alert episode, scheduled maintenance, on-call policy and security case rules, updatesToPost: the updates it posts after the record's first message, any of ResolvedAndReopened, StateChanges, SeverityChanges, PublicNotes, PrivateNotes, Owners, Paging, EpisodeMembers, Postmortem, AI, Edits, CustomFields, Reminders and Other - for an on-call policy rule, any of OnCallNow, OnCallNext, NobodyOnCall, Overrides, PolicyChanges and Owners; for a security case rule, any of PriorityReached, Assigned, ApprovalRequested and ResponseTargets - as its kind of rule allows. null posts every update, and so does leaving it out, except on an on-call policy rule: a new one that leaves it out posts the key updates, OnCallNow and NobodyOnCall, as one made in the Dashboard does, and an edit that leaves it out keeps the updates the rule posts. The first message, and a video call's Join button, are posted whatever it lists. An on-call policy rule that posts OnCallNow can set shouldMentionWhoIsOnCall to true: its hand-off line then mentions the person who comes on call, where they connected their chat account. Owner and On-Call Policy conditions in filters name people, teams and policies of this project by ID. A JSON value: read it with `jsondecode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
