---
page_title: "oneuptime_check_setting Data Source - oneuptime"
subcategory: "Data Security"
description: |-
  How a project tunes a built-in check: turn it off, or give its failures another severity.
---

# oneuptime_check_setting (Data Source)

How a project tunes a built-in check: turn it off, or give its failures another severity.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one check setting may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_check_setting" "example" {
  check_id = "example-check-id"
}

# Or by id:
data "oneuptime_check_setting" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `check_id` (String) The built-in check this setting changes.
- `created_by_user_id` (String) ID of the user who created this check setting. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether the check runs.
- `note` (String) Why the check was changed.
- `severity_override` (String) The severity every failure of the check gets, instead of the one its store's sensitivity sets.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this check setting belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
