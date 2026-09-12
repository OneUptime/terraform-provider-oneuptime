---
page_title: "oneuptime_v_center_label_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure rules for automatically attaching labels to vCenters when matching vCenters are created
---

# oneuptime_v_center_label_rule (Resource)

Configure rules for automatically attaching labels to vCenters when matching vCenters are created

## Example Usage

```terraform
resource "oneuptime_v_center_label_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this vCenter label rule..

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource...
- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this vCenter label rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `vmware_v_center_labels` (Set) Only trigger for vCenters that already have at least one of these labels. Leave empty to match regardless of labels...
- `vmware_v_center_name_pattern` (String) Regex (case-insensitive) matched against the vCenter name. Leave empty to match any name...
- `vmware_v_center_description_pattern` (String) Regex (case-insensitive) matched against the vCenter description. Leave empty to match any description...
- `labels_to_add` (Set) Labels to attach to the vCenter when this rule matches. Already-attached labels are not duplicated...
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
terraform import oneuptime_v_center_label_rule.example <id>
```
