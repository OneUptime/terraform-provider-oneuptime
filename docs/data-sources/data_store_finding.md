---
page_title: "oneuptime_data_store_finding Data Source - oneuptime"
subcategory: "Other"
description: |-
  What Data Security found in a store: one data type at one location, with its confidence, counts and masked samples.
---

# oneuptime_data_store_finding (Data Source)

What Data Security found in a store: one data type at one location, with its confidence, counts and masked samples.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data store finding may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_store_finding" "example" {
  data_store_id = oneuptime_data_store.example.id
}

# Or by id:
data "oneuptime_data_store_finding" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `ai_reason` (String) OneUptime AI's reason for its suggestion, in a sentence.
- `confidence` (String) How sure the scan is: High, Medium or Low.
- `data_category` (String) The data type's category.
- `data_store_id` (String) The store the finding is in. The ID of a `oneuptime_data_store`.
- `data_type_id` (String) The data type found: a built-in id such as credit-card-number, or custom: and a custom type's id.
- `data_type_name` (String) The data type's name when it was found.
- `distinct_values` (Number) How many different values matched, estimated from hashes that never leave the scan.
- `finding_key` (String) The data type and location, which identify the finding within its store.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_ai_suggested` (Boolean) OneUptime AI suggested this data type for a field no detector placed, from the field's name, the names beside it and the shapes of its values - never a value. It counts only once a person confirms it.
- `is_confirmed` (Boolean) A person confirmed the finding holds this data. A confirmed finding counts as High confidence toward the store's sensitivity and checks, whatever the scan or OneUptime AI was sure of.
- `is_found_by_name` (Boolean) Found from the location's name alone, for a data type no value pattern can confirm.
- `is_free_text` (Boolean) The location holds free text, so the finding is about values inside it rather than the whole field.
- `is_name_hinted` (Boolean) The location's name hints at the data type.
- `is_scrubbed` (Boolean) For telemetry: a scrub rule already masks this data at ingest.
- `is_suppressed` (Boolean) A person marked the finding as not sensitive. It no longer counts toward the store's sensitivity or checks.
- `last_scan_id` (String) The scan that last saw the finding. The ID of a `oneuptime_data_security_scan` (see the data source).
- `location` (String) Where in the store it was found: a table and column, a field path, an object path pattern, or a telemetry attribute.
- `location_kind` (String) What the location is: Column, Field, Object, Attribute, Body or Text.
- `match_count` (Number) How many sampled values matched.
- `match_percent` (Number) The share of sampled values that matched, from 0 to 100.
- `sensitivity` (String) The data type's sensitivity.
- `suppressed_reason` (String) Why the finding was marked not sensitive.
- `test_value_count` (Number) How many matches were well-known test values, which are never findings.
- `value_count` (Number) How many values were sampled at the location.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `first_seen_at` (String) When a scan first found it.
- `last_seen_at` (String) When a scan last found it.
- `masked_samples` (String) At most three masked samples, such as j***@e******.com. Never a raw value. A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of the project this finding belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
