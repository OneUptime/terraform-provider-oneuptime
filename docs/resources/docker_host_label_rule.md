---
page_title: "oneuptime_docker_host_label_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure rules for automatically attaching labels to Docker hosts when matching Docker hosts are created
---

# oneuptime_docker_host_label_rule (Resource)

Configure rules for automatically attaching labels to Docker hosts when matching Docker hosts are created

## Example Usage

```terraform
resource "oneuptime_docker_host_label_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this Docker host label rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this Docker host label rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `docker_host_labels` (Set) Only trigger for Docker hosts that already have at least one of these labels. Leave empty to match regardless of labels...
- `docker_host_name_pattern` (String) Regex (case-insensitive) matched against the Docker host name. Leave empty to match any name...
- `docker_host_description_pattern` (String) Regex (case-insensitive) matched against the Docker host description. Leave empty to match any description...
- `labels_to_add` (Set) Labels to attach to the Docker host when this rule matches. Already-attached labels are not duplicated...
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
terraform import oneuptime_docker_host_label_rule.example <id>
```
