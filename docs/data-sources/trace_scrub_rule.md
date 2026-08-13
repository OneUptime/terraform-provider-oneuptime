---
page_title: "oneuptime_trace_scrub_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Configure rules to automatically detect and scrub sensitive data (PII) from spans at ingest time.
---

# oneuptime_trace_scrub_rule (Data Source)

Configure rules to automatically detect and scrub sensitive data (PII) from spans at ingest time. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_trace_scrub_rule" "by_name" {
  name = "example-trace_scrub_rule"
}

data "oneuptime_trace_scrub_rule" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `description` (String) Description of what this scrub rule does... Computed.
- `pattern_type` (String) The type of sensitive data pattern to detect: email, creditCard, ssn, phoneNumber, ipAddress, or custom... Computed.
- `custom_regex` (String) A custom regular expression pattern to match. Only used when patternType is 'custom'... Computed.
- `scrub_action` (String) How to scrub matched data: 'mask' partially hides it, 'hash' replaces with a hash, 'redact' removes entirely... Computed.
- `fields_to_scrub` (String) Which span fields to scrub: 'name' (span name), 'attributes' (attribute values), 'events' (span event attributes), or 'all'... Computed.
- `is_enabled` (Bool) Whether this scrub rule is active... Computed.
- `sort_order` (Number) Determines the evaluation order of this rule relative to others... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
