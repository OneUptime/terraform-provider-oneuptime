---
page_title: "oneuptime_vcenter_owner_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure rules for automatically assigning owner users and teams when matching vCenters are created
---

# oneuptime_vcenter_owner_rule (Resource)

Configure rules for automatically assigning owner users and teams when matching vCenters are created

~> **Renamed:** this resource was called `oneuptime_v_center_owner_rule` before. The old name still works, but is deprecated. To switch, rename the resource in your configuration and add a `moved` block, so Terraform keeps the existing vcenter owner rule:

```terraform
moved {
  from = oneuptime_v_center_owner_rule.example
  to   = oneuptime_vcenter_owner_rule.example
}
```

## Example Usage

```terraform
resource "oneuptime_vcenter_owner_rule" "example" {
  name        = "Example vcenter owner rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this vCenter owner rule.

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this vCenter owner rule.
- `is_enabled` (Boolean) Whether this rule is enabled. Defaults to `true`.
- `notify_owners` (Boolean) Send notifications to owner users and teams when they are added by this rule. Defaults to `true`.
- `owner_teams` (Set of String) Teams to add as owners on the vCenter when this rule matches. IDs of `oneuptime_team` resources.
- `owner_users` (Set of String) Users to add as owners on the vCenter when this rule matches. IDs of `oneuptime_user` records.
- `vmware_v_center_description_pattern` (String) Regex (case-insensitive) matched against the vCenter description. Leave empty to match any description.
- `vmware_v_center_labels` (Set of String) Only trigger for vCenters that have at least one of these labels. Leave empty to match regardless of labels. IDs of `oneuptime_label` resources.
- `vmware_v_center_name_pattern` (String) Regex (case-insensitive) matched against the vCenter name. Leave empty to match any name.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing vcenter owner rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_vcenter_owner_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_vcenter_owner_rule.example <id>
```
