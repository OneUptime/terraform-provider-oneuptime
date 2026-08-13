---
page_title: "oneuptime_alert_state Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Manage alert states for your project (Created, Acknowledged for example). Add / edit or remove states.
---

# oneuptime_alert_state (Resource)

Manage alert states for your project (Created, Acknowledged for example). Add / edit or remove states.

## Example Usage

```terraform
resource "oneuptime_alert_state" "example" {
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
- `is_created_state` (Bool) Is it the created state of the alert?..
- `is_acknowledged_state` (Bool) Is it the acknowledged state of the alert?..
- `is_resolved_state` (Bool) Is it the resolved state of the alert?..
- `order` (Number) Order / Priority of this resource..

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
terraform import oneuptime_alert_state.example <id>
```
