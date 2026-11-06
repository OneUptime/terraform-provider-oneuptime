---
page_title: "oneuptime_data_product Data Source - oneuptime"
subcategory: "Other"
description: |-
  A named group of watched tables with an owner, its consumers and a freshness SLA, on a status page and tracked as an SLO.
---

# oneuptime_data_product (Data Source)

A named group of watched tables with an owner, its consumers and a freshness SLA, on a status page and tracked as an SLO.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data product may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_product" "example" {
  owner_team_id    = oneuptime_team.example.id
  owner_user_email = "alice@example.com"
}

# Or by id:
data "oneuptime_data_product" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `consumers` (String) Who uses the data: teams, dashboards, reports, AI agents.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) What the data is, and what it is for.
- `fresh_by_time_of_day` (String) The time of day, in UTC (HH:MM), by which each day's data lands. Empty for no daily deadline.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `late_table_count` (Number) How many of the tables are late.
- `max_data_age_in_minutes` (Number) How old the data may be, in minutes, before the product is late. Empty for no limit.
- `monitor_id` (String) ID of the Metrics monitor that is Degraded, with an incident, while a table is late. The ID of a `oneuptime_monitor`.
- `name` (String) What the data product is called, as its consumers know it.
- `owner_team_id` (String) ID of the team that owns the data product. The ID of a `oneuptime_team`.
- `owner_user_email` (String) Look it up by the email address of the person in `owner_user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `owner_user_id` (String) ID of the person who owns the data product. The ID of a `oneuptime_user` (see the data source). Or look them up by email address, in `owner_user_email`.
- `service_level_objective_id` (String) ID of the SLO that tracks the data product's SLA, with burn-rate alerts. The ID of a `oneuptime_service_level_objective`.
- `sla_target_percentage` (Number) The share of time the data must be on time, which its SLO tracks: 99 means late for at most 1% of the time.
- `status` (String) Whether the data is fresh by the SLA: On Time, Late, or Unknown while no table has been seen changing.

### Read-Only

- `badge_token` (String, Sensitive) The secret in the badge's URL, so a badge can be embedded where people cannot sign in to OneUptime.
- `created_at` (String) Date and Time when the object was created.
- `data_as_of` (String) The oldest last change among the tables: what the data is as of.
- `data_assets` (Set of String) The watched tables the data product is made of. IDs of `oneuptime_data_asset` resources.
- `late_data_asset_ids` (String) The IDs of the tables that are late. A JSON value: read it with `jsondecode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `status_changed_at` (String) When the status last changed.
- `status_checked_at` (String) When the tables were last checked against the SLA.
- `updated_at` (String) Date and Time when the object was updated.
