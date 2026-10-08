---
page_title: "oneuptime_ceph_cluster_feed Data Source - oneuptime"
subcategory: "Other"
description: |-
  Log of everything that happened to this Ceph cluster - creation, updates, owner changes and the rules that made them.
---

# oneuptime_ceph_cluster_feed (Data Source)

Log of everything that happened to this Ceph cluster - creation, updates, owner changes and the rules that made them.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one ceph cluster feed may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_ceph_cluster_feed" "example" {
  ceph_cluster_id = oneuptime_ceph_cluster.example.id
}

# Or by id:
data "oneuptime_ceph_cluster_feed" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `ceph_cluster_feed_event_type` (String) Ceph Cluster Feed Event.
- `ceph_cluster_id` (String) Relation to Ceph Cluster ID in which this resource belongs. The ID of a `oneuptime_ceph_cluster`.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `feed_info_in_markdown` (String) Log of the Ceph cluster change in Markdown.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `more_information_in_markdown` (String) More information in Markdown.
- `user_id` (String) User who this feed belongs to (if this feed belongs to a User). The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `display_color` (String) Display color for this feed item.
- `posted_at` (String) Date and time when the feed was posted.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
