---
page_title: "oneuptime_log_scrub_rule Resource - oneuptime"
subcategory: "Logs, Traces & Metrics"
description: |-
  Configure rules to automatically detect and scrub sensitive data (PII) from logs at ingest time.
---

# oneuptime_log_scrub_rule (Resource)

Configure rules to automatically detect and scrub sensitive data (PII) from logs at ingest time.

## Example Usage

```terraform
resource "oneuptime_log_scrub_rule" "example" {
  name         = "Example log scrub rule"
  pattern_type = "email"
  description  = "Managed by Terraform"
}
```

<a id="order"></a>
## Order

`sort_order` is this log scrub rule's place among the log scrub rules of its project, lowest first. Setting it places this log scrub rule and moves no other, so a configuration that numbers every log scrub rule of a project is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.log_scrub_rules`:

```terraform
sort_order = index(local.log_scrub_rules, each.key) + 1
```

- Give each log scrub rule a number of its own: two log scrub rules at one number are shown the older first.
- A log scrub rule created without `sort_order` goes after the last one of its project.
- Deleting a log scrub rule leaves a gap in the numbers, which changes nothing: log scrub rules go by their numbers, not by counting.
- Moving a log scrub rule in the OneUptime Dashboard numbers the project's log scrub rules 1..n again; the next plan shows that as changes to `sort_order`.


## Schema

### Required

- `name` (String) Friendly name for this scrub rule.
- `pattern_type` (String) The type of sensitive data to detect: email, creditCard, ssn, phoneNumber, ipAddress, sensitiveKeys (the whole value of every attribute whose key looks sensitive, such as password or token), or custom (the regular expression in customRegex). Any other value is refused: it would scrub nothing. Allowed values: `email`, `creditCard`, `ssn`, `phoneNumber`, `ipAddress`, `sensitiveKeys`, `custom`.

### Optional

- `custom_regex` (String) The regular expression a 'custom' rule scrubs, written without slashes or flags and matched case-sensitively. Required when patternType is 'custom': a pattern that is empty, does not compile, or matches empty text is refused. Ignored for the other pattern types.
- `description` (String) Description of what this scrub rule does.
- `fields_to_scrub` (String) Which log fields to scrub: 'body' (the log message), 'attributes' (attribute values), or 'both', the default. A sensitiveKeys rule always scrubs attribute values, whatever this says. Defaults to `both`.
- `is_enabled` (Boolean) Whether this scrub rule is active. Defaults to `true`.
- `scrub_action` (String) How to scrub matched data: 'redact' replaces it with [REDACTED] (the default), 'mask' partially hides it, 'hash' replaces it with a short hash of the value. Allowed values: `mask`, `hash`, `redact`. Defaults to `redact`.
- `sort_order` (Number) Where this rule is applied among the project's log scrub rules, lowest number first. A new rule is added to the end of the list. Setting a number places it there and moves no other one: two at one number are read the older first. In the dashboard, drag the rows to reorder them. See [Order](#order).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this log scrub rule. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this log scrub rule belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing log scrub rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_log_scrub_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_log_scrub_rule.example <id>
```
