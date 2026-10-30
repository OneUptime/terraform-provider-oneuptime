---
page_title: "oneuptime_watchlist_value Resource - oneuptime"
subcategory: "Other"
description: |-
  The values on a watchlist the team keeps, each with an optional note, matched by the detection rules that read the list.
---

# oneuptime_watchlist_value (Resource)

The values on a watchlist the team keeps, each with an optional note, matched by the detection rules that read the list.

## Example Usage

```terraform
resource "oneuptime_watchlist_value" "example" {
  security_watchlist_id = oneuptime_watchlist.example.id
  value                 = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `security_watchlist_id` (String) ID of the watchlist this value is on. The ID of a `oneuptime_watchlist`.
- `value` (String) The value, at most 256 characters, matched whatever the case. * stands for any run of characters (svc-* is every service account), and an IP address list takes ranges such as 10.20.0.0/16.

### Optional

- `note` (String) Why the value is on the list, for whoever reads it later.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the person who added the value. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this value belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing watchlist value by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_watchlist_value.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_watchlist_value.example <id>
```
