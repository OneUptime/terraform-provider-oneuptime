---
page_title: "oneuptime_incident_state Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Manage incident states for your project (Created, Acknowledged for example). Add / edit or remove states.
---

# oneuptime_incident_state (Data Source)

Manage incident states for your project (Created, Acknowledged for example). Add / edit or remove states. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_incident_state" "by_name" {
  name = "example-incident_state"
}

data "oneuptime_incident_state" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `description` (String) Friendly description that will help you remember.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `color` (String) Color object. Computed.
- `is_created_state` (Bool) Is it the created state of the incident?.. Computed.
- `is_acknowledged_state` (Bool) Is it the acknowledged state of the incident?.. Computed.
- `is_resolved_state` (Bool) Is it the resolved state of the incident?.. Computed.
- `order` (Number) Order / Priority of this resource.. Computed.
