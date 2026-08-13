---
page_title: "oneuptime_host_label_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure rules for automatically attaching labels to hosts when matching hosts are created
---

# oneuptime_host_label_rule (Resource)

Configure rules for automatically attaching labels to hosts when matching hosts are created

## Example Usage

```terraform
resource "oneuptime_host_label_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this host label rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this host label rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `host_labels` (Set) Only trigger for hosts that already have at least one of these labels. Leave empty to match regardless of labels...
- `host_name_pattern` (String) Regex (case-insensitive) matched against the host name. Leave empty to match any name...
- `host_description_pattern` (String) Regex (case-insensitive) matched against the host description. Leave empty to match any description...
- `labels_to_add` (Set) Labels to attach to the host when this rule matches. Already-attached labels are not duplicated...
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
terraform import oneuptime_host_label_rule.example <id>
```
