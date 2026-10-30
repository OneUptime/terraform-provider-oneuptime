---
page_title: "oneuptime_watchlist Data Source - oneuptime"
subcategory: "Other"
description: |-
  Lists a team keeps that detection rules match against, such as VIP users, critical hosts, known-bad file hashes and allowed admin addresses.
---

# oneuptime_watchlist (Data Source)

Lists a team keeps that detection rules match against, such as VIP users, critical hosts, known-bad file hashes and allowed admin addresses.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one watchlist may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_watchlist" "example" {
  name = "Example watchlist"
}

# Or by id:
data "oneuptime_watchlist" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the person who created the watchlist. The ID of a `oneuptime_user` (see the data source).
- `description` (String) What the list is for, for whoever reads it later.
- `entity_type` (String) What the list holds: User, Host, IP Address, File Hash, Domain, URL, Process or Other. It decides which values the list takes - an IP address list takes addresses and ranges such as 10.20.0.0/16 - and which indicators a threat intel feed gives it.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `name` (String) What rules call the list, between % signs: a rule reads vip_users as principalUser|expand: '%vip_users%'. Letters, digits and underscores, starting with a letter; unique in the project whatever the case.
- `source` (String) Where the values come from: Values (added by hand or imported) or Threat Intel Feed (the feed's active indicators of the list's kind).
- `threat_intel_feed_id` (String) For a list that follows a feed: the ID of the feed whose active indicators are the list's values. The ID of a `oneuptime_threat_intel_feed`.
- `value_count` (Number) How many values the list holds, counted by OneUptime as values are added and removed. 0 for a list that follows a threat intel feed.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this watchlist belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
