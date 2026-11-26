---
page_title: "oneuptime_data_asset_feed Data Source - oneuptime"
subcategory: "Other"
description: |-
  The timeline of a data asset: when it was found, which columns changed, and when it was dropped or restored.
---

# oneuptime_data_asset_feed (Data Source)

The timeline of a data asset: when it was found, which columns changed, and when it was dropped or restored.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data asset feed may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_asset_feed" "example" {
  last_modified_by_user_id = data.oneuptime_user.example.id
  user_email               = "alice@example.com"
}

# Or by id:
data "oneuptime_data_asset_feed" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `data_asset_feed_event_type` (String) What kind of change this entry records.
- `data_asset_id` (String) ID of the data asset this entry is about. The ID of a `oneuptime_data_asset`.
- `feed_info_in_markdown` (String) What happened to the data asset, in Markdown.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `more_information_in_markdown` (String) More information in Markdown.
- `user_email` (String) Look it up by the email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) ID of the user this entry is about or was done by (if a person was involved). The ID of a `oneuptime_user` (see the data source). Or look them up by email address, in `user_email`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `display_color` (String) Display color for this feed item.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `posted_at` (String) Date and time when the feed was posted.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
