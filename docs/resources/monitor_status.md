---
page_title: "oneuptime_monitor_status Resource - oneuptime"
subcategory: "Monitors"
description: |-
  Manage monitor status in your project. Monitor Status are Operational, Degraded and Offline for example. Add custom status like Monitoring or more.
---

# oneuptime_monitor_status (Resource)

Manage monitor status in your project. Monitor Status are Operational, Degraded and Offline for example. Add custom status like Monitoring or more.

## Example Usage

```terraform
resource "oneuptime_monitor_status" "example" {
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
- `is_operational_state` (Bool) Is this monitor in operational state?..
- `is_offline_state` (Bool) Is this monitor in offline state?..
- `priority` (Number) Order / Priority of this status. Behaves like an insertion slot: creating a status with priority P shifts every existing status with priority >= P up by one, and priority cannot be changed after creation (delete and recreate instead). When managing statuses declaratively, use high, gapped values (e.g. 101, 102, 103) created in ascending order so existing statuses are never shifted...

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
terraform import oneuptime_monitor_status.example <id>
```
