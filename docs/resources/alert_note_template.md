---
page_title: "oneuptime_alert_note_template Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Manage alert note templates for your project
---

# oneuptime_alert_note_template (Resource)

Manage alert note templates for your project

## Example Usage

```terraform
resource "oneuptime_alert_note_template" "example" {
  template_name = "Example short text"
  template_description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `template_name` (String) Name of the Alert Template..
- `template_description` (String) Description of the Alert Template..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `note` (String) Note template for public or private notes. This is in markdown...
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
terraform import oneuptime_alert_note_template.example <id>
```
