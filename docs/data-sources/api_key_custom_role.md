---
page_title: "oneuptime_api_key_custom_role Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  The custom roles an API key is given. The key holds every permission of the role except those with the Owned scope: a key owns nothing.
---

# oneuptime_api_key_custom_role (Data Source)

The custom roles an API key is given. The key holds every permission of the role except those with the Owned scope: a key owns nothing.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one api key custom role may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_api_key_custom_role" "example" {
  api_key_id = oneuptime_api_key.example.id
}

# Or by id:
data "oneuptime_api_key_custom_role" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `api_key_id` (String) ID of the API key given the role. The ID of a `oneuptime_api_key`.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `custom_role_id` (String) ID of the custom role given to the API key. The ID of a `oneuptime_custom_role`.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
