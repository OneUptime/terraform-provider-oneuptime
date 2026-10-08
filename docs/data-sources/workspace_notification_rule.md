---
page_title: "oneuptime_workspace_notification_rule Data Source - oneuptime"
subcategory: "Other"
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
- `notification_rule` (String) Notification Rules for the Workspace. A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
