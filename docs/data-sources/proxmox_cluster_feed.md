---
page_title: "oneuptime_proxmox_cluster_feed Data Source - oneuptime"
subcategory: "Other"
description: |-
  Log of everything that happened to this Proxmox cluster - creation, updates, owner changes and the rules that made them.
---

# oneuptime_proxmox_cluster_feed (Data Source)

Log of everything that happened to this Proxmox cluster - creation, updates, owner changes and the rules that made them.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one proxmox cluster feed may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_proxmox_cluster_feed" "example" {
  proxmox_cluster_id = oneuptime_proxmox_cluster.example.id
}

# Or by id:
data "oneuptime_proxmox_cluster_feed" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `feed_info_in_markdown` (String) Log of the Proxmox cluster change in Markdown.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `more_information_in_markdown` (String) More information in Markdown.
- `proxmox_cluster_feed_event_type` (String) Proxmox Cluster Feed Event.
- `proxmox_cluster_id` (String) Relation to Proxmox Cluster ID in which this resource belongs. The ID of a `oneuptime_proxmox_cluster`.
- `user_id` (String) User who this feed belongs to (if this feed belongs to a User). The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `display_color` (String) Display color for this feed item.
- `posted_at` (String) Date and time when the feed was posted.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
