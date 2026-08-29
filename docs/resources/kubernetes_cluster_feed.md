---
page_title: "oneuptime_kubernetes_cluster_feed Resource - oneuptime"
subcategory: "Other"
description: |-
  Log of everything that happened to this Kubernetes cluster - creation, updates, owner changes and the rules that made them.
---

# oneuptime_kubernetes_cluster_feed (Resource)

Log of everything that happened to this Kubernetes cluster - creation, updates, owner changes and the rules that made them.

## Example Usage

```terraform
resource "oneuptime_kubernetes_cluster_feed" "example" {
  kubernetes_cluster_id = "123e4567-e89b-12d3-a456-426614174000"
  feed_info_in_markdown = "# Heading

This is **markdown** content"
  kubernetes_cluster_feed_event_type = "Example short text"
  display_color = jsonencode({
    "_type": "Color",
    "value": "#ff0000"
  })
}
```

## Schema

### Required

- `kubernetes_cluster_id` (String) A unique identifier for an object, represented as a UUID..
- `feed_info_in_markdown` (String) Log of the Kubernetes cluster change in Markdown..
- `kubernetes_cluster_feed_event_type` (String) Kubernetes Cluster Feed Event..
- `display_color` (String) Color object.

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `more_information_in_markdown` (String) More information in Markdown..
- `user_id` (String) A unique identifier for an object, represented as a UUID..
- `posted_at` (String) A date time object..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_kubernetes_cluster_feed.example <id>
```
