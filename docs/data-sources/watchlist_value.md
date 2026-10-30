---
page_title: "oneuptime_watchlist_value Data Source - oneuptime"
subcategory: "Other"
description: |-
  The values on a watchlist the team keeps, each with an optional note, matched by the detection rules that read the list.
---

# oneuptime_watchlist_value (Data Source)

The values on a watchlist the team keeps, each with an optional note, matched by the detection rules that read the list.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one watchlist value may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_watchlist_value" "example" {
  security_watchlist_id = oneuptime_watchlist.example.id
}

# Or by id:
data "oneuptime_watchlist_value" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the person who added the value. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `note` (String) Why the value is on the list, for whoever reads it later.
- `security_watchlist_id` (String) ID of the watchlist this value is on. The ID of a `oneuptime_watchlist`.
- `value` (String) The value, at most 256 characters, matched whatever the case. * stands for any run of characters (svc-* is every service account), and an IP address list takes ranges such as 10.20.0.0/16.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this value belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
