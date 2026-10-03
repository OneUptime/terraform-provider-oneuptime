---
page_title: "oneuptime_incident_severity Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Manage incident severity for your project (Created, Acknowledged for example). Add / edit or remove severities.
---

# oneuptime_incident_severity (Resource)

Manage incident severity for your project (Created, Acknowledged for example). Add / edit or remove severities.

## Example Usage

```terraform
resource "oneuptime_incident_severity" "example" {
  name = "Example short text"
  color = jsonencode({
    "_type": "Color",
    "value": "#ff0000"
  })
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object..
- `color` (String) Color object.

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `order` (Number) Where this severity ranks among the project's incident severities: 1 is the most severe. A new severity without a number goes to the end of the list. Setting a number moves the severity to that place, and the ones in between shift by one. In the dashboard, drag the rows to reorder them...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_incident_severity.example <id>
```
