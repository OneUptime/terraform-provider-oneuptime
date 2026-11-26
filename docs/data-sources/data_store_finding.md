---
page_title: "oneuptime_data_store_finding Data Source - oneuptime"
subcategory: "Data Security"
description: |-
  What Data Security found in a store: one data type at one location, with its confidence, counts and masked samples.
---

# oneuptime_data_store_finding (Data Source)

What Data Security found in a store: one data type at one location, with its confidence, counts and masked samples.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data store finding may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_store_finding" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_data_store_finding" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `ai_reason` (String) OneUptime AI's reason for its suggestion, in a sentence.
- `ai_verdict` (String) OneUptime AI's verdict on a finding that waits for a person: Confirm, NotSensitive or CannotTell. A person decides.
- `confidence` (String) How sure the scan is: High, Medium or Low.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `data_category` (String) The data type's category.
- `data_store_id` (String) The store the finding is in. The ID of a `oneuptime_data_store`.
- `data_type_id` (String) The data type found: a built-in id such as credit-card-number, or custom: and a custom type's id.
- `data_type_name` (String) The data type's name when it was found.
- `distinct_values` (Number) How many different values matched, estimated from hashes that never leave the scan.
- `finding_key` (String) The data type and location, which identify the finding within its store.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_ai_suggested` (Boolean) OneUptime AI suggested this data type for a field no detector placed, from the field's name, the names beside it and the shapes of its values - never a value. It counts only once a person confirms it.
- `is_confirmed` (Boolean) A person confirmed the finding holds this data. A confirmed finding counts as High confidence toward the store's sensitivity and checks, whatever the scan or OneUptime AI was sure of.
- `is_confirmed_by_ai` (Boolean) OneUptime AI confirmed the finding without asking, as one of the project's AI rules allows. Undo it from Data Security > AI.
- `is_found_by_model` (Boolean) Found in free text by the names-and-addresses model that runs inside the scanner, not by a pattern.
- `is_found_by_name` (Boolean) Found from the location's name alone, for a data type no value pattern can confirm.
- `is_free_text` (Boolean) The location holds free text, so the finding is about values inside it rather than the whole field.
- `is_name_hinted` (Boolean) The location's name hints at the data type.
- `is_scrubbed` (Boolean) For telemetry: a scrub rule already masks this data at ingest.
- `is_suppressed` (Boolean) A person marked the finding as not sensitive. It no longer counts toward the store's sensitivity or checks.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
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

- `ai_verdict_analysis` (String) OneUptime AI's reason for its verdict, a reason a person could give, and the earlier decisions on findings like it that it cites. A JSON value: read it with `jsondecode()`.
- `ai_verdict_at` (String) When OneUptime AI last gave its verdict on the finding.
- `created_at` (String) Date and Time when the object was created.
- `decided_at` (String) When the finding was last confirmed or marked not sensitive.
- `first_seen_at` (String) When a scan first found it.
- `fixed_at` (String) For telemetry: when a code fix was verified to have stopped this finding in production. It no longer counts toward its issue unless a scan finds it again.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_seen_at` (String) When a scan last found it.
- `masked_samples` (String) At most three masked samples, such as j***@e******.com. Never a raw value. A JSON value: read it with `jsondecode()`.
- `project_id` (String) ID of the project this finding belongs to. The ID of a `oneuptime_project`.
- `sample_records` (String) Telemetry only: up to ten records the last scan found it in, by when each was written and its trace, span or session. Never a value. A JSON value: read it with `jsondecode()`.
- `updated_at` (String) Date and Time when the object was updated.
