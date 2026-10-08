---
page_title: "oneuptime_network_device_label_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure rules for automatically attaching labels to network devices when matching network devices are created
---

# oneuptime_network_device_label_rule (Resource)

Configure rules for automatically attaching labels to network devices when matching network devices are created

## Example Usage

```terraform
resource "oneuptime_network_device_label_rule" "example" {
  name        = "Example network device label rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this network device label rule.

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this network device label rule.
- `is_enabled` (Boolean) Whether this rule is enabled. Defaults to `true`.
- `labels_to_add` (Set of String) Labels to attach to the network device when this rule matches. Already-attached labels are not duplicated. IDs of `oneuptime_label` resources.
- `network_device_description_pattern` (String) Regex or * wildcard pattern (case-insensitive) matched against the network device description. Leave empty to match any description.
- `network_device_labels` (Set of String) Only trigger for network devices that already have at least one of these labels. Leave empty to match regardless of labels. IDs of `oneuptime_label` resources.
- `network_device_name_pattern` (String) Regex or * wildcard pattern (case-insensitive) matched against the network device name. Leave empty to match any name.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing network device label rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_network_device_label_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_network_device_label_rule.example <id>
```
