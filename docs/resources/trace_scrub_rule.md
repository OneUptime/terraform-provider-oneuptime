---
page_title: "oneuptime_trace_scrub_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure rules to automatically detect and scrub sensitive data (PII) from spans at ingest time.
---

# oneuptime_trace_scrub_rule (Resource)

Configure rules to automatically detect and scrub sensitive data (PII) from spans at ingest time.

## Example Usage

```terraform
resource "oneuptime_trace_scrub_rule" "example" {
  name = jsonencode({
    "_type": "Name",
    "value": "John Doe"
  })
  pattern_type = "Example short text"
  scrub_action = "Example short text"
  fields_to_scrub = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name object.
- `pattern_type` (String) The type of sensitive data pattern to detect: email, creditCard, ssn, phoneNumber, ipAddress, or custom...
- `scrub_action` (String) How to scrub matched data: 'mask' partially hides it, 'hash' replaces with a hash, 'redact' removes entirely...
- `fields_to_scrub` (String) Which span fields to scrub: 'name' (span name), 'attributes' (attribute values), 'events' (span event attributes), or 'all'...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of what this scrub rule does...
- `custom_regex` (String) A custom regular expression pattern to match. Only used when patternType is 'custom'...
- `is_enabled` (Bool) Whether this scrub rule is active...
- `sort_order` (Number) Determines the evaluation order of this rule relative to others...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_trace_scrub_rule.example <id>
```
