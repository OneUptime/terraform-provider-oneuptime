---
page_title: "oneuptime_watchlist Resource - oneuptime"
subcategory: "Other"
description: |-
  Lists a team keeps that detection rules match against, such as VIP users, critical hosts, known-bad file hashes and allowed admin addresses.
---

# oneuptime_watchlist (Resource)

Lists a team keeps that detection rules match against, such as VIP users, critical hosts, known-bad file hashes and allowed admin addresses.

## Example Usage

```terraform
resource "oneuptime_watchlist" "example" {
  name        = "Example watchlist"
  entity_type = "Example short text"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `entity_type` (String) What the list holds: User, Host, IP Address, File Hash, Domain, URL, Process or Other. It decides which values the list takes - an IP address list takes addresses and ranges such as 10.20.0.0/16 - and which indicators a threat intel feed gives it.
- `name` (String) What rules call the list, between % signs: a rule reads vip_users as principalUser|expand: '%vip_users%'. Letters, digits and underscores, starting with a letter; unique in the project whatever the case.

### Optional

- `description` (String) What the list is for, for whoever reads it later.
- `source` (String) Where the values come from: Values (added by hand or imported) or Threat Intel Feed (the feed's active indicators of the list's kind). Defaults to `Values`.
- `threat_intel_feed_id` (String) For a list that follows a feed: the ID of the feed whose active indicators are the list's values. The ID of a `oneuptime_threat_intel_feed`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the person who created the watchlist. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
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
