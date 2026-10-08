---
page_title: "oneuptime_scheduled_maintenance_note_template Data Source - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Manage scheduled maintenance note templates for your project
---

# oneuptime_scheduled_maintenance_note_template (Data Source)

Manage scheduled maintenance note templates for your project

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one scheduled maintenance note template may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_scheduled_maintenance_note_template" "example" {
  note = "example-note"
}

# Or by id:
data "oneuptime_scheduled_maintenance_note_template" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `note` (String) Note template for public or private notes. This is in markdown.
- `template_description` (String) Description of the Incident Template.
- `template_name` (String) Name of the Incident Template.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
