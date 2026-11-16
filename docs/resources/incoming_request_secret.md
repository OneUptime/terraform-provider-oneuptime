---
page_title: "oneuptime_incoming_request_secret Resource - oneuptime"
subcategory: "Other"
description: |-
  A secret that requests to the project's Incoming Request monitors must send, in the X-OneUptime-Secret header or as a Bearer token, while the project requires one.
---

# oneuptime_incoming_request_secret (Resource)

A secret that requests to the project's Incoming Request monitors must send, in the X-OneUptime-Secret header or as a Bearer token, while the project requires one.

## Example Usage

```terraform
resource "oneuptime_incoming_request_secret" "example" {
  name = "Example incoming request secret"
}
```

## Schema

### Required

- `name` (String) What the secret is for, so you can tell it apart from the project's other incoming request secrets - the senders that use it, say.

### Optional

- `secret_value` (String, Sensitive) The value senders send in the X-OneUptime-Secret header, or as a Bearer token in the Authorization header. Leave it empty and OneUptime generates a strong one. Set your own with 16 to 256 characters: letters, digits and keyboard symbols, no spaces. It cannot be changed after the secret is created - create a new secret instead.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_used_at` (String) About when a request was last accepted with this secret, to the nearest five minutes. Empty means none has been accepted with it yet. Use it to find secrets that are safe to delete.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incoming request secret by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incoming_request_secret.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incoming_request_secret.example <id>
```
