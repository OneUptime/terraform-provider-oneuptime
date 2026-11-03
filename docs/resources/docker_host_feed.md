---
page_title: "oneuptime_docker_host_feed Resource - oneuptime"
subcategory: "Infrastructure"
description: |-
  Log of everything that happened to this Docker host - creation, updates, owner changes and the rules that made them.
---

# oneuptime_docker_host_feed (Resource)

Log of everything that happened to this Docker host - creation, updates, owner changes and the rules that made them.

## Example Usage

```terraform
resource "oneuptime_docker_host_feed" "example" {
  docker_host_id              = oneuptime_docker_host.example.id
  feed_info_in_markdown       = "# Heading\n\nThis is **markdown** content"
  docker_host_feed_event_type = "DockerHostCreated"
  display_color               = "#ff0000"
}
```

## Schema

### Required

- `display_color` (String) Display color for this feed item.
- `docker_host_feed_event_type` (String) Docker Host Feed Event. Allowed values: `DockerHostCreated`, `DockerHostUpdated`, `DockerHostArchived`, `DockerHostRestored`, `OwnerUserAdded`, `OwnerUserRemoved`, `OwnerTeamAdded`, `OwnerTeamRemoved`, `OwnerRuleExecuted`, `LabelRuleExecuted`.
- `docker_host_id` (String) Relation to Docker Host ID in which this resource belongs. The ID of a `oneuptime_docker_host`.
- `feed_info_in_markdown` (String) Log of the Docker host change in Markdown.

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

Import an existing docker host feed by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_docker_host_feed.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_docker_host_feed.example <id>
```
