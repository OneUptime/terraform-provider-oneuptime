---
page_title: "oneuptime_data_asset_usage Data Source - oneuptime"
subcategory: "Other"
description: |-
  A service seen reading or writing a data asset in its database spans, with the operations it ran and when it last ran them.
---

# oneuptime_data_asset_usage (Data Source)

A service seen reading or writing a data asset in its database spans, with the operations it ran and when it last ran them.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data asset usage may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_asset_usage" "example" {
  data_asset_id = oneuptime_data_asset.example.id
}

# Or by id:
data "oneuptime_data_asset_usage" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `access_kind` (String) Whether the service writes the table's rows or reads them.
- `call_count` (Number) How many calls the service made to the table in the last hour.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `data_asset_id` (String) ID of the table or view the service read or wrote. The ID of a `oneuptime_data_asset`.
- `data_store_id` (String) ID of the data store the table is in.
- `error_count` (Number) How many of those calls failed.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_ai_app` (Boolean) Whether the service also calls a model: reading this table, it is an AI app whose answers depend on it.
- `service_id` (String) ID of the service whose spans read or wrote the table. The ID of a `oneuptime_service`.
- `service_name` (String) The service's name when it was last seen using the table.
- `usage_key` (String) The service and the way it uses the table, one row per pair: what the next discovery updates.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `endpoints` (String) The service's busiest endpoints that made these calls, from the server spans of the same traces. A JSON value: read it with `jsondecode()`.
- `first_seen_at` (String) When the service was first seen using the table this way.
- `last_seen_at` (String) When the service was last seen using the table this way.
- `operations` (String) The statements the service ran on the table, such as INSERT or SELECT. A JSON value: read it with `jsondecode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
