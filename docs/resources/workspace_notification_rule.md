---
page_title: "oneuptime_workspace_notification_rule Resource - oneuptime"
subcategory: "Slack, Teams & Video Calls"
description: |-
  Notification Rule for Third Party Workspaces
---

# oneuptime_workspace_notification_rule (Resource)

Notification Rule for Third Party Workspaces

## Example Usage

```terraform
resource "oneuptime_workspace_notification_rule" "example" {
  name           = "Example workspace notification rule"
  event_type     = "Incident"
  workspace_type = "Slack"
  description    = "Managed by Terraform"
}
```

## Schema

### Required

- `event_type` (String) Event Type for the Workspace like Incident Created, Monitor Status Updated, etc. Allowed values: `Incident`, `Monitor`, `Alert`, `Alert Episode`, `Incident Episode`, `Scheduled Maintenance`, `On-Call Duty Policy`, `Security Case`.
- `name` (String) Name of the Notification Rule.
- `workspace_type` (String) Type of Workspace - slack, microsoft teams etc. Allowed values: `Slack`, `MicrosoftTeams`.

### Optional

- `description` (String) Description of the Notification Rule.
- `notification_rule` (String) Notification Rules for the Workspace. A JSON value: write it with `jsonencode()`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing workspace notification rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_workspace_notification_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_workspace_notification_rule.example <id>
```
