---
page_title: "oneuptime_vcenter_label_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure rules for automatically attaching labels to vCenters when matching vCenters are created
---

# oneuptime_vcenter_label_rule (Resource)

Configure rules for automatically attaching labels to vCenters when matching vCenters are created

~> **Renamed:** this resource was called `oneuptime_v_center_label_rule` before. The old name still works, but is deprecated. To switch, rename the resource in your configuration and add a `moved` block, so Terraform keeps the existing vcenter label rule:

```terraform
moved {
  from = oneuptime_v_center_label_rule.example
  to   = oneuptime_vcenter_label_rule.example
}
```

## Example Usage

```terraform
resource "oneuptime_vcenter_label_rule" "example" {
  name        = "Example vcenter label rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this vCenter label rule.

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this vCenter label rule.
- `is_enabled` (Boolean) Whether this rule is enabled. Defaults to `true`.
- `labels_to_add` (Set of String) Labels to attach to the vCenter when this rule matches. Already-attached labels are not duplicated. IDs of `oneuptime_label` resources.
- `vmware_v_center_description_pattern` (String) Regex (case-insensitive) matched against the vCenter description. Leave empty to match any description.
- `vmware_v_center_labels` (Set of String) Only trigger for vCenters that already have at least one of these labels. Leave empty to match regardless of labels. IDs of `oneuptime_label` resources.
- `vmware_v_center_name_pattern` (String) Regex (case-insensitive) matched against the vCenter name. Leave empty to match any name.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing vcenter label rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_vcenter_label_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_vcenter_label_rule.example <id>
```
