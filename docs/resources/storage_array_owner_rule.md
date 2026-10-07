---
page_title: "oneuptime_storage_array_owner_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure rules for automatically assigning owner users and teams when matching storage arrays are created
---

# oneuptime_storage_array_owner_rule (Resource)

Configure rules for automatically assigning owner users and teams when matching storage arrays are created

## Example Usage

```terraform
resource "oneuptime_storage_array_owner_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this storage array owner rule..

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource...
- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this storage array owner rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `notify_owners` (Bool) Send notifications to owner users and teams when they are added by this rule..
- `storage_array_labels` (Set) Only trigger for storage arrays that have at least one of these labels. Leave empty to match regardless of labels...
- `storage_array_name_pattern` (String) Regex (case-insensitive) matched against the storage array name. Leave empty to match any name...
- `storage_array_description_pattern` (String) Regex (case-insensitive) matched against the storage array description. Leave empty to match any description...
- `owner_users` (Set) Users to add as owners on the storage array when this rule matches...
- `owner_teams` (Set) Teams to add as owners on the storage array when this rule matches...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_storage_array_owner_rule.example <id>
```
