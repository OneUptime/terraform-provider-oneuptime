---
page_title: "oneuptime_openstack_cloud_owner_rule Resource - oneuptime"
subcategory: "Infrastructure"
description: |-
  Configure rules for automatically assigning owner users and teams when matching OpenStack clouds are created
---

# oneuptime_openstack_cloud_owner_rule (Resource)

Configure rules for automatically assigning owner users and teams when matching OpenStack clouds are created

~> **Renamed:** this resource was called `oneuptime_open_stack_cloud_owner_rule` before. The old name still works, but is deprecated. To switch, rename the resource in your configuration and add a `moved` block, so Terraform keeps the existing openstack cloud owner rule:

```terraform
moved {
  from = oneuptime_open_stack_cloud_owner_rule.example
  to   = oneuptime_openstack_cloud_owner_rule.example
}
```

## Example Usage

```terraform
resource "oneuptime_openstack_cloud_owner_rule" "example" {
  name        = "Example openstack cloud owner rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this OpenStack cloud owner rule.

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this OpenStack cloud owner rule.
- `is_enabled` (Boolean) Whether this rule is enabled. Defaults to `true`.
- `notify_owners` (Boolean) Send notifications to owner users and teams when they are added by this rule. Defaults to `true`.
- `open_stack_cloud_description_pattern` (String) Regex (case-insensitive) matched against the OpenStack cloud description. Leave empty to match any description.
- `open_stack_cloud_labels` (Set of String) Only trigger for OpenStack clouds that have at least one of these labels. Leave empty to match regardless of labels. IDs of `oneuptime_label` resources.
- `open_stack_cloud_name_pattern` (String) Regex (case-insensitive) matched against the OpenStack cloud name. Leave empty to match any name.
- `open_stack_domain_id` (String) Keystone's ID for one domain. Set, this rule adds its owners to every alert and incident that reaches a project in that domain instead of to the cloud.
- `open_stack_project_id` (String) Keystone's ID for one project. Set, this rule adds its owners to every alert and incident that reaches that project instead of to the cloud.
- `owner_teams` (Set of String) Teams to add as owners on the OpenStack cloud when this rule matches. IDs of `oneuptime_team` resources.
- `owner_user_emails` (Set of String) The email addresses of the people in `owner_users`, instead of their IDs: each one of someone in this project, or invited to it. Give one of the two.
- `owner_users` (Set of String) Users to add as owners on the OpenStack cloud when this rule matches. IDs of `oneuptime_user` records. Or name them by email address, in `owner_user_emails`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing openstack cloud owner rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_openstack_cloud_owner_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_openstack_cloud_owner_rule.example <id>
```
