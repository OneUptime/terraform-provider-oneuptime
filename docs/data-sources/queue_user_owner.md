---
page_title: "oneuptime_queue_user_owner Data Source - oneuptime"
subcategory: "Infrastructure"
description: |-
  Add users as owners to your queues.
---

# oneuptime_queue_user_owner (Data Source)

Add users as owners to your queues.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one queue user owner may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_queue_user_owner" "example" {
  message_queue_id = oneuptime_queue.example.id
  user_email       = "alice@example.com"
}

# Or by id:
data "oneuptime_queue_user_owner" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `message_queue_id` (String) ID of the queue the user owns. The ID of a `oneuptime_queue`.
- `user_email` (String) Look it up by the email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) ID of the user who owns the queue. The ID of a `oneuptime_user` (see the data source). Or look them up by email address, in `user_email`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
