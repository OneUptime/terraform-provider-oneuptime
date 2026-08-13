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
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this network device label rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this network device label rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `network_device_labels` (Set) Only trigger for network devices that already have at least one of these labels. Leave empty to match regardless of labels...
- `network_device_name_pattern` (String) Regex or * wildcard pattern (case-insensitive) matched against the network device name. Leave empty to match any name...
- `network_device_description_pattern` (String) Regex or * wildcard pattern (case-insensitive) matched against the network device description. Leave empty to match any description...
- `labels_to_add` (Set) Labels to attach to the network device when this rule matches. Already-attached labels are not duplicated...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_network_device_label_rule.example <id>
```
