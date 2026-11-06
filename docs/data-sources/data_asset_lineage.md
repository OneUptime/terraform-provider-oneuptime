---
page_title: "oneuptime_data_asset_lineage Data Source - oneuptime"
subcategory: "Other"
description: |-
  One table feeding another, from the pipelines that read the first and write the second, the database's query history, or a view's definition.
---

# oneuptime_data_asset_lineage (Data Source)

One table feeding another, from the pipelines that read the first and write the second, the database's query history, or a view's definition.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data asset lineage may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_asset_lineage" "example" {
  upstream_data_asset_id = oneuptime_data_asset.example.id
}

# Or by id:
data "oneuptime_data_asset_lineage" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `data_pipeline_id` (String) ID of the pipeline that moves the data. The ID of a `oneuptime_data_pipeline` (see the data source).
- `downstream_data_asset_id` (String) ID of the table built from it. The ID of a `oneuptime_data_asset`.
- `edge_key` (String) The upstream table, the downstream table and the pipeline together: one row each.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `seen_by` (String) For an edge from query history: a login or service account whose queries moved the data.
- `source` (String) Where the edge was seen: a pipeline's run, the database's query history, or a view's definition. Empty means a pipeline's run.
- `upstream_data_asset_id` (String) ID of the table the data comes from; empty for a pipeline that writes the table from data OneUptime does not watch. The ID of a `oneuptime_data_asset`.

### Read-Only

- `column_mappings` (String) Which columns of the upstream table feed which columns of the downstream one, when its source names them. A JSON value: read it with `jsondecode()`.
- `created_at` (String) Date and Time when the object was created.
- `first_seen_at` (String) When the edge was first seen.
- `last_seen_at` (String) When a run or query last moved data along the edge, or a view was last seen defined on it.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
