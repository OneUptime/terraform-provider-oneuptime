---
page_title: "oneuptime_watchlist Resource - oneuptime"
subcategory: "Security"
description: |-
  Lists a team keeps that detection rules match against, such as VIP users, critical hosts, known-bad file hashes and allowed admin addresses.
---

# oneuptime_watchlist (Resource)

Lists a team keeps that detection rules match against, such as VIP users, critical hosts, known-bad file hashes and allowed admin addresses.

## Example Usage

```terraform
resource "oneuptime_watchlist" "example" {
  name        = "Example watchlist"
  entity_type = "User"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `entity_type` (String) What the list holds: User, Host, IP Address, File Hash, Domain, URL, Process or Other. It decides which values the list takes - an IP address list takes addresses and ranges such as 10.20.0.0/16 - and which indicators a threat intel feed gives it. Allowed values: `User`, `Host`, `IP Address`, `Domain`, `URL`, `File Hash`, `Process`, `Data Store`, `Other`.
- `name` (String) What rules call the list, between % signs: a rule reads vip_users as principalUser|expand: '%vip_users%'. Letters, digits and underscores, starting with a letter; unique in the project whatever the case.

### Optional

- `description` (String) What the list is for, for whoever reads it later.
- `import_schedule` (String) For a list read from an address: how often it is read - Every Hour (the default), Every 6 Hours or Every Day. Allowed values: `Every Hour`, `Every 6 Hours`, `Every Day`.
- `import_url` (String) For a list read from an address: the https address of a text or CSV file, one value a line, with an optional note and when the value leaves the list. Read again on the list's schedule; the list is what the file held at its last read.
- `source` (String) Where the values come from: Values (added by hand or imported), URL (a file at an https address, read again on a schedule) or Threat Intel Feed (the feed's active indicators of the list's kind). Allowed values: `Values`, `Threat Intel Feed`, `URL`. Defaults to `Values`.
- `threat_intel_feed_id` (String) For a list that follows a feed: the ID of the feed whose active indicators are the list's values. The ID of a `oneuptime_threat_intel_feed`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the person who created the watchlist. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_import_at` (String) For a list read from an address: when OneUptime last read it, or tried to.
- `last_import_error` (String) For a list read from an address: why its last read failed. The list keeps the values of the read before. Empty when it worked.
- `last_import_error_detail` (String) For a list read from an address: what the address or the network said when its last read failed, as it said it.
- `last_import_lines_left_out` (Number) For a list read from an address: how many lines of the file its last read left out - lines that are not values of the list's kind, or whose date to leave the list had passed.
- `last_import_succeeded_at` (String) For a list read from an address: when it was last read without an error. The list holds what the file held then.
- `next_import_at` (String) For a list read from an address: when it is next read. Empty for any other list.
- `project_id` (String) ID of the project this watchlist belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
- `value_count` (Number) How many values the list holds, counted by OneUptime as values are added and removed. 0 for a list that follows a threat intel feed.

## Import

Import an existing watchlist by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_watchlist.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_watchlist.example <id>
```
