---
page_title: "oneuptime_incident_member Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Assign users with specific roles to incidents.
---

# oneuptime_incident_member (Resource)

Assign users with specific roles to incidents.

## Example Usage

```terraform
resource "oneuptime_incident_member" "example" {
  user_id = "123e4567-e89b-12d3-a456-426614174000"
  incident_id = "123e4567-e89b-12d3-a456-426614174000"
  incident_role_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `user_id` (String) A unique identifier for an object, represented as a UUID..
- `incident_id` (String) A unique identifier for an object, represented as a UUID..
- `incident_role_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `notes` (String) Assignment context or notes..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_incident_member.example <id>
```
