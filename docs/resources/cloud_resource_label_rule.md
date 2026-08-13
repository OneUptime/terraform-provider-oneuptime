---
page_title: "oneuptime_cloud_resource_label_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Rules for automatically attaching labels to cloud resources when matching resources are created.
---

# oneuptime_cloud_resource_label_rule (Resource)

Rules for automatically attaching labels to cloud resources when matching resources are created.

## Example Usage

```terraform
resource "oneuptime_cloud_resource_label_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `match_labels` (Set) Only trigger for resources that already have at least one of these labels. Leave empty to match regardless of labels...
- `name_regex_pattern` (String) Regex (case-insensitive) matched against the resource name. Leave empty to match any name...
- `description_regex_pattern` (String) Regex (case-insensitive) matched against the resource description. Leave empty to match any description...
- `labels_to_add` (Set) Labels to attach to the resource when this rule matches. Already-attached labels are not duplicated...
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
terraform import oneuptime_cloud_resource_label_rule.example <id>
```
