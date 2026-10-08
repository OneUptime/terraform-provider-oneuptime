---
page_title: "oneuptime_proxmox_cluster_feed Resource - oneuptime"
subcategory: "Other"
description: |-
  Log of everything that happened to this Proxmox cluster - creation, updates, owner changes and the rules that made them.
---

# oneuptime_proxmox_cluster_feed (Resource)

Log of everything that happened to this Proxmox cluster - creation, updates, owner changes and the rules that made them.

## Example Usage

```terraform
resource "oneuptime_proxmox_cluster_feed" "example" {
  proxmox_cluster_id              = oneuptime_proxmox_cluster.example.id
  feed_info_in_markdown           = "# Heading\n\nThis is **markdown** content"
  proxmox_cluster_feed_event_type = "Example short text"
  display_color                   = "#ff0000"
}
```

## Schema

### Required

- `display_color` (String) Display color for this feed item.
- `feed_info_in_markdown` (String) Log of the Proxmox cluster change in Markdown.
- `proxmox_cluster_feed_event_type` (String) Proxmox Cluster Feed Event.
- `proxmox_cluster_id` (String) Relation to Proxmox Cluster ID in which this resource belongs. The ID of a `oneuptime_proxmox_cluster`.

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

Import an existing proxmox cluster feed by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_proxmox_cluster_feed.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_proxmox_cluster_feed.example <id>
```
