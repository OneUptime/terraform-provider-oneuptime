---
page_title: "oneuptime_incoming_call_policy_team_owner Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Add teams as owners to your incoming call policies.
---

# oneuptime_incoming_call_policy_team_owner (Resource)

Add teams as owners to your incoming call policies.

## Example Usage

```terraform
resource "oneuptime_incoming_call_policy_team_owner" "example" {
  team_id                 = oneuptime_team.example.id
  incoming_call_policy_id = oneuptime_incoming_call_policy.example.id
}
```

## Schema

### Required

- `incoming_call_policy_id` (String) ID of the incoming call policy the team owns. The ID of a `oneuptime_incoming_call_policy`.
- `team_id` (String) ID of the team that owns the incoming call policy. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incoming call policy team owner by its incoming call policy and its `team_id`, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incoming_call_policy_team_owner.example
  id = "<incoming_call_policy_id>/<team_id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incoming_call_policy_team_owner.example <incoming_call_policy_id>/<team_id>
```

Its own ID works too.
