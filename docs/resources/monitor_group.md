---
page_title: "oneuptime_monitor_group Resource - oneuptime"
subcategory: "Monitors"
description: |-
  Monitor Groups are a way to organize your monitors into groups. You can create as many groups as you want and add as many monitors as you want to each group.
---

# oneuptime_monitor_group (Resource)

Monitor Groups are a way to organize your monitors into groups. You can create as many groups as you want and add as many monitors as you want to each group.

## Example Usage

```terraform
resource "oneuptime_monitor_group" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Any friendly name for this monitor group..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `labels` (Set) Relation to Labels Array where this object is categorized in...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_monitor_group.example <id>
```
