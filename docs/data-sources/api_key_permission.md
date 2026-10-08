---
page_title: "oneuptime_api_key_permission Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  Permissions for your API Keys
---

# oneuptime_api_key_permission (Data Source)

Permissions for your API Keys

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one api key permission may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_api_key_permission" "example" {
  api_key_id = oneuptime_api_key.example.id
}

# Or by id:
data "oneuptime_api_key_permission" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `api_key_id` (String) ID of API Key resource in which this object belongs. The ID of a `oneuptime_api_key`.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_block_permission` (Boolean) Permissions - Create: [Project Owner, Project Admin, Edit API Key Permissions], Read: [Project Owner, Project Admin, Read API Key], Update: [Project Owner, Project Admin, Edit API Key Permissions]

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `labels` (Set of String) Relation to Labels Array where this permission is scoped at. IDs of `oneuptime_label` resources.
- `permission` (String) Permission. You can find list of permissions on the Permissions page. A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
