---
page_title: "oneuptime_openstack_cloud_label_rule Resource - oneuptime"
subcategory: "Infrastructure"
description: |-
  Configure rules for automatically attaching labels to OpenStack clouds when matching OpenStack clouds are created
---

# oneuptime_openstack_cloud_label_rule (Resource)

Configure rules for automatically attaching labels to OpenStack clouds when matching OpenStack clouds are created

~> **Renamed:** this resource was called `oneuptime_open_stack_cloud_label_rule` before. The old name still works, but is deprecated. To switch, rename the resource in your configuration and add a `moved` block, so Terraform keeps the existing openstack cloud label rule:

```terraform
moved {
  from = oneuptime_open_stack_cloud_label_rule.example
  to   = oneuptime_openstack_cloud_label_rule.example
}
```

## Example Usage

```terraform
resource "oneuptime_openstack_cloud_label_rule" "example" {
  name        = "Example openstack cloud label rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this OpenStack cloud label rule.

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this OpenStack cloud label rule.
- `is_enabled` (Boolean) Whether this rule is enabled. Defaults to `true`.
- `labels_to_add` (Set of String) Labels to attach to the OpenStack cloud when this rule matches. Already-attached labels are not duplicated. IDs of `oneuptime_label` resources.
- `open_stack_cloud_description_pattern` (String) Regex (case-insensitive) matched against the OpenStack cloud description. Leave empty to match any description.
- `open_stack_cloud_labels` (Set of String) Only trigger for OpenStack clouds that already have at least one of these labels. Leave empty to match regardless of labels. IDs of `oneuptime_label` resources.
- `open_stack_cloud_name_pattern` (String) Regex (case-insensitive) matched against the OpenStack cloud name. Leave empty to match any name.
- `open_stack_domain_id` (String) Keystone's ID for one domain. Set, this rule adds its labels to every alert and incident that reaches a project in that domain instead of to the cloud.
- `open_stack_project_id` (String) Keystone's ID for one project. Set, this rule adds its labels to every alert and incident that reaches that project instead of to the cloud.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing openstack cloud label rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_openstack_cloud_label_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_openstack_cloud_label_rule.example <id>
```
