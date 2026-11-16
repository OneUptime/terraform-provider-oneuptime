---
page_title: "oneuptime_incoming_request_secret Data Source - oneuptime"
subcategory: "Other"
description: |-
  A secret that requests to the project's Incoming Request monitors must send, in the X-OneUptime-Secret header or as a Bearer token, while the project requires one.
---

# oneuptime_incoming_request_secret (Data Source)

A secret that requests to the project's Incoming Request monitors must send, in the X-OneUptime-Secret header or as a Bearer token, while the project requires one.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incoming request secret may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incoming_request_secret" "example" {
  name = "Example incoming request secret"
}

# Or by id:
data "oneuptime_incoming_request_secret" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `name` (String) What the secret is for, so you can tell it apart from the project's other incoming request secrets - the senders that use it, say.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_used_at` (String) About when a request was last accepted with this secret, to the nearest five minutes. Empty means none has been accepted with it yet. Use it to find secrets that are safe to delete.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `secret_value` (String, Sensitive) The value senders send in the X-OneUptime-Secret header, or as a Bearer token in the Authorization header. Leave it empty and OneUptime generates a strong one. Set your own with 16 to 256 characters: letters, digits and keyboard symbols, no spaces. It cannot be changed after the secret is created - create a new secret instead.
- `updated_at` (String) Date and Time when the object was updated.
