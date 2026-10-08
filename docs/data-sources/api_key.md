---
page_title: "oneuptime_api_key Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  Manage API Keys for your project
---

# oneuptime_api_key (Data Source)

Manage API Keys for your project

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one api key may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_api_key" "example" {
  name = "Example api key"
}

# Or by id:
data "oneuptime_api_key" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `api_key` (String) Secret API Key.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `name` (String) Any friendly name of this object.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `expires_at` (String) Date and Time when this API Key expires.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
