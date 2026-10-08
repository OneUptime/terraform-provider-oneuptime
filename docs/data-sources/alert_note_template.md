---
page_title: "oneuptime_alert_note_template Data Source - oneuptime"
subcategory: "Alerts"
description: |-
  Manage alert note templates for your project
---

# oneuptime_alert_note_template (Data Source)

Manage alert note templates for your project

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one alert note template may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_alert_note_template" "example" {
  note = "example-note"
}

# Or by id:
data "oneuptime_alert_note_template" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `note` (String) Note template for public or private notes. This is in markdown.
- `template_description` (String) Description of the Alert Template.
- `template_name` (String) Name of the Alert Template.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
