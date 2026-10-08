---
page_title: "oneuptime_ai_agent_owner_team Resource - oneuptime"
subcategory: "Other"
description: |-
  Add teams as owners to your AI agents.
---

# oneuptime_ai_agent_owner_team (Resource)

Add teams as owners to your AI agents.

## Example Usage

```terraform
resource "oneuptime_ai_agent_owner_team" "example" {
  team_id     = oneuptime_team.example.id
  ai_agent_id = oneuptime_ai_agent.example.id
}
```

## Schema

### Required

- `ai_agent_id` (String) ID of your OneUptime AI Agent in which this object belongs. The ID of a `oneuptime_ai_agent`.
- `team_id` (String) ID of your OneUptime Team in which this object belongs. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing ai agent owner team by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_ai_agent_owner_team.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_ai_agent_owner_team.example <id>
```
