---
page_title: "oneuptime_podman_host_label_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure rules for automatically attaching labels to Podman hosts when matching Podman hosts are created
---

# oneuptime_podman_host_label_rule (Resource)

Configure rules for automatically attaching labels to Podman hosts when matching Podman hosts are created

## Example Usage

```terraform
resource "oneuptime_podman_host_label_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this Podman host label rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this Podman host label rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `podman_host_labels` (Set) Only trigger for Podman hosts that already have at least one of these labels. Leave empty to match regardless of labels...
- `podman_host_name_pattern` (String) Regex (case-insensitive) matched against the Podman host name. Leave empty to match any name...
- `podman_host_description_pattern` (String) Regex (case-insensitive) matched against the Podman host description. Leave empty to match any description...
- `labels_to_add` (Set) Labels to attach to the Podman host when this rule matches. Already-attached labels are not duplicated...
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
terraform import oneuptime_podman_host_label_rule.example <id>
```
