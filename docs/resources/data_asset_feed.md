---
page_title: "oneuptime_data_asset_feed Resource - oneuptime"
subcategory: "Other"
description: |-
  The timeline of a data asset: when it was found, which columns changed, and when it was dropped or restored.
---

# oneuptime_data_asset_feed (Resource)

The timeline of a data asset: when it was found, which columns changed, and when it was dropped or restored.

## Example Usage

```terraform
resource "oneuptime_data_asset_feed" "example" {
  data_asset_id              = oneuptime_data_asset.example.id
  feed_info_in_markdown      = "# Heading\n\nThis is **markdown** content"
  data_asset_feed_event_type = "AssetDiscovered"
  display_color              = "#ff0000"
}
```

## Schema

### Required

- `data_asset_feed_event_type` (String) What kind of change this entry records. Allowed values: `AssetDiscovered`, `AssetDropped`, `AssetRestored`, `ColumnAdded`, `ColumnRemoved`, `ColumnRenamed`, `ColumnTypeChanged`, `ColumnNullabilityChanged`, `DescriptionUpdated`, `OwnerUserAdded`, `OwnerUserRemoved`, `OwnerTeamAdded`, `OwnerTeamRemoved`, `OwnerRuleExecuted`, `LabelRuleExecuted`.
- `data_asset_id` (String) ID of the data asset this entry is about. The ID of a `oneuptime_data_asset`.
- `display_color` (String) Display color for this feed item.
- `feed_info_in_markdown` (String) What happened to the data asset, in Markdown.

### Optional

- `more_information_in_markdown` (String) More information in Markdown.
- `posted_at` (String) Date and time when the feed was posted.
- `user_email` (String) The email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) ID of the user this entry is about or was done by (if a person was involved). The ID of a `oneuptime_user` (see the data source). Or name them by email address, in `user_email`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing data asset feed by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_data_asset_feed.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_data_asset_feed.example <id>
```
