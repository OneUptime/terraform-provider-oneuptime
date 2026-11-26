---
page_title: "oneuptime_data_asset Resource - oneuptime"
subcategory: "Other"
description: |-
  Every table and view of your connected data stores: its columns, how many rows it holds, how big it is and when it last changed.
---

# oneuptime_data_asset (Resource)

Every table and view of your connected data stores: its columns, how many rows it holds, how big it is and when it last changed.

## Example Usage

```terraform
resource "oneuptime_data_asset" "example" {
  data_store_id = oneuptime_data_store.example.id
  asset_key     = "This is an example of longer text content that might be stored in this field."
  name          = "Example data asset"
  description   = "Managed by Terraform"
}
```

## Schema

### Required

- `asset_key` (String) The asset's name with its schema, such as "public.orders": what its metrics and checks name it by.
- `data_store_id` (String) ID of the data store this asset is in. The ID of a `oneuptime_data_store`.
- `name` (String) The table's or view's own name, such as orders.

### Optional

- `description` (String) What the asset holds and who relies on it, in your own words.
- `freshness_column` (String) A timestamp column that says when a row was written. Its latest value says when the table last changed, where the database's own counters cannot.
- `labels` (Set of String) Labels to group and filter tables by. OneUptime adds to it too, on its own: the entries listed here are added and taken off, and every other entry stays. Left out, it keeps the entries OneUptime gives it. IDs of `oneuptime_label` resources.
- `schema_name` (String) The schema (or database, or dataset) the asset is in, such as public.

### Read-Only

- `asset_kind` (String) Whether it is a table, a view, a materialized view, a partitioned table or a foreign table.
- `column_count` (Number) How many columns the asset has.
- `columns` (String) The asset's columns at its last collection: each one's name, type, whether it allows empty values, and its position. A JSON value: read it with `jsondecode()`.
- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the user who created this asset, if a person did rather than OneUptime. The ID of a `oneuptime_user` (see the data source).
- `data_store_kind` (String) What kind of store this asset is in, such as PostgreSQL or Snowflake.
- `data_store_name` (String) The name of the data store this asset is in.
- `dbt_node` (String) The dbt model, seed, snapshot or source that is this table, from the project's last uploaded artifacts: its id, materialization, tags, owner, and how its last build went. A JSON value: read it with `jsondecode()`.
- `dropped_at` (String) When OneUptime saw that the asset was gone.
- `entity_set_name` (String) A Microsoft Dataverse table's entity set: the name its Web API and Power BI's OData connectors call it by, such as accounts for the account table. Empty for every other database's tables.
- `first_seen_at` (String) When OneUptime first found the asset.
- `id` (String) Unique identifier for the resource.
- `is_dropped` (Boolean) The asset is gone from its store. It is kept with its history, and comes back if the table does.
- `last_changed_at` (String) When a row was last inserted, updated or deleted, as near as the engine tells.
- `last_collected_at` (String) When OneUptime last read the asset's health.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `last_read_at` (String) When the table was last seen read, to the hour: by the engine's read counters, or by a service's database spans.
- `late_at` (String) When the table is late if it has not changed by then: twice its usual interval after its last change, and at least fifteen minutes past it. Empty for a view, or a table whose rhythm is not known yet.
- `project_id` (String) ID of the project this data asset belongs to. The ID of a `oneuptime_project`.
- `row_count` (Number) How many rows the asset holds, from the engine's own statistics at the last collection.
- `row_count_change_tolerance_percent` (Number) A change in rows this small counts as none here, raised each time such an alert is marked not an issue.
- `size_in_bytes` (Number) How much space the asset and its indexes take.
- `unread_since` (String) The time from which the table is known to have had no read: its last read, or when the engine's read counters started when it has had none since. Empty when that cannot be said.
- `update_interval_in_seconds` (Number) How often the asset usually changes, learned from the changes OneUptime has seen.
- `updated_at` (String) Date and Time when the object was updated.
- `write_bytes_billed` (Number) BigQuery bytes billed for the jobs that write the table in the last day, a job's shared among the tables it writes.
- `write_credits` (Number) Snowflake credits the jobs that write the table spent in the last day, a job's shared among the tables it writes.

## Import

Import an existing data asset by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_data_asset.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_data_asset.example <id>
```
