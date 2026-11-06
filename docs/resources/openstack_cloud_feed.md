---
page_title: "oneuptime_openstack_cloud_feed Resource - oneuptime"
subcategory: "Infrastructure"
description: |-
  Log of everything that happened to this OpenStack cloud - creation, updates, owner changes and the rules that made them.
---

# oneuptime_openstack_cloud_feed (Resource)

Log of everything that happened to this OpenStack cloud - creation, updates, owner changes and the rules that made them.

~> **Renamed:** this resource was called `oneuptime_open_stack_cloud_feed` before. The old name still works, but is deprecated. To switch, rename the resource in your configuration and add a `moved` block, so Terraform keeps the existing openstack cloud feed:

```terraform
moved {
  from = oneuptime_open_stack_cloud_feed.example
  to   = oneuptime_openstack_cloud_feed.example
}
```

## Example Usage

```terraform
resource "oneuptime_openstack_cloud_feed" "example" {
  open_stack_cloud_id              = oneuptime_openstack_cloud.example.id
  feed_info_in_markdown            = "# Heading\n\nThis is **markdown** content"
  open_stack_cloud_feed_event_type = "OpenStackCloudCreated"
  display_color                    = "#ff0000"
}
```

## Schema

### Required

- `display_color` (String) Display color for this feed item.
- `feed_info_in_markdown` (String) Log of the OpenStack cloud change in Markdown.
- `open_stack_cloud_feed_event_type` (String) OpenStack Cloud Feed Event. Allowed values: `OpenStackCloudCreated`, `OpenStackCloudUpdated`, `OpenStackCloudArchived`, `OpenStackCloudRestored`, `OwnerUserAdded`, `OwnerUserRemoved`, `OwnerTeamAdded`, `OwnerTeamRemoved`, `OwnerRuleExecuted`, `LabelRuleExecuted`, `LifecycleEvent`, `AiFix`.
- `open_stack_cloud_id` (String) Relation to OpenStack Cloud ID in which this resource belongs. The ID of a `oneuptime_openstack_cloud`.

### Optional

- `more_information_in_markdown` (String) More information in Markdown.
- `posted_at` (String) Date and time when the feed was posted.
- `user_email` (String) The email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) User who this feed belongs to (if this feed belongs to a User). The ID of a `oneuptime_user` (see the data source). Or name them by email address, in `user_email`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing openstack cloud feed by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_openstack_cloud_feed.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_openstack_cloud_feed.example <id>
```
