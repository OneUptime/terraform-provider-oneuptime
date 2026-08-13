---
page_title: "oneuptime_scheduled_maintenance_note_template Data Source - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Manage scheduled maintenance note templates for your project
---

# oneuptime_scheduled_maintenance_note_template (Data Source)

Manage scheduled maintenance note templates for your project Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_scheduled_maintenance_note_template" "by_name" {
  name = "example-scheduled_maintenance_note_template"
}

data "oneuptime_scheduled_maintenance_note_template" "by_id" {
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
- `note` (String) Note template for public or private notes. This is in markdown... Computed.
- `template_name` (String) Name of the Incident Template.. Computed.
- `template_description` (String) Description of the Incident Template.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
