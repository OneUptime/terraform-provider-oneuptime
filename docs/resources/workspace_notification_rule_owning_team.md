---
page_title: "oneuptime_workspace_notification_rule_owning_team Resource - oneuptime"
subcategory: "Slack, Teams & Video Calls"
description: |-
  Teams that own a workspace notification rule: only their members, and project admins, can change or delete it.
---

# oneuptime_workspace_notification_rule_owning_team (Resource)

Teams that own a workspace notification rule: only their members, and project admins, can change or delete it.

## Example Usage

```terraform
resource "oneuptime_workspace_notification_rule_owning_team" "example" {
  team_id                        = oneuptime_team.example.id
  workspace_notification_rule_id = oneuptime_workspace_notification_rule.example.id
}
```

## Schema

### Required

- `team_id` (String) ID of the team that owns the rule. The ID of a `oneuptime_team`.
- `workspace_notification_rule_id` (String) ID of the workspace notification rule the team owns. The ID of a `oneuptime_workspace_notification_rule`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing workspace notification rule owning team by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_workspace_notification_rule_owning_team.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_workspace_notification_rule_owning_team.example <id>
```
