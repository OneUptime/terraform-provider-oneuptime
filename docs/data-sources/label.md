---
page_title: "oneuptime_label Data Source - oneuptime"
subcategory: "Organization"
description: |-
  Organize resources for your project by using labels / tags.
---

# oneuptime_label (Data Source)

Organize resources for your project by using labels / tags.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one label may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_label" "example" {
  name = "Example label"
}

# Or by id:
data "oneuptime_label" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `name` (String) Any friendly name of this object.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `color` (String) Color of this resource in Hex (#32a852 for example).
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
