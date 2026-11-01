---
page_title: "oneuptime_data_flow Data Source - oneuptime"
subcategory: "Other"
description: |-
  Where sensitive data goes: the services that read each store, where those services send data, and copies of a store's data found in another store.
---

# oneuptime_data_flow (Data Source)

Where sensitive data goes: the services that read each store, where those services send data, and copies of a store's data found in another store.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data flow may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_flow" "example" {
  flow_key = "example-flow-key"
}

# Or by id:
data "oneuptime_data_flow" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `call_count` (Number) How many calls the service map counted along the flow in its last window.
- `destination` (String) Where the data goes, when it is not a store: a host or an AI provider.
- `destination_data_store_id` (String) The store the data reaches, for a write or a copy. The ID of a `oneuptime_data_store`.
- `flow_key` (String) The flow's kind, source and destination, which identify it within the project.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_confirmed` (Boolean) A scan found the data at the destination, not only the path to it.
- `kind` (String) What the flow is: Read, Write, AI Provider, External Host or Copy.
- `overlap_percent` (Number) For a copy: how much of the source's values were found in the destination, from 0 to 100.
- `sensitivity` (String) The sensitivity of the data that can travel along it.
- `service_id` (String) The service the data passes through. The ID of a `oneuptime_service`.
- `source_data_store_id` (String) The store the data comes from. The ID of a `oneuptime_data_store`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `data_type_ids` (String) The sensitive data types that can travel along the flow. A JSON value: write it with `jsonencode()`.
- `first_seen_at` (String) When the flow was first seen.
- `last_seen_at` (String) When the flow was last seen.
- `project_id` (String) ID of the project this data flow belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
