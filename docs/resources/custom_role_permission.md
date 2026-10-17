---
page_title: "oneuptime_custom_role_permission Resource - oneuptime"
subcategory: "Other"
description: |-
  The permissions in a custom role, each with the scope and labels a team's own permission takes.
---

# oneuptime_custom_role_permission (Resource)

The permissions in a custom role, each with the scope and labels a team's own permission takes.

## Example Usage

```terraform
resource "oneuptime_custom_role_permission" "example" {
  custom_role_id = oneuptime_custom_role.example.id
}
```

## Schema

### Required

- `custom_role_id` (String) ID of the custom role this permission is part of. The ID of a `oneuptime_custom_role`.

### Optional

- `labels` (Set of String) With the scope Labels: the labels this permission is limited to. A resource carrying at least one of them is reached. IDs of `oneuptime_label` resources.
- `permission` (String) A built-in role, such as IncidentMember, or a single permission, such as ReadProjectIncident. You can find the list of permissions in the Permission Reference. ProjectOwner cannot be part of a custom role. A JSON value: write it with `jsonencode()`.
- `scope` (String) Which resources this permission reaches. One of: All (every resource in the project), Owned (resources the team given the role, or one of its members, owns), Labels (resources carrying one of the labels). Defaults to All. Defaults to `All`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing custom role permission by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_custom_role_permission.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_custom_role_permission.example <id>
```
