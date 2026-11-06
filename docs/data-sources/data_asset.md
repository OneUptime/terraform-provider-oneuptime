---
page_title: "oneuptime_data_asset Data Source - oneuptime"
subcategory: "Other"
description: |-
  Every table and view of your connected data stores: its columns, how many rows it holds, how big it is and when it last changed.
---

# oneuptime_data_asset (Data Source)

Every table and view of your connected data stores: its columns, how many rows it holds, how big it is and when it last changed.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data asset may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_asset" "example" {
  name = "Example data asset"
}

# Or by id:
data "oneuptime_data_asset" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `asset_key` (String) The asset's name with its schema, such as "public.orders": what its metrics and checks name it by.
- `asset_kind` (String) Whether it is a table, a view, a materialized view, a partitioned table or a foreign table.
- `column_count` (Number) How many columns the asset has.
- `created_by_user_id` (String) ID of the user who created this asset, if a person did rather than OneUptime. The ID of a `oneuptime_user` (see the data source).
- `data_store_id` (String) ID of the data store this asset is in. The ID of a `oneuptime_data_store`.
- `data_store_kind` (String) What kind of store this asset is in, such as PostgreSQL or Snowflake.
- `data_store_name` (String) The name of the data store this asset is in.
- `description` (String) What the asset holds and who relies on it, in your own words.
- `freshness_column` (String) A timestamp column that says when a row was written. Its latest value says when the table last changed, where the database's own counters cannot.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_dropped` (Boolean) The asset is gone from its store. It is kept with its history, and comes back if the table does.
- `name` (String) The table's or view's own name, such as orders.
- `row_count` (Number) How many rows the asset holds, from the engine's own statistics at the last collection.
- `row_count_change_tolerance_percent` (Number) A change in rows this small counts as none here, raised each time such an alert is marked not an issue.
- `schema_name` (String) The schema (or database, or dataset) the asset is in, such as public.
- `size_in_bytes` (Number) How much space the asset and its indexes take.
- `update_interval_in_seconds` (Number) How often the asset usually changes, learned from the changes OneUptime has seen.
- `write_bytes_billed` (Number) BigQuery bytes billed for the jobs that write the table in the last day, a job's shared among the tables it writes.
- `write_credits` (Number) Snowflake credits the jobs that write the table spent in the last day, a job's shared among the tables it writes.

### Read-Only

- `columns` (String) The asset's columns at its last collection: each one's name, type, whether it allows empty values, and its position. A JSON value: read it with `jsondecode()`.
- `created_at` (String) Date and Time when the object was created.
- `dbt_node` (String) The dbt model, seed, snapshot or source that is this table, from the project's last uploaded artifacts: its id, materialization, tags, owner, and how its last build went. A JSON value: read it with `jsondecode()`.
- `dropped_at` (String) When OneUptime saw that the asset was gone.
- `first_seen_at` (String) When OneUptime first found the asset.
- `labels` (Set of String) Labels to group and filter tables by. IDs of `oneuptime_label` resources.
- `last_changed_at` (String) When a row was last inserted, updated or deleted, as near as the engine tells.
- `last_collected_at` (String) When OneUptime last read the asset's health.
- `last_read_at` (String) When the table was last seen read, to the hour: by the engine's read counters, or by a service's database spans.
- `project_id` (String) ID of the project this data asset belongs to. The ID of a `oneuptime_project`.
- `unread_since` (String) The time from which the table is known to have had no read: its last read, or when the engine's read counters started when it has had none since. Empty when that cannot be said.
- `updated_at` (String) Date and Time when the object was updated.
