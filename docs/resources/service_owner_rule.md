---
page_title: "oneuptime_service_owner_rule Resource - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Configure rules for automatically assigning owner users and teams when matching services are created
---

# oneuptime_service_owner_rule (Resource)

Configure rules for automatically assigning owner users and teams when matching services are created

## Example Usage

```terraform
resource "oneuptime_service_owner_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this service owner rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this service owner rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `notify_owners` (Bool) Send notifications to owner users and teams when they are added by this rule..
- `service_labels` (Set) Only trigger for services that have at least one of these labels. Leave empty to match regardless of labels...
- `service_name_pattern` (String) Regex (case-insensitive) matched against the service name. Leave empty to match any name...
- `service_description_pattern` (String) Regex (case-insensitive) matched against the service description. Leave empty to match any description...
- `owner_users` (Set) Users to add as owners on the service when this rule matches...
- `owner_teams` (Set) Teams to add as owners on the service when this rule matches...
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
terraform import oneuptime_service_owner_rule.example <id>
```
