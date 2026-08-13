---
page_title: "oneuptime_monitor_owner_rule Resource - oneuptime"
subcategory: "Monitors"
description: |-
  Configure rules for automatically assigning owner users and teams when matching monitors are created
---

# oneuptime_monitor_owner_rule (Resource)

Configure rules for automatically assigning owner users and teams when matching monitors are created

## Example Usage

```terraform
resource "oneuptime_monitor_owner_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this monitor owner rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this monitor owner rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `notify_owners` (Bool) Send notifications to owner users and teams when they are added by this rule..
- `monitor_labels` (Set) Only trigger for monitors that have at least one of these labels. Leave empty to match regardless of labels...
- `monitor_name_pattern` (String) Regex (case-insensitive) matched against the monitor name. Leave empty to match any name...
- `monitor_description_pattern` (String) Regex (case-insensitive) matched against the monitor description. Leave empty to match any description...
- `owner_users` (Set) Users to add as owners on the monitor when this rule matches...
- `owner_teams` (Set) Teams to add as owners on the monitor when this rule matches...
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
terraform import oneuptime_monitor_owner_rule.example <id>
```
