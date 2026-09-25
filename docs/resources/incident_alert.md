---
page_title: "oneuptime_incident_alert Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Alerts linked to an incident. Link alerts to an incident to track them as part of that incident's response.
---

# oneuptime_incident_alert (Resource)

Alerts linked to an incident. Link alerts to an incident to track them as part of that incident's response.

## Example Usage

```terraform
resource "oneuptime_incident_alert" "example" {
  incident_id = "123e4567-e89b-12d3-a456-426614174000"
  alert_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `incident_id` (String) A unique identifier for an object, represented as a UUID..
- `alert_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
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
terraform import oneuptime_incident_alert.example <id>
```
