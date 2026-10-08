---
page_title: "oneuptime_log_scrub_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Configure rules to automatically detect and scrub sensitive data (PII) from logs at ingest time.
---

# oneuptime_log_scrub_rule (Data Source)

Configure rules to automatically detect and scrub sensitive data (PII) from logs at ingest time.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one log scrub rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_log_scrub_rule" "example" {
  name = "Example log scrub rule"
}

# Or by id:
data "oneuptime_log_scrub_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the user who created this log scrub rule. The ID of a `oneuptime_user` (see the data source).
- `custom_regex` (String) The regular expression a 'custom' rule scrubs, written without slashes or flags and matched case-sensitively. Required when patternType is 'custom': a pattern that is empty, does not compile, or matches empty text is refused. Ignored for the other pattern types.
- `description` (String) Description of what this scrub rule does.
- `fields_to_scrub` (String) Which log fields to scrub: 'body' (the log message), 'attributes' (attribute values), or 'both', the default. A sensitiveKeys rule always scrubs attribute values, whatever this says.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this scrub rule is active.
- `name` (String) Friendly name for this scrub rule.
- `pattern_type` (String) The type of sensitive data to detect: email, creditCard, ssn, phoneNumber, ipAddress, sensitiveKeys (the whole value of every attribute whose key looks sensitive, such as password or token), or custom (the regular expression in customRegex). Any other value is refused: it would scrub nothing.
- `scrub_action` (String) How to scrub matched data: 'redact' replaces it with [REDACTED] (the default), 'mask' partially hides it, 'hash' replaces it with a short hash of the value.
- `sort_order` (Number) Where this rule is applied among the project's log scrub rules, lowest number first. A new rule is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this log scrub rule belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
