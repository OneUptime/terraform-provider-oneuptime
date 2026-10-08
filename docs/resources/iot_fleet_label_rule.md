---
page_title: "oneuptime_iot_fleet_label_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure rules for automatically attaching labels to IoT fleets when matching IoT fleets are created
---

# oneuptime_iot_fleet_label_rule (Resource)

Configure rules for automatically attaching labels to IoT fleets when matching IoT fleets are created

~> **Renamed:** this resource was called `oneuptime_io_t_fleet_label_rule` before. The old name still works, but is deprecated. To switch, rename the resource in your configuration and add a `moved` block, so Terraform keeps the existing iot fleet label rule:

```terraform
moved {
  from = oneuptime_io_t_fleet_label_rule.example
  to   = oneuptime_iot_fleet_label_rule.example
}
```

## Example Usage

```terraform
resource "oneuptime_iot_fleet_label_rule" "example" {
  name        = "Example iot fleet label rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this IoT fleet label rule.

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this IoT fleet label rule.
- `iot_fleet_description_pattern` (String) Regex (case-insensitive) matched against the IoT fleet description. Leave empty to match any description.
- `iot_fleet_labels` (Set of String) Only trigger for IoT fleets that already have at least one of these labels. Leave empty to match regardless of labels. IDs of `oneuptime_label` resources.
- `iot_fleet_name_pattern` (String) Regex (case-insensitive) matched against the IoT fleet name. Leave empty to match any name.
- `is_enabled` (Boolean) Whether this rule is enabled. Defaults to `true`.
- `labels_to_add` (Set of String) Labels to attach to the IoT fleet when this rule matches. Already-attached labels are not duplicated. IDs of `oneuptime_label` resources.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing iot fleet label rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_iot_fleet_label_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_iot_fleet_label_rule.example <id>
```
