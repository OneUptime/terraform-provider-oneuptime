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
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name object.
- `pattern_type` (String) The type of sensitive data to detect: email, creditCard, ssn, phoneNumber, ipAddress, sensitiveKeys (the whole value of every attribute whose key looks sensitive, such as password or token), or custom (the regular expression in customRegex). Any other value is refused: it would scrub nothing...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of what this scrub rule does...
- `custom_regex` (String) The regular expression a 'custom' rule scrubs, written without slashes or flags and matched case-sensitively. Required when patternType is 'custom': a pattern that is empty, does not compile, or matches empty text is refused. Ignored for the other pattern types...
- `scrub_action` (String) How to scrub matched data: 'redact' replaces it with [REDACTED] (the default), 'mask' partially hides it, 'hash' replaces it with a short hash of the value...
- `fields_to_scrub` (String) Which span fields to scrub: 'name' (the span name), 'attributes' (attribute values), 'events' (span event attributes), or 'all', the default. A sensitiveKeys rule always scrubs attribute and event attribute values, whatever this says...
- `is_enabled` (Bool) Whether this scrub rule is active...
- `sort_order` (Number) Where this rule is applied among the project's span scrub rules, lowest number first. A new rule is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them...

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
