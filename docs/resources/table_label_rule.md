---
page_title: "oneuptime_table_label_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Rules for automatically attaching labels to tables when matching resources are created.
---

# oneuptime_table_label_rule (Resource)

Rules for automatically attaching labels to tables when matching resources are created.

## Example Usage

```terraform
resource "oneuptime_table_label_rule" "example" {
  name        = "Example table label rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this rule.

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `database_regex_pattern` (String) Case-insensitive regex matched against the name of the database the table is in.
- `description` (String) Description of this rule.
- `is_enabled` (Boolean) Whether this rule is enabled. Defaults to `true`.
- `labels_to_add` (Set of String) Labels to attach to the resource when this rule matches. Already-attached labels are not duplicated. IDs of `oneuptime_label` resources.
- `match_labels` (Set of String) Only tables with one of these labels. Leave empty to skip the filter. IDs of `oneuptime_label` resources.
- `name_regex_pattern` (String) Case-insensitive regex matched against the table's name.
- `schema_regex_pattern` (String) Case-insensitive regex matched against the table's schema.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing table label rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_table_label_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_table_label_rule.example <id>
```
