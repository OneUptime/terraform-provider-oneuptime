---
page_title: "oneuptime_vcenter_feed Resource - oneuptime"
subcategory: "Other"
description: |-
  Log of everything that happened to this vCenter - creation, updates, owner changes and the rules that made them.
---

# oneuptime_vcenter_feed (Resource)

Log of everything that happened to this vCenter - creation, updates, owner changes and the rules that made them.

~> **Renamed:** this resource was called `oneuptime_v_center_feed` before. The old name still works, but is deprecated. To switch, rename the resource in your configuration and add a `moved` block, so Terraform keeps the existing vcenter feed:

```terraform
moved {
  from = oneuptime_v_center_feed.example
  to   = oneuptime_vcenter_feed.example
}
```

## Example Usage

```terraform
resource "oneuptime_vcenter_feed" "example" {
  vmware_v_center_id              = oneuptime_vcenter.example.id
  feed_info_in_markdown           = "# Heading\n\nThis is **markdown** content"
  vmware_v_center_feed_event_type = "Example short text"
  display_color                   = "#ff0000"
}
```

## Schema

### Required

- `display_color` (String) Display color for this feed item.
- `feed_info_in_markdown` (String) Log of the vCenter change in Markdown.
- `vmware_v_center_feed_event_type` (String) vCenter Feed Event.
- `vmware_v_center_id` (String) Relation to vCenter ID in which this resource belongs. The ID of a `oneuptime_vcenter`.

### Optional

- `more_information_in_markdown` (String) More information in Markdown.
- `posted_at` (String) Date and time when the feed was posted.
- `user_id` (String) User who this feed belongs to (if this feed belongs to a User). The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing vcenter feed by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_vcenter_feed.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_vcenter_feed.example <id>
```
